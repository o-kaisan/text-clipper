package model

import (
	"time"

	"gorm.io/gorm"
)

type Clip struct {
	gorm.Model
	Title      string
	Content    string
	IsActive   *bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LastUsedAt time.Time
}

func NewClip(title, content string, isActive *bool, createdAt time.Time, updatedAt time.Time, lastUsedAt time.Time) *Clip {
	return &Clip{
		Title:      title,
		Content:    content,
		IsActive:   isActive,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		LastUsedAt: lastUsedAt,
	}
}
