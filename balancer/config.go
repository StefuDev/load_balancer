package balancer

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert"`
	KeyFile  string `yaml:"key"`
}

type HealthCheckConfig struct {
	Enabled        bool `yaml:"enabled"`
	Interval       int  `yaml:"interval"`
	MaxConcurrency int  `yaml:"max_concurrency"`
}

type Config struct {
	IP              string            `yaml:"ip"`
	Port            int               `yaml:"port"`
	ServerList      []string          `yaml:"server_list"`
	Balancer        string            `yaml:"balancer"`
	ReadTimeout     int               `yaml:"read_timeout"`
	WriteTimeout    int               `yaml:"write_timeout"`
	UpstreamTimeout int               `yaml:"upstream_timeout"`
	TLS             TLSConfig         `yaml:"tls"`
	HealthCheck     HealthCheckConfig `yaml:"healthcheck"`
}

func (config *Config) ReadFromYAML(path string) error {
	buff, err := os.ReadFile(path)

	if err != nil {
		return fmt.Errorf("Error reading config file: %s", err)
	}

	if err = yaml.Unmarshal(buff, config); err != nil {
		return fmt.Errorf("Error reading yaml: %s", err)
	}

	return nil
}

func (config *Config) Validate() error {
	if len(config.ServerList) == 0 {
		return fmt.Errorf("server_list cannot be empty")
	}

	if config.IP == "" {
		return fmt.Errorf("invalid ip")
	}

	if config.Port <= 0 {
		return fmt.Errorf("invalid port")
	}

	if config.ReadTimeout <= 0 || config.WriteTimeout <= 0 || config.UpstreamTimeout <= 0 {
		return fmt.Errorf("timeouts must be > 0")
	}

	if config.TLS.Enabled {
		if config.TLS.CertFile == "" || config.TLS.KeyFile == "" {
			return fmt.Errorf("TLS enabled but cert/key missing")
		}
	}

	if config.HealthCheck.Enabled {
		if config.HealthCheck.Interval <= 0 {
			return fmt.Errorf("Health check interval must be > 0")
		}

		if config.HealthCheck.MaxConcurrency <= 0 {
			return fmt.Errorf("Health check max concurrency must be > 0")
		}
	}

	return nil
}
