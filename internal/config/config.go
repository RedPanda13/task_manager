package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App  AppConfig  `yaml:"app"`
	HTTP HTTPConfig `yaml:"http"`
	DB   DBConfig   `yaml:"db"`
}

type AppConfig struct {
	Env string `yaml:"env"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type DBConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Name            string        `yaml:"name"`
	SSLMode         string        `yaml:"sslmode"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
	PingTimeout     time.Duration `yaml:"ping_timeout"`
	AutoMigrate     bool          `yaml:"auto_migrate"`
}

func Load() (Config, error) {
	return LoadFromFile("configs/application.yaml")
}

func LoadFromFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse yaml: %w", err)
	}

	applyDefaults(&cfg)

	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.App.Env == "" {
		cfg.App.Env = "dev"
	}

	if cfg.HTTP.Port == 0 {
		cfg.HTTP.Port = 8080
	}

	if cfg.DB.Host == "" {
		cfg.DB.Host = "localhost"
	}
	if cfg.DB.Port == 0 {
		cfg.DB.Port = 5432
	}
	if cfg.DB.User == "" {
		cfg.DB.User = "postgres"
	}
	if cfg.DB.Password == "" {
		cfg.DB.Password = "postgres"
	}
	if cfg.DB.Name == "" {
		cfg.DB.Name = "app"
	}
	if cfg.DB.SSLMode == "" {
		cfg.DB.SSLMode = "disable"
	}
	if cfg.DB.MaxOpenConns == 0 {
		cfg.DB.MaxOpenConns = 25
	}
	if cfg.DB.MaxIdleConns == 0 {
		cfg.DB.MaxIdleConns = 10
	}
	if cfg.DB.ConnMaxLifetime == 0 {
		cfg.DB.ConnMaxLifetime = 30 * time.Minute
	}
	if cfg.DB.ConnMaxIdleTime == 0 {
		cfg.DB.ConnMaxIdleTime = 10 * time.Minute
	}
	if cfg.DB.PingTimeout == 0 {
		cfg.DB.PingTimeout = 5 * time.Second
	}
}

func validate(cfg Config) error {
	if cfg.DB.Host == "" || cfg.DB.User == "" || cfg.DB.Name == "" {
		return fmt.Errorf("database configuration is incomplete")
	}

	return nil
}
