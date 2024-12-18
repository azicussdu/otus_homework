package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	ServerConf   ServerConf   `yaml:"server"`
	DatabaseConf DatabaseConf `yaml:"database"`
	Logger       LoggerConf   `yaml:"logging"`
	StorageType  string       `yaml:"storage_type"`
}

type ServerConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DatabaseConf struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	Type          string `yaml:"type"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	DBName        string `yaml:"db_name"`
	MigrationPath string `yaml:"migration_path"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

func NewConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Printf("failed to close config file: %v\n", err)
		}
	}(file)

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// Validate the configuration
	if err = validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	if cfg.ServerConf.Host == "" {
		return errors.New("server host is required")
	}
	if cfg.ServerConf.Port <= 0 || cfg.ServerConf.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if cfg.DatabaseConf.Host == "" {
		return errors.New("database host is required")
	}
	if cfg.Logger.Level == "" {
		return errors.New("logger level is required")
	}
	if cfg.StorageType != "memory" && cfg.StorageType != "database" {
		return errors.New("storage type must be either 'memory' or 'database'")
	}
	return nil
}
