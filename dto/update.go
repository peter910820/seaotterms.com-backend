package dto

import "time"

type SystemTodoUpdate struct {
	SystemName  string     `json:"systemName"`
	Title       string     `json:"title"`
	Detail      string     `json:"detail"`
	Status      uint       `json:"status"`
	Deadline    *time.Time `json:"deadline"`
	Urgency     uint       `json:"urgency"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	UpdatedName string     `json:"updatedName"`
}
