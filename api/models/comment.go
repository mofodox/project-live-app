package models

import (
	"time"
)

type Comment struct {
	//	gorm.Model
	CommentID  uint32 `gorm:"primary_key;auto_increment" json:"id"`
	BusinessID uint32
	UserID     uint32
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
