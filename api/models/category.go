package models

import (
	"time"
)

type Category struct {
	//gorm.Model
	ID          uint32    `gorm:"primary_key;auto_increment;" json:"id"`
	ParentID    int32     `gorm:"size:255;default:1" json:"parentid"`
	Name        string    `gorm:"size:255;not_nulls" json:"name"`
	Description string    `gorm:"size:255;not_nulls" json:"description"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
