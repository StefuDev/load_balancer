package balancer

import (
	"io"
	"load_balancer/algorithms"
	"log"
	"net/http"
)

type RequestHandler struct {
	client           *http.Client
	config           *Config
	balancing_method BalancingMethod
}

func NewRequestHandler(config *Config) *RequestHandler {
	handler := RequestHandler{
		client: &http.Client{},
		config: config,
	}

	switch config.Balancer {
	case "round_robin":
		handler.balancing_method = &algorithms.RoundRobin{}
	case "ip_hashing":
		handler.balancing_method = &algorithms.IPHashing{}
	default:
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
		log.Fatalln("Error requesting to the server")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	server_req.Header = request.Header

	server_res, err := handler.client.Do(server_req)

	if err != nil {
		log.Fatalln(err)
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
		log.Fatalln("Error copying the server Body!")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
