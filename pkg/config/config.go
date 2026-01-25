package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Tailscale TailscaleConfig `yaml:"tailscale"`
	Logging   LoggingConfig   `yaml:"logging"`
	Server    ServerConfig    `yaml:"server"`
}

type TailscaleConfig struct {
	BinaryPath string `yaml:"binary_path"`
	DaemonPath string `yaml:"daemon_path"`
	SocketPath string `yaml:"socket_path"`
	AutoStart  bool   `yaml:"auto_start"`
}

type LoggingConfig struct {
	AppLogPath    string `yaml:"app_log_path"`
	DaemonLogPath string `yaml:"daemon_log_path"`
	Level         string `yaml:"level"`
	Format        string `yaml:"format"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
	Mode string `yaml:"mode"`
}

var (
	cfg        *Config
	configPath string
)

func Load(path string) (*Config, error) {
	configPath = path
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg = &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func Save() error {
	if cfg == nil || configPath == "" {
		return nil
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func Get() *Config {
	return cfg
}

func SetAutoStart(enabled bool) error {
	if cfg == nil {
		return nil
	}

	cfg.Tailscale.AutoStart = enabled
	return Save()
}
