package balancer

import (
	"load_balancer/algorithms"
	"load_balancer/manager"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type RequestHandler struct {
	client          *http.Client
	config          *Config
	balancingMethod BalancingMethod
	serverManager   *manager.ServerManager
}

func NewRequestHandler(config *Config) *RequestHandler {
	serverManager := &manager.ServerManager{}
	serverManager.SetServers(config.ServerList)

	if config.HealthCheck.Enabled {
		healthCheck := manager.NewHealthCheck(
			serverManager,
			time.Duration(config.HealthCheck.Interval)*time.Second,
			config.HealthCheck.MaxConcurrency)

		healthCheck.StartTimer()
		healthCheck.Update()
	}

	handler := RequestHandler{
		client: &http.Client{
			Timeout: time.Duration(config.UpstreamTimeout) * time.Second,
		},
		config:        config,
		serverManager: serverManager,
	}

	switch config.Balancer {
	case "round_robin":
		handler.balancingMethod = &algorithms.RoundRobin{}
	case "ip_hashing":
		handler.balancingMethod = &algorithms.IPHashing{}
	case "random":
		handler.balancingMethod = &algorithms.Random{}
	default:
		log.Printf("Balancing method '%s' not found, using 'Random' method\n", config.Balancer)
		handler.balancingMethod = &algorithms.Random{}
	}

	return &handler
}

func (handler *RequestHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	aliveServers := handler.serverManager.GetAliveServers()

	if len(aliveServers) <= 0 {
		log.Println("No alive servers")
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	targetServer := handler.balancingMethod.GetServer(aliveServers, req)
	targetURL, err := url.Parse(targetServer)

	if err != nil {
		log.Println("Error parsing target url: ", err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("(%s) %s -> %s\n", req.RemoteAddr, req.URL, targetURL)

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.ServeHTTP(rw, req)
}
