package manager

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type HealthCheck struct {
	serverManager  *ServerManager
	ticker         *time.Ticker
	maxConcurrency int
}

func NewHealthCheck(serverManager *ServerManager, interval time.Duration, maxConcurrency int) *HealthCheck {
	hc := &HealthCheck{
		serverManager:  serverManager,
		maxConcurrency: maxConcurrency,
	}

	if interval <= 0 {
		log.Println("health check interval must be > 0, setting to 5 minutes")
		interval = 5 * time.Minute
	}

	hc.ticker = time.NewTicker(interval)

	return hc
}

func (hc *HealthCheck) StartTimer() {
	go hc.tickerListener()
}

func (hc *HealthCheck) StopTimer() {
	hc.ticker.Stop()
}

func (hc *HealthCheck) tickerListener() {
	for range hc.ticker.C {
		hc.Update()
	}
}

func (hc *HealthCheck) Update() {
	log.Println("Running servers healthcheck")

	var wg sync.WaitGroup
	var mu sync.Mutex

	aliveList := make([]string, 0)

	sem := make(chan struct{}, hc.maxConcurrency)

	for _, server := range hc.serverManager.GetServers() {
		wg.Add(1)
		sem <- struct{}{}

		go func(s string) {
			defer wg.Done()
			defer func() { <-sem }()

			if checkServer(s) {
				mu.Lock()
				aliveList = append(aliveList, s)
				mu.Unlock()
			} else {
				log.Printf("Server '%s' is down\n", s)
			}
		}(server)
	}

	wg.Wait()

	hc.serverManager.SetAliveServers(aliveList)
}

func checkServer(server string) bool {
	resp, err := http.Get(server + "/health")
	if err != nil {
		log.Printf("Error cheking server '%s' with error: %s\n", server, err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
