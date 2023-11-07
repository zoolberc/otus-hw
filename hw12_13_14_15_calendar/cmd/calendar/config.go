package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	sqlstorage "github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

type Config struct {
	LogLevel                string `yaml:"logLevel" env-default:"debug"`
	StorageType             string `yaml:"storageType" env-default:"memory"`
	HTTPServer              `yaml:"httpServer"`
	sqlstorage.DataBaseConf `yaml:"dataBase"`
}

type HTTPServer struct {
	Host    string        `yaml:"host" env-default:"localhost"`
	Port    string        `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
}

func NewConfig() (Config, error) {
	if configFile == "" {
		return Config{}, errors.New("config file is`t set")
	}
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("config file does`t exist: %s", configFile)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configFile, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
