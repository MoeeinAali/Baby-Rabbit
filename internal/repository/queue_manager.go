package repository

import (
	"Baby-Rabbit/internal/domain"
	"Baby-Rabbit/internal/pkg/logger"
	"errors"
	"sync"
)

type QueueManager struct {
	queues map[string]domain.Queue
	mutex  sync.RWMutex
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		queues: make(map[string]domain.Queue),
	}
}

func (m *QueueManager) CreateQueue(name string, capacity int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.queues[name]; exists {
		logger.Log.Warnf("Queue already exists: %s", name)
		return errors.New("queue already exists")
	}

	m.queues[name] = NewRingBufferQueue(capacity)
	logger.Log.Infof("Created queue: %s, capacity: %d", name, capacity)
	return nil
}

func (m *QueueManager) GetQueue(name string) (domain.Queue, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	q, ok := m.queues[name]
	if !ok {
		return nil, errors.New("queue not found")
	}

	return q, nil
}

func (m *QueueManager) ListQueues() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	names := make([]string, 0, len(m.queues))
	for k := range m.queues {
		names = append(names, k)
	}
	return names
}
