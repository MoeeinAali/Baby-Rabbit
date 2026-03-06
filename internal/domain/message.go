package domain

import "time"

type Message struct {
	ID        string
	Value     string
	CreatedAt time.Time
	TTL       time.Duration
}

func (m Message) Expired() bool {
	return time.Since(m.CreatedAt) > m.TTL
}
