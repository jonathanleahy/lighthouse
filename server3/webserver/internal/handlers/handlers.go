package handlers

import (
	"github.com/jonathanleahy/project/jobscheduler"
)

// APIHandler wraps all API handlers
type APIHandler struct {
	jobsHandler  *JobsHandler
	statsHandler *StatsHandler
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(scheduler *jobscheduler.Scheduler) *APIHandler {
	return &APIHandler{
		jobsHandler:  NewJobsHandler(scheduler),
		statsHandler: NewStatsHandler(scheduler),
	}
}

// JobsHandler returns the jobs handler
func (h *APIHandler) JobsHandler() *JobsHandler {
	return h.jobsHandler
}

// StatsHandler returns the stats handler
func (h *APIHandler) StatsHandler() *StatsHandler {
	return h.statsHandler
}
