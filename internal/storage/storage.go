package storage

import "time"

type TaskStatus string

const (
	StatusNew  TaskStatus = "new"
	StatusDone TaskStatus = "done"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Deadline    time.Time `json:"deadline"`
}
