package repository

import (
	"Baby-Rabbit/cmd/internal/domain"
	"Baby-Rabbit/cmd/internal/pkg/logger"
	"errors"
	"sync"
)

type RingBufferQueue struct {
	buffer   []domain.Message
	head     int
	tail     int
	size     int
	capacity int

	mutex sync.Mutex
	cond  *sync.Cond
}

func NewRingBufferQueue(capacity int) *RingBufferQueue {
	q := &RingBufferQueue{
		buffer:   make([]domain.Message, capacity),
		capacity: capacity,
	}
	q.cond = sync.NewCond(&q.mutex)
	return q
}

func (q *RingBufferQueue) Push(msg domain.Message) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.size == q.capacity {
		logger.Log.Warn("Queue full, push rejected")
		return errors.New("queue full")
	}

	q.buffer[q.tail] = msg
	q.tail = (q.tail + 1) % q.capacity
	q.size++

	logger.Log.Debugf("Push message into queue: %s", msg.ID)
	q.cond.Signal()
	return nil
}

func (q *RingBufferQueue) Pop() (domain.Message, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for q.size == 0 {
		q.cond.Wait()
	}

	msg := q.buffer[q.head]
	q.head = (q.head + 1) % q.capacity
	q.size--

	logger.Log.Debugf("Pop message from queue: %s", msg.ID)
	return msg, nil
}

func (q *RingBufferQueue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.size
}

func (q *RingBufferQueue) Capacity() int {
	return q.capacity
}

func (q *RingBufferQueue) RemoveExpired() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	newBuffer := make([]domain.Message, q.capacity)
	index := 0

	for i := 0; i < q.size; i++ {
		pos := (q.head + i) % q.capacity
		msg := q.buffer[pos]
		if !msg.Expired() {
			newBuffer[index] = msg
			index++
		}
	}

	q.buffer = newBuffer
	q.head = 0
	q.tail = index
	q.size = index
}
