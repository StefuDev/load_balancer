package algorithms

import "net/http"

type RoundRobin struct {
	selected int
}

func (m *RoundRobin) GetServer(server_list []string, request *http.Request) string {
	server := server_list[m.selected]

	m.selected = (m.selected + 1) % len(server_list)

	return server
}
