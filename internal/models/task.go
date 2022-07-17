package models

import "time"

type Task struct {
	ID         int       `json:"id" example:"1"`
	ObjectID   int       `json:"object_id" example:"1234"`
	IsDone     bool      `json:"is_done" example:"false"`
	CreatedAt  time.Time `json:"created_at" example:"2021-05-25T00:53:16.535668Z" format:"date-time"`
	FinishedAt time.Time `json:"ended_at" example:"2021-05-25T00:53:16.535668Z" format:"date-time"`
}
