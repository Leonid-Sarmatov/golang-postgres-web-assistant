package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	EnvMode          string `yaml:"environment_mode" env-default:"local"`
	HTTPServerConfig `yaml:"http_server"`
}

type HTTPServerConfig struct {
	PagePath          string        `yaml:"page_path" env-default:"./index.html"`
	Address           string        `yaml:"address" env-default:"localhost:8080"`
	RequestTimeout    time.Duration `yaml:"request_timeout" env-default:"4s"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	// Проверяем задан ли путь до конфиг-файла
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is empty")
	}

	// Проверяем существует ли конфиг-файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file is not exist: %v", configPath)
	}

	// Считываем конфиги
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Reading config was failed: %v", err)
	}
	return &cfg
}
