// Package config предоставляет загрузку конфигурации приложения
// из JSON-файла с учётом приоритетов источников.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

// Config хранит настройки приложения.
type Config struct {
	Port            string `json:"port"`
	DatabaseURL     string `json:"database_url"`
	JWTSecret       string `json:"jwt_secret"`
	ShutdownTimeout int    `json:"shutdown_timeout"`
	GinMode         string `json:"gin_mode"`
}

// Default возвращает конфигурацию со значениями по умолчанию.
// Используется как основа, если файл конфигурации не задан
// или задан частично.
func Default() Config {

	return Config{
		Port:            ":8080",
		DatabaseURL:     "postgres://user:password@localhost:5432/mydb?sslmode=disable",
		JWTSecret:       "change-me",
		ShutdownTimeout: 5,
		GinMode:         "debug",
	}
}

// Load читает конфигурацию из JSON-файла по указанному пути.
// Если файл не существует — возвращает дефолтную конфигурацию
// без ошибки. Если файл невалиден — возвращает ошибку.
func Load(path string) (Config, error) {

	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("не удалось прочитать конфиг: %w", err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("не удалось распарсить конфиг: %w", err)
	}
	return cfg, nil
}

// ShutdownDuration возвращает таймаут graceful shutdown
// в формате time.Duration.
// context.WithTimeout требует time.Duration, а не int.
func (c Config) ShutdownDuration() time.Duration {
	return time.Duration(c.ShutdownTimeout) * time.Second
}

var (
	configPathShort string
	configPathLong  string
)

func init() {
	// flag.StringVar - зарегистрируй флаг
	// flag.Parse() - прочитай флаги из os.Args
	flag.StringVar(&configPathShort, "c", "", "путь к файлу конфигурации")
	flag.StringVar(&configPathLong, "config", "", "путь к файлу конфигурации")
}

// GetConfigPath определяет путь к файлу конфигурации.
// Приоритет: флаг -c/-config → переменная окружения CONFIG →
// дефолтное значение "config.json".
func GetConfigPath() string {
	// flag.Parse() уже вызвали в main до этого (можно вызвать один раз за программу)
	// если flag.Parse() не вызвать — флаги не заполнятся
	// flag.StringVar только регистрирует, но не парсит
	if configPathShort != "" {
		return configPathShort
	}
	if configPathLong != "" {
		return configPathLong
	}
	if env := os.Getenv("CONFIG"); env != "" {
		return env
	}
	return "config.json"
}
