package models

import (
	"time"
)

type Comment struct {
	//	gorm.Model
	CommentID  uint32    `gorm:"primary_key;auto_increment" json:"id"`
	BusinessID uint32    `gorm:"not_null" json:"business_id"`
	UserID     uint32    `gorm:"not_null" json:"user_id"`
	Content    string    `gorm:"size:255" json:"content"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
