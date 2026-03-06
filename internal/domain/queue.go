package domain

type Queue interface {
	Push(Message) error
	Pop() (Message, error)
	Size() int
	Capacity() int
	RemoveExpired()
}

type QueueMetadata struct {
	ID       string
	Name     string
	Capacity int
}
