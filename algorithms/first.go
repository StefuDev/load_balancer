package algorithms

import "net/http"

type First struct {
}

func (m *First) GetServer(server_list []string, request *http.Request) string {
	return server_list[0]
}
