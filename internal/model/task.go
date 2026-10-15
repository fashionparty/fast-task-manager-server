package model

import (
	"time"
)

type Priority string

const (
	PriorityLow    Priority = "priority_low"
	PriorityMedium Priority = "priority_medium"
	PriorityHigh   Priority = "priority_high"
)

type Task struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	Title     string     `gorm:"not null" json:"title"`
	Content   string     `gorm:"type:text" json:"content"`
	DueDate   *time.Time `json:"dueDate"`
	Priority  Priority   `json:"priority"`
	IsDone    bool       `gorm:"type:bool" json:"isDone"`
	OwnerID   uint       `gorm:"index" json:"ownerId"`
}
