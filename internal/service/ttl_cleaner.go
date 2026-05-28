package service

import (
	"context"
	"time"

	"Baby-Rabbit/internal/domain"
	"Baby-Rabbit/internal/pkg/logger"
)

// TTLCleaner periodically asks every queue to drop expired messages.
// It depends only on the domain.QueueManager port.
type TTLCleaner struct {
	manager  domain.QueueManager
	interval time.Duration
}

func NewTTLCleaner(m domain.QueueManager, interval time.Duration) *TTLCleaner {
	if interval <= 0 {
		interval = time.Second
	}
	return &TTLCleaner{manager: m, interval: interval}
}

// Run blocks until ctx is canceled, periodically sweeping expired messages.
func (c *TTLCleaner) Run(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("ttl cleaner stopped")
			return
		case <-ticker.C:
			c.sweep()
		}
	}
}

func (c *TTLCleaner) sweep() {
	for _, meta := range c.manager.ListQueues() {
		q, err := c.manager.GetQueue(meta.ID)
		if err != nil {
			continue
		}
		if removed := q.RemoveExpired(); removed > 0 {
			logger.Log.Debugf("ttl cleaner removed %d expired messages from %s", removed, meta.Name)
		}
	}
}
