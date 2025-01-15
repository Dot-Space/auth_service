package db

import (
	"database/sql"

	"fmt"

	"github.com/Dot-Space/auth_service/config"
)

type Storage struct {
	db *sql.DB
}

func New(config *config.DBConfig) (*Storage, error) {
	const op = "db.New"

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=&s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
