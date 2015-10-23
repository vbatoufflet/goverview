package main

import (
	"github.com/vbatoufflet/goverview/pkg/poller"
)

// Config represents the service configuration struct
type Config struct {
	LogPath   string        `yaml:"log_path"`
	LogLevel  string        `yaml:"log_level"`
	BindAddr  string        `yaml:"bind_addr"`
	StaticDir string        `yaml:"static_dir"`
	Poller    poller.Config `yaml:"poller"`
}
