package dto

import "time"

type SystemTodoCreateRequest struct {
	SystemName  string     `json:"systemName"`
	Title       string     `json:"title"`
	Detail      string     `json:"detail"`
	Status      uint       `json:"status"`
	Deadline    *time.Time `json:"deadline"`
	Urgency     uint       `json:"urgency"`
	CreatedName string     `json:"createdName"`
}

type SystemTodoUpdateRequest struct {
	SystemName  string     `json:"systemName"`
	Title       string     `json:"title"`
	Detail      string     `json:"detail"`
	Status      uint       `json:"status"`
	Deadline    *time.Time `json:"deadline"`
	Urgency     uint       `json:"urgency"`
	UpdatedName string     `json:"updatedName"`
}

type QuickSystemTodoUpdateRequest struct {
	Status      uint      `json:"status"`
	UpdatedAt   time.Time `json:"updatedAt"`
	UpdatedName string    `json:"updatedName"`
}
