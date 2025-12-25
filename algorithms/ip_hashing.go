package algorithms

import (
	"hash/fnv"
	"log"
	"net/http"
	"strings"
)

type IPHashing struct {
	selected int
}

func (m *IPHashing) GetServer(server_list []string, request *http.Request) string {
	ip := strings.Split(request.RemoteAddr, ":")[0]

	hasher := fnv.New32a()

	hasher.Write([]byte(ip))

	ipHash := hasher.Sum32()

	index := ipHash % uint32(len(server_list))

	server := server_list[index]

	log.Printf("%s -> %s\n", ip, server)

	return server
}
