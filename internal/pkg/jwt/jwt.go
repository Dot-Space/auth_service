package jwt

import (
	"fmt"
	"time"

	"github.com/Dot-Space/auth_service/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Id        int64  `json:"id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	Exp       int64  `json:"exp"`
	jwt.MapClaims
}

func CreateToken(userClaims *UserClaims, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = *userClaims

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewToken(user models.User, secret string, duration time.Duration, tokenType string) (string, error) {
	userClaims := &UserClaims{
		Id:        user.ID,
		Email:     user.Email,
		TokenType: tokenType,
		Exp:       time.Now().Add(duration).Unix(),
	}

	tokenString, err := CreateToken(userClaims, secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(token string, secret string, tokenType string) (bool, error) {
	userClaims := &UserClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, userClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if userClaims.TokenType != tokenType {
		return false, fmt.Errorf("Wrong type of token")
	}

	if err != nil {
		return false, err
	}

	return parsedToken.Valid, nil
}

func RefreshToken(token string, secret string, duration time.Duration) (string, error) {
	userClaims := &UserClaims{}

	_, err := jwt.ParseWithClaims(token, userClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	userClaims.Exp = time.Now().Add(duration).Unix()

	refreshedToken, err := CreateToken(userClaims, secret)
	if err != nil {
		return "", err
	}

	return refreshedToken, nil
}
