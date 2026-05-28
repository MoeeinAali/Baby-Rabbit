package repository

import (
	"sync"

	"Baby-Rabbit/internal/domain"
)

// QueueManager is an in-memory implementation of domain.QueueManager.
type QueueManager struct {
	mu        sync.RWMutex
	queues    map[string]domain.Queue
	metadata  map[string]domain.QueueMetadata
	nameIndex map[string]string
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		queues:    make(map[string]domain.Queue),
		metadata:  make(map[string]domain.QueueMetadata),
		nameIndex: make(map[string]string),
	}
}

func (m *QueueManager) CreateQueue(meta domain.QueueMetadata) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.nameIndex[meta.Name]; exists {
		return domain.ErrQueueAlreadyExists
	}
	m.queues[meta.ID] = NewRingBufferQueue(meta.Capacity)
	m.metadata[meta.ID] = meta
	m.nameIndex[meta.Name] = meta.ID
	return nil
}

func (m *QueueManager) GetQueue(id string) (domain.Queue, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	q, ok := m.queues[id]
	if !ok {
		return nil, domain.ErrQueueNotFound
	}
	return q, nil
}

func (m *QueueManager) GetMetadata(id string) (domain.QueueMetadata, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	meta, ok := m.metadata[id]
	if !ok {
		return domain.QueueMetadata{}, domain.ErrQueueNotFound
	}
	return meta, nil
}

func (m *QueueManager) ListQueues() []domain.QueueMetadata {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]domain.QueueMetadata, 0, len(m.metadata))
	for _, meta := range m.metadata {
		result = append(result, meta)
	}
	return result
}
