package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://todo_user:todo_pass@localhost:5432/todo_db?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("не удалось разобрать строку подключения: %w", err)
	}

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("не удалось создать пул соединений: %w", err)
	}

	if err = Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	// Создаём таблицы, если их нет
	return createTables()
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id    SERIAL PRIMARY KEY,
            login TEXT UNIQUE NOT NULL,
            pass  TEXT NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS tasks (
            id      SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL,
            text    TEXT NOT NULL,
            done    BOOLEAN NOT NULL DEFAULT FALSE
        );`,
	}
	for _, q := range queries {
		if _, err := Pool.Exec(context.Background(), q); err != nil {
			return err
		}
	}
	return nil
}

func CloseDB() {
	if Pool != nil {
		Pool.Close()
	}
}
