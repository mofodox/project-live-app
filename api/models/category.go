package models

import (
	"time"
)

type Category struct {
	//gorm.Model
	ID uint `gorm:"primary_key;not null;"`
	//ParentID    uint   `gorm:"primary_key;not null;auto_increment:false"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
