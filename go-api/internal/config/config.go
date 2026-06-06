// Package config centraliza la carga y exposición de variables de entorno.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config contiene todos los parámetros de configuración de la aplicación.
type Config struct {
	AppPort    string
	AppEnv     string
	NodeAPIURL string
}

// Load lee el archivo .env (si existe) y devuelve un Config poblado.
// Los valores ausentes se sustituyen por defaults razonables.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[config] .env file not found, using environment variables")
	}

	return &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),
		NodeAPIURL: getEnv("NODE_API_URL", "http://localhost:3000/api/v1/statistics"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
