package main

import (
	"sync"
	"time"
)

// QueueManager manages multiple work queues
type QueueManager struct {
	queues map[string]*WorkQueue
	config map[string]QueueConfig
	mu     sync.RWMutex
}

// WorkQueue represents a single queue
type WorkQueue struct {
	name          string
	items         []*QueuedCheck
	processing    map[string]bool
	maxSize       int
	maxConcurrent int
	mu            sync.RWMutex
}

// QueueStatus represents the current status of a queued job
type QueueStatus struct {
	Status     string    `json:"status"`
	Position   int       `json:"position"`
	CacheKey   string    `json:"cache_key"`
	QueueTime  time.Time `json:"queue_time"`
	StartTime  time.Time `json:"start_time"`
	StepsToRun []string  `json:"steps_to_run,omitempty"`
}

// Ensure QueueStatus implements ServiceResponse
func (qs *QueueStatus) ResponseType() string {
	return qs.Status
}

// Helper function to check if any element in needles exists in haystack
func containsAny(haystack []string, needles ...string) bool {
	for _, needle := range needles {
		for _, item := range haystack {
			if item == needle {
				return true
			}
		}
	}
	return false
}
