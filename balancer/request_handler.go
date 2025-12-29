package balancer

import (
	"io"
	"load_balancer/algorithms"
	"log"
	"net"
	"net/http"
	"time"
)

type RequestHandler struct {
	client           *http.Client
	config           *Config
	balancing_method BalancingMethod
}

func NewRequestHandler(config *Config) *RequestHandler {
	handler := RequestHandler{
		client: &http.Client{
			Timeout: time.Duration(config.UpstreamTimeout) * time.Second,
		},
		config: config,
	}

	switch config.Balancer {
	case "round_robin":
		handler.balancing_method = &algorithms.RoundRobin{}
	case "ip_hashing":
		handler.balancing_method = &algorithms.IPHashing{}
	case "first":
		handler.balancing_method = &algorithms.First{}
	default:
		log.Printf("Balancing method '%s' not found, using 'First' method\n", config.Balancer)
		handler.balancing_method = &algorithms.First{}
	}

	return &handler
}

func (handler *RequestHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	selected_server := handler.balancing_method.GetServer(handler.config.Server_list, request)
	new_url := selected_server + request.URL.Path

	log.Printf("(%s) %s -> %s\n", request.RemoteAddr, request.URL, new_url)

	server_req, err := http.NewRequest(request.Method, new_url, request.Body)

	if err != nil {
		log.Println("Error requesting to the server")
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	server_req.Header = request.Header.Clone()

	if err := addXForwardedFor(request.RemoteAddr, server_req); err != nil {
		log.Println(err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	server_res, err := handler.client.Do(server_req)
	if err != nil {
		log.Println(err)
		http.Error(response, "Bad Gateway", http.StatusBadGateway)
		return
	}

	defer server_res.Body.Close()

	for key, values := range server_res.Header {
		for _, value := range values {
			response.Header().Add(key, value)
		}
	}

	response.WriteHeader(server_res.StatusCode)

	_, err = io.Copy(response, server_res.Body)
	if err != nil {
		log.Println("Error copying the server Body!")
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func addXForwardedFor(remoteAddr string, req *http.Request) error {
	clientIP, _, err := net.SplitHostPort(remoteAddr)

	if err != nil {
		return err
	}

	xff := req.Header.Get("X-Forwarded-For")

	if xff == "" {
		xff = clientIP
	} else {
		xff = xff + ", " + clientIP
	}

	req.Header.Set("X-Forwarded-For", xff)

	return nil
}
