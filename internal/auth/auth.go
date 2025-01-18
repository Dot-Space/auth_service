package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Dot-Space/auth_service/internal/db"
	"github.com/Dot-Space/auth_service/internal/models"
	"github.com/Dot-Space/auth_service/internal/pkg/ecfs"
	"github.com/Dot-Space/auth_service/internal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrUserAlreadyExists  = errors.New("User already exists!")
)

type Auth struct {
	log             *slog.Logger
	storageProvider StorageProvider
	tokenTTL        time.Duration
	jwtSecret       string
}

type StorageProvider interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
	GetUser(ctx context.Context, email string) (models.User, error)
}

func New(
	log *slog.Logger,
	storageProvider StorageProvider,
	tokenTTL time.Duration,
	jwtSecret string,
) *Auth {
	return &Auth{
		log:             log,
		storageProvider: storageProvider,
		tokenTTL:        tokenTTL,
		jwtSecret:       jwtSecret,
	}
}

func (a *Auth) RegisterUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("Registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to generate hash for password", ecfs.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.storageProvider.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, db.ErrUserExists) {
			a.log.Warn("User already exists", ecfs.Err(err))

			return 0, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}
		log.Error("Failed to save user", ecfs.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) Login(ctx context.Context, email string, password string) (string, string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	user, err := a.storageProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			a.log.Warn("User not found", ecfs.Err(err))

			return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("Failed to get user")

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("Invalid credentials", ecfs.Err(err))

		return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("User logged in successfully")

	token, err := jwt.NewToken(user, a.jwtSecret, a.tokenTTL, "access")
	if err != nil {
		a.log.Error("Couldnt create access token", ecfs.Err(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err := jwt.NewToken(user, a.jwtSecret, a.tokenTTL, "refresh")
	if err != nil {
		a.log.Error("Couldnt create refresh token", ecfs.Err(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return token, refreshToken, nil
}

func (a *Auth) CheckToken(ctx context.Context, token string, tokenType string) (bool, error) {
	const op = "auth.CheckToken"

	isValid, err := jwt.ValidateToken(token, a.jwtSecret, tokenType)
	if err != nil {
		return isValid, fmt.Errorf("%s: %w", op, err)
	}

	return isValid, nil
}

func (a *Auth) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	const op = "auth.RefreshToken"

	isValid, err := jwt.ValidateToken(refreshToken, a.jwtSecret, "refresh")

	if !isValid {
		return "", "", fmt.Errorf("%s: %s", op, "InvalidToken")
	}

	newAccessToken, newRefreshToken, err := jwt.RefreshToken(refreshToken, a.jwtSecret, a.tokenTTL)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return newAccessToken, newRefreshToken, nil
}
