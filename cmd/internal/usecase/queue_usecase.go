package usecase

import (
	"Baby-Rabbit/cmd/internal/domain"
	"time"

	"github.com/google/uuid"
)

type QueueUseCase struct {
	manager domain.QueueManager
}

func NewQueueUseCase(m domain.QueueManager) *QueueUseCase {
	return &QueueUseCase{manager: m}
}

func (u *QueueUseCase) CreateQueue(name string, capacity int) error {
	return u.manager.CreateQueue(name, capacity)
}

func (u *QueueUseCase) Push(queue string, value string, ttl int) error {
	q, err := u.manager.GetQueue(queue)
	if err != nil {
		return err
	}

	msg := domain.Message{
		ID:        uuid.New().String(),
		Value:     value,
		CreatedAt: time.Now(),
		TTL:       time.Duration(ttl) * time.Second,
	}

	return q.Push(msg)
}

func (u *QueueUseCase) Pop(queue string) (domain.Message, error) {
	q, err := u.manager.GetQueue(queue)
	if err != nil {
		return domain.Message{}, err
	}

	return q.Pop()
}

func (u *QueueUseCase) ListQueues() []string {
	return u.manager.ListQueues()
}
