package usecase

import (
	domain2 "Baby-Rabbit/internal/domain"
	"time"

	"github.com/google/uuid"
)

type QueueUseCase struct {
	manager domain2.QueueManager
}

func NewQueueUseCase(m domain2.QueueManager) *QueueUseCase {
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

	msg := domain2.Message{
		ID:        uuid.New().String(),
		Value:     value,
		CreatedAt: time.Now(),
		TTL:       time.Duration(ttl) * time.Second,
	}

	return q.Push(msg)
}

func (u *QueueUseCase) Pop(queue string) (domain2.Message, error) {
	q, err := u.manager.GetQueue(queue)
	if err != nil {
		return domain2.Message{}, err
	}

	return q.Pop()
}

func (u *QueueUseCase) ListQueues() []string {
	return u.manager.ListQueues()
}
