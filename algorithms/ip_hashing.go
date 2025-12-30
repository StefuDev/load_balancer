package algorithms

import (
	"hash/fnv"
	"log"
	"net"
	"net/http"
)

type IPHashing struct {
}

func (alg *IPHashing) GetServer(serverList []string, request *http.Request) string {
	ip := request.Header.Get("X-Forwarded-For")

	if ip == "" {
		host, _, err := net.SplitHostPort(request.RemoteAddr)

		if err != nil {
			log.Println(err)
			return serverList[0]
		}

		ip = host
	}

	hasher := fnv.New32a()

	hasher.Write([]byte(ip))

	ipHash := hasher.Sum32()

	index := ipHash % uint32(len(serverList))

	server := serverList[index]

	return server
}
