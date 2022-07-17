package models

import "time"

type Letter struct {
	ID         int       `json:"id" example:"1"`
	TaskID     int       `json:"task_id" example:"1234"`
	Email      string    `json:"email" example:"test@task.com"`
	IsAccepted bool      `json:"is_accepted" example:"false"`
	SentAt     time.Time `json:"sent_at" example:"2021-05-25T00:53:16.535668Z" format:"date-time"`
	FinishedAt time.Time `json:"ended_at" example:"2021-05-25T00:53:16.535668Z" format:"date-time"`
}
