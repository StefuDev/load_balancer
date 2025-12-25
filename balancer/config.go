package balancer

import (
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	IP          string   `json:"ip"`
	Port        int      `json:"port"`
	Server_list []string `json:"server_list"`
	Balancer    string   `json:"balancer"`
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
