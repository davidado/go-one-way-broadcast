// Package config provides configuration for the application.
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Envs holds the configuration for the application.
var Envs = initConfig()

// Config struct
type Config struct {
	Listen        string
	RedisHost     string
	RedisPassword string
	RabbitMQHost  string
}

func initConfig() Config {
	godotenv.Load()

	return Config{
		Listen:        getEnv("LISTEN", ":8080"),
		RedisHost:     fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379")),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RabbitMQHost:  fmt.Sprintf("amqp://%s:%s@%s:%s/", getEnv("RABBITMQ_USER", "guest"), getEnv("RABBITMQ_PASSWORD", "guest"), getEnv("RABBITMQ_HOST", "localhost"), getEnv("RABBITMQ_PORT", "5672")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInSeconds(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
