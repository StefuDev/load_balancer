# Load Balancer

This project showcases a load balancer written in Go, with simple algorithms and config.

## ✨ Features

* **HTTPS support:** Built-in HTTPS support with TLS.
* **Multiple algorithms:** This project implements the Round Robin, IP Hashing and First server as default algorithms.
* **Multiple servers support:** You can add as many servers as you want.
* **Easy to contribute:** Written in a structured, minimal way that makes it easy to add new algorithms and features.
* **Config in YAML:** The configuration is written in YAML.

## ⚙️ Configuration Guide

### 1. Create an "config.yaml"
To make the configuration, you can use "example_config.yaml" as reference.

The current algorithms are:
-  round_robin
-  ip_hashing
-  first

example_config.yaml:
```yaml
ip: 127.0.0.1
port: 8080

server_list:
  - "https://server1.example.com"
  - "https://server2.example.com"
  - "https://server3.example.com"

balancer: "round_robin"

read_timeout: 10
write_timeout: 10

# Optional
tls:
  enabled: true
  cert: "cert.pem"
  key: "key.pem"
```

### 2. Setup servers
Ensure all servers is running and the selected port and IP is not already in use.

### 3. Run the load balancer
In the terminal run the project using ```go run load_balancer```
