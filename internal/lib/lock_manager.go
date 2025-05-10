package lib

import "sync"

type LockManager interface {
	GetLock(key string) *sync.Mutex
}

type lockManager struct {
	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

func NewLockManager() LockManager {
	return &lockManager{
		locks: make(map[string]*sync.Mutex),
	}
}

func (m *lockManager) GetLock(key string) *sync.Mutex {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.locks[key]; !ok {
		m.locks[key] = &sync.Mutex{}
	}
	return m.locks[key]
}
