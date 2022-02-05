package monitoring

import (
	"sync"
)

type Monitor struct {
	mu  sync.Mutex
	Rqs int64
}

func (m *Monitor) IncRqs() {
	m.mu.Lock()
	m.Rqs++
	m.mu.Unlock()
}

func (m *Monitor) GetRqs() int64 {
	return m.Rqs
}
