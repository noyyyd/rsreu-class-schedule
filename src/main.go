package main

import (
	"flag"
	"path/filepath"
	"rsreu-class-schedule/config"
	"rsreu-class-schedule/server"
)

const (
	flagConfig = "config"
)

var (
	defaultConfigPath = filepath.Join("config", "config.json")
)

func main() {
	var configPath string

	flag.StringVar(&configPath, flagConfig, defaultConfigPath, "path to JSON configPath. Default path is `config/config.json`")
	flag.Parse()

	server.New(config.NewConfig(configPath)).Run()
}
