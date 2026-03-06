package domain

type QueueManager interface {
	CreateQueue(name string, capacity int) (string, error)
	GetQueue(id string) (Queue, error)
	GetQueueByName(name string) (Queue, error)
	ListQueues() []QueueMetadata
}
