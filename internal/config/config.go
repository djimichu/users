package config

import (
	"fmt"
	"github.com/caarlos0/env"
	"time"
)

type Config struct {
	AppEnv string `env:"APP_ENV" envDefault:"local"`

	// HTTP
	HTTPPort string `env:"HTTP_PORT" envDefault:":9091"`

	// JWT
	JWTSecret   string        `env:"JWT_SECRET" envDefault:"fallback-secret-change-me"`
	JWTExpHours time.Duration `env:"JWT_EXP_HOURS" envDefault:"15m"`

	// Database
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     int    `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER" envDefault:"postgres"`
	DBPassword string `env:"DB_PASSWORD" envDefault:""`
	DBName     string `env:"DB_NAME" envDefault:"myapp"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}
