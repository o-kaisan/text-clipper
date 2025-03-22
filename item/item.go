package item

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Title      string
	Content    string
	IsActive   *bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LastUsedAt time.Time
}

func NewItem(title, content string, isActive *bool, createdAt time.Time, updatedAt time.Time, lastUsedAt time.Time) *Item {
	return &Item{
		Title:      title,
		Content:    content,
		IsActive:   isActive,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		LastUsedAt: lastUsedAt,
	}
}
