package repository

import (
	"Baby-Rabbit/internal/domain"
	"Baby-Rabbit/internal/pkg/logger"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type QueueManager struct {
	queues    map[string]domain.Queue         // key: UUID
	metadata  map[string]domain.QueueMetadata // key: UUID
	nameIndex map[string]string               // key: name, value: UUID
	mutex     sync.RWMutex
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		queues:    make(map[string]domain.Queue),
		metadata:  make(map[string]domain.QueueMetadata),
		nameIndex: make(map[string]string),
	}
}

func (m *QueueManager) CreateQueue(name string, capacity int) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.nameIndex[name]; exists {
		logger.Log.Warnf("Queue already exists: %s", name)
		return "", errors.New("queue already exists")
	}

	queueID := uuid.New().String()
	m.queues[queueID] = NewRingBufferQueue(capacity)
	m.metadata[queueID] = domain.QueueMetadata{
		ID:       queueID,
		Name:     name,
		Capacity: capacity,
	}
	m.nameIndex[name] = queueID
	logger.Log.Infof("Created queue: %s (ID: %s), capacity: %d", name, queueID, capacity)
	return queueID, nil
}

func (m *QueueManager) GetQueue(id string) (domain.Queue, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	q, ok := m.queues[id]
	if !ok {
		return nil, errors.New("queue not found")
	}

	return q, nil
}

func (m *QueueManager) GetQueueByName(name string) (domain.Queue, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	id, ok := m.nameIndex[name]
	if !ok {
		return nil, errors.New("queue not found")
	}

	return m.queues[id], nil
}

func (m *QueueManager) ListQueues() []domain.QueueMetadata {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make([]domain.QueueMetadata, 0, len(m.metadata))
	for _, meta := range m.metadata {
		result = append(result, meta)
	}
	return result
}
