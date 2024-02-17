package entity

const (
	StatusNotDone = "not_done"
	StatusDone    = "done"
)

type Task struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
