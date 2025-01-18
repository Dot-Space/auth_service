package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Dot-Space/auth_service/internal/models"
	"github.com/lib/pq"
)

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "db.SaveUser"

	fmt.Println(email, passHash)

	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email, passHash)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
			}
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, email string) (models.User, error) {
	const op = "db.GetUser"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = $1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	user := models.User{}
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
