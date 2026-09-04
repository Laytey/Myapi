package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const ( // нельзя изменить после объявления
	ctxTimeout = 5 * time.Second // контекст с таймаутом (ctxTimeout)
	// программа будет ждать ответа от базы данных. в течение этого времени
)

type Storage struct {
	Pool *pgxpool.Pool
	// Storage хранит пул соединений (*pgxpool.Pool)
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	// Data Source Name (имя источника данных)
	// dsn — это строка подключения к базе данных. Она приходит из main.go (из переменной окружения DATABASE_URL)
	// New создаёт пул, проверяет подключение и возвращает его
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
		// %w - wrap оборачивает ошибку (позволяет проверять её позже)
	}

	config.MaxConns = 10
	// config — это структура, содержащая все настройки подключения к PostgreSQL:
	// хост, порт, пользователь, пароль, база данных, параметры SSL, ограничения пула и т.д.

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		// ping — проверяет подключение к PostgreSQL

		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	return &Storage{
		Pool: pool,
	}, nil
}
