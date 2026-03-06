package service

import (
	"Baby-Rabbit/internal/domain"
	"time"
)

func StartTTLCleaner(manager domain.QueueManager) {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {
			for _, name := range manager.ListQueues() {
				q, _ := manager.GetQueue(name)
				q.RemoveExpired()
			}
		}
	}()
}
