package algorithms

import (
	"math/rand"
	"net/http"
	"time"
)

type Random struct {
}

func (alg *Random) GetServer(serverList []string, request *http.Request) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return serverList[r.Intn(len(serverList))]
}
