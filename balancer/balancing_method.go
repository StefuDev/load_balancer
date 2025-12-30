package balancer

import (
	"net/http"
)

type BalancingMethod interface {
	GetServer(serverList []string, request *http.Request) string
}
