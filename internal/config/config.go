package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Connect struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		Database int    `yaml:"database"`
		Protocol int    `yaml:"protocol"`
	} `yaml:"connect"`
	Interval time.Duration `yaml:"interval"`
	Limit    int64         `yaml:"limit"`
}

func InitConfig(filename string) (*Config, error) {
	var cnf Config
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can't read file: %w", err)
	}

	err = yaml.Unmarshal(file, &cnf)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal file: %w", err)
	}
	return &cnf, nil
}
