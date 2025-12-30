package manager

import (
	"sync"
	"sync/atomic"
)

type ServerManager struct {
	mu           sync.RWMutex
	allServers   []string
	aliveServers atomic.Value
}

func (sm *ServerManager) GetServers() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	cpy := make([]string, len(sm.allServers))
	copy(cpy, sm.allServers)
	return cpy
}

func (sm *ServerManager) SetServers(list []string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	cpy := make([]string, len(list))
	copy(cpy, list)
	sm.allServers = cpy
}

func (sm *ServerManager) AddServer(server string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.allServers = append(sm.allServers, server)
}

func (sm *ServerManager) GetAliveServers() []string {
	val := sm.aliveServers.Load()
	if val == nil {
		return []string{}
	}
	list := val.([]string)
	cpy := make([]string, len(list))
	copy(cpy, list)
	return cpy
}

func (sm *ServerManager) SetAliveServers(list []string) {
	cpy := make([]string, len(list))
	copy(cpy, list)
	sm.aliveServers.Store(cpy)
}
