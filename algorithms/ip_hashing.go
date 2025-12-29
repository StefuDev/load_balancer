package algorithms

import (
	"hash/fnv"
	"log"
	"net"
	"net/http"
)

type IPHashing struct {
}

func (m *IPHashing) GetServer(server_list []string, request *http.Request) string {
	ip := request.Header.Get("X-Forwarded-For")

	if ip == "" {
		host, _, err := net.SplitHostPort(request.RemoteAddr)

		if err != nil {
			log.Println(err)
			return server_list[0]
		}

		ip = host
	}

	hasher := fnv.New32a()

	hasher.Write([]byte(ip))

	ipHash := hasher.Sum32()

	index := ipHash % uint32(len(server_list))

	server := server_list[index]

	return server
}
