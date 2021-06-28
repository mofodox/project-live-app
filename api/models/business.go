package models

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"googlemaps.github.io/maps"
)

type Business struct {
	ID               uint32    `gorm:"PRIMARY_KEY;auto_increment" json:"id"`
	UserID           uint32    `gorm:"type:int(10);not null" json:"userId"` // last user that edited. We're only keeping last edit.
	Name             string    `gorm:"size:255;UNIQUE;not null" json:"name"`
	ShortDescription string    `gorm:"size:255;not null" json:"shortDescription"`
	Description      string    `gorm:"type:text;not null" json:"description"` // 65,535 characters, rich text editor
	Address          string    `gorm:"size:255;not null" json:"address"`
	UnitNo           string    `gorm:"size:255;not null" json:"unitNo"`
	Zipcode          string    `gorm:"size:255;not null" json:"zipcode"`
	Lat              float64   `gorm:"default:0" json:"lat"`
	Lng              float64   `gorm:"default:0" json:"lng"`
	Status           string    `gorm:"default:'active'" json:"status"` // active / inactive (deleted)
	Website          string    `gorm:"size:255;" json:"website"`
	Instagram        string    `gorm:"size:255;" json:"instagram"`
	Facebook         string    `gorm:"size:255;" json:"facebook"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (business *Business) Validate() error {
	if len(strings.TrimSpace(business.Name)) < 5 {
		return errors.New("enter a business name that has at least 5 characters")
	}

	if len(strings.TrimSpace(business.ShortDescription)) < 20 {
		return errors.New("enter a short business description that has at least 20 characters")
	}

	if len(strings.TrimSpace(business.Description)) < 20 {
		return errors.New("enter a business description that has at least 20 characters")
	}

	return nil
}

// https://pkg.go.dev/googlemaps.github.io/maps?utm_source=godoc
// address to lat / lng
func (business *Business) Geocode() (lat float64, lnt float64, err error) {

	var gMapsAPI = os.Getenv("GMapsAPI")
	c, err := maps.NewClient(maps.WithAPIKey(gMapsAPI))

	if err != nil {
		return 0, 0, errors.New("unable to fetch lat/lng")
	}

	fullAddress := business.Address + " " + business.UnitNo + " " + business.Zipcode

	// https://pkg.go.dev/googlemaps.github.io/maps?utm_source=godoc#GeocodingRequest
	r := &maps.GeocodingRequest{
		Address: fullAddress,
		Region:  "SG",
	}

	// https://pkg.go.dev/googlemaps.github.io/maps?utm_source=godoc#Client.Geocode
	// https://pkg.go.dev/googlemaps.github.io/maps?utm_source=godoc#GeocodingResult
	results, err := c.Geocode(context.Background(), r)

	if err == nil && len(results) > 0 {

		business.Lat = results[0].Geometry.Location.Lat
		business.Lng = results[0].Geometry.Location.Lng

		return business.Lat, business.Lng, nil
	}

	return 0, 0, errors.New("unable to fetch lat/lng")
}
