package models

import (
	//"time"

	"github.com/jinzhu/gorm"
)

type Category struct {
	gorm.Model
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// createdAt time.Time `json:"createdat"`
	// updatedAt time.Time `json:"updatedat"`
	// parentID int `json:"parentid"`
}
