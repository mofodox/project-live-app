package models

import "gorm.io/gorm"

type Business struct {
	gorm.Model
	Name string `json:"name"`
}
