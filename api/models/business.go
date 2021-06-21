package models

import "time"

type Business struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	UnitNo    string    `json:"unitNo"`
	Zipcode   string    `json:"zipcode"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
