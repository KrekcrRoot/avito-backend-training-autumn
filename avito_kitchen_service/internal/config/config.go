package config

import "os"

type Config struct {
	AppHost    string
	AppPort    string
	ApiKey     string
	JWTSecret  string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
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
		AppPort:    getEnv("APP_PORT", "3001"),
		JWTSecret:  getEnv("JWT_SECRET", "some-secret"),
		ApiKey:     getEnv("API_KEY", ""),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "avito_kitchen"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
	}
}
