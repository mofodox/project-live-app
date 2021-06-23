package models

import (
	"time"
)

type File struct {
	ID          uint32    `gorm:"PRIMARY_KEY;auto_increment" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"size:255;not null" json:"description"`
	URL         string    `gorm:"size:255;not null" json:"url"`
	Type        string    `json:"type"`                           // type: user == avatar, business == avatar/logo, comment = images that accompanies comments
	Status      string    `gorm:"default:'active'" json:"status"` // active / inactive (deleted)
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
