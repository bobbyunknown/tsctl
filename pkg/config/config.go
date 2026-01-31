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
	StateDir  string `yaml:"state_dir"`
	Hostname  string `yaml:"hostname"`
	AuthKey   string `yaml:"auth_key"`
	Ephemeral bool   `yaml:"ephemeral"`
}

type LoggingConfig struct {
	AppLogPath string `yaml:"app_log_path"`
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
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
