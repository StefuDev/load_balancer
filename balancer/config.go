package balancer

import (
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert"`
	KeyFile  string `yaml:"key"`
}

type Config struct {
	IP          string    `yaml:"ip"`
	Port        int       `yaml:"port"`
	Server_list []string  `yaml:"server_list"`
	Balancer    string    `yaml:"balancer"`
	TLS         TLSConfig `yaml:"tls"`
}

func (config *Config) ReadFromYAML(path string) {
	buff, err := os.ReadFile(path)

	if err != nil {
		log.Fatalf("Error reading config file: %s\n", err)
		return
	}

	if err = yaml.Unmarshal(buff, config); err != nil {
		log.Fatalf("Error reading yaml: %s\n", err)
		return
	}
}
