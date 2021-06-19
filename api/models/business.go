package models

import "github.com/jinzhu/gorm"

type Business struct {
	gorm.Model
	Name string `json:"name"`
}
