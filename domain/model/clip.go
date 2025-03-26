package model

import (
	"log"
	"strings"
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

func (c *Clip) TruncateContent(height int) string {
	lines := strings.Split(c.Content, "\n")

	log.Print(height)
	if height <= 0 {
		return c.Content
	}

	if len(lines) <= height {
		return c.Content
	}

	if len(lines) > height {
		return strings.Join(lines[:height], "\n") + "\n…"
	}

	return c.Content
}
