package models

import "time"

type Business struct {
	ID        uint32    `gorm:"PRIMARY_KEY;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;UNIQUE;not null" json:"name"`
	Address   string    `gorm:"size:255;not null" json:"address"`
	UnitNo    string    `json:"unitNo"`
	Zipcode   string    `gorm:"size:255;not null" json:"zipcode"`
	Lat       float64   `gorm:"default:0" json:"lat"`
	Lng       float64   `gorm:"default:0" json:"lng"`
	Status    string    `gorm:"default:'active'" json:"status"` // active, 0 = inactive
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
