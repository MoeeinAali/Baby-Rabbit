package domain

type Queue interface {
	Push(Message) error
	Pop() (Message, error)
	Size() int
	Capacity() int
	RemoveExpired()
}
