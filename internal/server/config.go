package server

import "awesomeProject1/internal/store"

type Config struct {
	BindAddr string `yaml:"bind_addr"`
	LogLevel string `yaml:"log_level"`
	LogDir   string `yaml:"log_dir"`
	Storage  *store.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		Storage:  store.NewConfig(),
	}
}
