package algorithms

import (
	"net/http"
	"sync/atomic"
)

type RoundRobin struct {
	counter int32
}

func (alg *RoundRobin) GetServer(serverList []string, request *http.Request) string {
	next := atomic.AddInt32(&alg.counter, 1)
	index := int(next-1) % len(serverList)

	return serverList[index]
}
