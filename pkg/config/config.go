package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Tailscale TailscaleConfig `yaml:"tailscale"`
	Logging   LoggingConfig   `yaml:"logging"`
	Server    ServerConfig    `yaml:"server"`
	Proxy     ProxyConfig     `yaml:"proxy"`
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

type ProxyConfig struct {
	Mode         string `yaml:"mode"`
	Port         int    `yaml:"port,omitempty"`
	Ports        []int  `yaml:"ports,omitempty"`
	ScanInterval int    `yaml:"scan_interval"`
	ExcludePorts []int  `yaml:"exclude_ports,omitempty"`
}

var (
	cfg *Config
)

func Load(path string) (*Config, error) {
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

func Get() *Config {
	return cfg
}
