package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const ( // нельзя изменить после объявления
	ctxTimeout = 5 * time.Second // контекст с таймаутом (ctxTimeout)
	// программа будет ждать ответа от базы данных. в течение этого времени
	maxRetries = 5               // количество попыток подключения
	retryDelay = 2 * time.Second // задержка между попытками
)

type Storage struct {
	Pool *pgxpool.Pool
	// Storage хранит пул соединений (*pgxpool.Pool)
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	// Data Source Name (имя источника данных)
	// dsn — это строка подключения к базе данных. Она приходит из main.go (из переменной окружения DATABASE_URL)
	// New создаёт пул, проверяет подключение и возвращает его

	var pool *pgxpool.Pool
	var err error

	// Retry-цикл: пытаемся подключиться до maxRetries раз
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Попытка подключения к PostgreSQL %d/%d", attempt, maxRetries)

		config, parseErr := pgxpool.ParseConfig(dsn)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse config: %w", parseErr)
			// %w - wrap оборачивает ошибку (позволяет проверять её позже)
		}

		config.MaxConns = 10
		// config — это структура, содержащая все настройки подключения к PostgreSQL:
		// хост, порт, пользователь, пароль, база данных, параметры SSL, ограничения пула и т.д.

		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil {
			if pingErr := pool.Ping(ctx); pingErr == nil {
				// ping — проверяет подключение к PostgreSQL
				log.Println(" Успешное подключение к PostgreSQL")
				return &Storage{
					Pool: pool,
				}, nil
			} else {
				err = pingErr
			}
		}

		// Если не удалось — логируем и ждём
		log.Printf("Не удалось подключиться: %v", err)
		if attempt < maxRetries {
			log.Printf("Ждём %v перед следующей попыткой...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	// Все попытки провалились
	return nil, fmt.Errorf("не удалось подключиться после %d попыток: %w", maxRetries, err)
}
