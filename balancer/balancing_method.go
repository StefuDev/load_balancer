package balancer

import "net/http"

type BalancingMethod interface {
	GetServer(server_list []string, request *http.Request) string
}
