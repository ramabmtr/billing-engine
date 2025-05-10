package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var (
	envOnce sync.Once
	env     *Env
)

type Env struct {
	Server   ServerEnv
	Database DatabaseEnv
}

type ServerEnv struct {
	Port   int
	ApiKey string
}

type DatabaseEnv struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (c *DatabaseEnv) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func InitEnv() {
	envOnce.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env file not found, using default environment variables")
		}

		env = &Env{
			Server: ServerEnv{
				Port:   getAsInt("SERVER_PORT", 8080),
				ApiKey: get("SERVER_API_KEY", ""),
			},
			Database: DatabaseEnv{
				Host:     get("DB_HOST", "localhost"),
				Port:     getAsInt("DB_PORT", 5432),
				User:     get("DB_USER", "admin"),
				Password: get("DB_PASSWORD", "admin"),
				DBName:   get("DB_NAME", "billing_engine"),
				SSLMode:  get("DB_SSLMODE", "disable"),
			},
		}
	})
}

func GetEnv() *Env {
	if env == nil {
		InitEnv()
	}

	return env
}

func get(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
