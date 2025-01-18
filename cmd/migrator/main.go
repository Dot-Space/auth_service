package main

import (
	"flag"
	"fmt"

	"github.com/Dot-Space/auth_service/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationsPath, migrationsTable string

	flag.StringVar(&migrationsPath, "migrations_path", "", "Path to migrations")
	flag.StringVar(&migrationsTable, "migrations_table", "", "Name of migrations table")
	flag.Parse()

	if migrationsPath == "" {
		panic("Migrations path is required!")
	}

	dbConfig := config.LoadDbConfig()
	pgConnInfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)

	fmt.Println(pgConnInfo)

	m, err := migrate.New(
		"file://"+migrationsPath,
		pgConnInfo,
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		panic(err)
	}
}
