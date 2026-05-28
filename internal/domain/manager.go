package domain

// QueueManager is the port for the collection of queues.
type QueueManager interface {
	CreateQueue(meta QueueMetadata) error
	GetQueue(id string) (Queue, error)
	GetMetadata(id string) (QueueMetadata, error)
	ListQueues() []QueueMetadata
}
