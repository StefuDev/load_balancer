package algorithms

import (
	"net/http"
	"sync/atomic"
)

type RoundRobin struct {
	counter int32
}

func (m *RoundRobin) GetServer(server_list []string, request *http.Request) string {
	next := atomic.AddInt32(&m.counter, 1)
	index := int(next-1) % len(server_list)
	return server_list[index]
}
