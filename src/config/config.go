package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	SchedulesPath = "schedules"
)

//go:generate go-enum -f=$GOFILE
/*
ENUM(
	file
)
*/
type RepositoryType string

type Config struct {
	RepositoryType RepositoryType `json:"repository_type,omitempty"`
	SchedulesPath  string         `json:"schedules_path,omitempty"`
}

func NewConfig(configPath string) *Config {
	file, err := os.OpenFile(configPath, os.O_RDONLY, 0777)
	if err != nil {
		log.Fatalf("failed open cofig file %s: %v", configPath, err)
		return nil
	}

	configBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed read cofig file %s: %v", configPath, err)
		return nil
	}

	config := new(Config)

	err = json.Unmarshal(configBytes, config)
	if err != nil {
		log.Fatalf("failed unmarshal JSON from cofig file %s: %v", configPath, err)
		return nil
	}

	config.normalize()

	return config
}

func (c *Config) normalize() {
	if c.RepositoryType == "" {
		c.RepositoryType = RepositoryTypeFile
	}
	if c.SchedulesPath == "" {
		c.SchedulesPath = filepath.Join("..", SchedulesPath)
	}
}
