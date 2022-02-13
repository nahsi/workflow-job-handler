package main

import (
	"time"

	"github.com/go-playground/webhooks/v6/github"
)

// WorkflowJob holds fields from GitHub payload we care about
type WorkflowJob struct {
	ID          int64
	RunID       int64
	Status      string
	Conclusion  string
	StartedAt   time.Time
	CompletedAt time.Time
	Name        string
	Repository  string
	User        string
}

// RedisWorkflowJob implements a mapping between WorkflowJob and Redis hash
type RedisWorkflowJob struct {
	ID          int64     `redis:"id"`
	RunID       int64     `redis:"run_id"`
	Status      string    `redis:"status"`
	Conclusion  string    `redis:"conclusion"`
	StartedAt   time.Time `redis:"started_at"`
	CompletedAt time.Time `redis:"completed_at"`
	Name        string    `redis:"name"`
	Repository  string    `redis:"repository"`
	User        string    `redis:"user"`
}

// Adapter between github.WorkflowJobPayload and WorkflowJob
func intoWorkFlowJob(p *github.WorkflowJobPayload) WorkflowJob {
	return WorkflowJob{
		ID:          p.WorkflowJob.ID,
		RunID:       p.WorkflowJob.RunID,
		Status:      p.WorkflowJob.Status,
		Conclusion:  p.WorkflowJob.Conclusion,
		StartedAt:   p.WorkflowJob.StartedAt,
		CompletedAt: p.WorkflowJob.CompletedAt,
		Name:        p.WorkflowJob.Name,
		Repository:  p.Repository.FullName,
		User:        p.Sender.Login,
	}
}

// Adapter between WorkflowJob and RedisWorkflowJob
func intoRedisWorkFlowJob(p *WorkflowJob) RedisWorkflowJob {
	return RedisWorkflowJob{
		ID:          p.ID,
		RunID:       p.RunID,
		Status:      p.Status,
		Conclusion:  p.Conclusion,
		StartedAt:   p.StartedAt,
		CompletedAt: p.CompletedAt,
		Name:        p.Name,
		Repository:  p.Repository,
		User:        p.User,
	}
}
