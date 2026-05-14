package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUsername string
	RabbitMQPassword string

	NewDir     string
	ProcessDir string
	DoneDir    string
	ErrorDir   string
	ErrorLog   string

	DataBase string

	MovieSource  string
	SeriesSource string
}

func Load() Config {
	cfg := Config{
		RabbitMQHost:     getEnv("RABBITMQ_HOST", "localhost"),
		RabbitMQPort:     getEnv("RABBITMQ_PORT", "5672"),
		RabbitMQUsername: getEnv("RABBITMQ_USERNAME", "guest"),
		RabbitMQPassword: getEnv("RABBITMQ_PASSWORD", "guest"),

		NewDir:     getEnv("NEW_DIR", "/tmp/new"),
		ProcessDir: getEnv("PROCESS_DIR", "/tmp/process"),
		DoneDir:    getEnv("DONE_DIR", "/tmp/done"),
		ErrorDir:   getEnv("ERROR_DIR", "/tmp/error"),
		ErrorLog:   getEnv("ERROR_LOG", "/tmp/error/error.log"),


		MovieSource:  getEnv("MOVIE_SOURCE", "/tmp/movies"),
		SeriesSource: getEnv("SERIES_SOURCE", "/tmp/series"),
	}

	validate(cfg)

	return cfg
}

func (c Config) RabbitMQURL() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		c.RabbitMQUsername,
		c.RabbitMQPassword,
		c.RabbitMQHost,
		c.RabbitMQPort,
	)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		fmt.Printf("%s env is missing", key)
		return fallback
	}

	return value
}

func validate(cfg Config) {
	required := map[string]string{
		"RABBITMQ_HOST": cfg.RabbitMQHost,
		"RABBITMQ_PORT": cfg.RabbitMQPort,
		"NEW_DIR":       cfg.NewDir,
		"PROCESS_DIR":   cfg.ProcessDir,
		"DONE_DIR":      cfg.DoneDir,
		"ERROR_DIR":     cfg.ErrorDir,
	}

	for key, value := range required {
		if value == "" {
			log.Fatalf("%s is required", key)
		}
	}
}