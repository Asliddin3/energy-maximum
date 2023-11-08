package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	Host string
	Port int

	DB_HOST string
	DB_PORT int
	DB_NAME string
	DB_PASS string
	DB_USER string
}

func Load() {
	c := Config{}
	c.Host = cast.ToString(getOrReturnDefault("HOST", "localhost"))
	c.Port = cast.ToInt(getOrReturnDefault("PORT", 8000))
	c.DB_HOST = cast.ToString(getOrReturnDefault("DB_HOST", "localhost"))
	c.DB_NAME = cast.ToString(getOrReturnDefault("DB_NAME", "market"))
	c.DB_PASS = cast.ToString(getOrReturnDefault("DB_PASS", "admin"))
	c.DB_USER = cast.ToString(getOrReturnDefault("DB_USER", "postgres"))
	c.DB_PORT = cast.ToInt(getOrReturnDefault("DB_PORT", 5432))
}
func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	err := godotenv.Load(".env")
	if err != nil {
		// log.Fatalf("Error loading .env file")
		return defaultValue
	}
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue
}
