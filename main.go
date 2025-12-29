package main

import (
	"fmt"
	"load_balancer/balancer"
	"log"
	"net/http"
	"time"
)

func main() {
	config := balancer.Config{}

	if err := config.ReadFromYAML("config.yaml"); err != nil {
		log.Fatalln(err)
	}

	if err := config.Validate(); err != nil {
		log.Fatalln(err)
	}

	handler := balancer.NewRequestHandler(&config)

	addr := fmt.Sprintf("%s:%d", config.IP, config.Port)

	server := &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if config.TLS.Enabled {
		log.Printf("Listening on (HTTPS) https://%s\n", addr)
		log.Fatal(server.ListenAndServeTLS(config.TLS.CertFile, config.TLS.KeyFile))
	} else {
		log.Printf("Listening on (HTTP) http://%s\n", addr)
		log.Fatal(server.ListenAndServe())
	}
}
