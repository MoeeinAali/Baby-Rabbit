package domain

type QueueManager interface {
	CreateQueue(name string, capacity int) error
	GetQueue(name string) (Queue, error)
	ListQueues() []string
}
