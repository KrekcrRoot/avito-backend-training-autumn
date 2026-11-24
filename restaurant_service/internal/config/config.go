package config

import "os"

type Config struct {
	AppHost    string
	AppPort    string
	ApiKey     string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	KitchenURL string
}

func getEnv(key, placeholder string) string {

	value := os.Getenv(key)

	if value == "" {
		return placeholder
	}

	return value

}

func Load() *Config {
	return &Config{
		AppHost:    getEnv("APP_HOST", "localhost"),
		AppPort:    getEnv("APP_PORT", "3000"),
		ApiKey:     getEnv("API_KEY", "secret"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "restaurant"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		KitchenURL: getEnv("KITCHEN_URL", "http://localhost:3001"),
	}
}
