package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

const path = "config/config.yaml"

type Config struct {
	App struct {
		Env string `yaml:"env"`
	} `yaml:"app"`

	JWT struct {
		Secret   string `yaml:"secret"`
		ExpHours int    `yaml:"exp_hours"`
	} `yaml:"jwt"`

	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
}

func LoadConfig(path string) (*Config, error) {

	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
	)
}
