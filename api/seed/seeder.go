package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/mofodox/project-live-app/api/models"
)

var users = []models.User{
	models.User{
		Fullname: "John Doe",
		Email:    "johndoe@gmail.com",
		Password: "password",
	},
	models.User{
		Fullname: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

var businesses = []models.Business{
	models.Business{
		Name:    "Shake Shack Orchard Road",
		Address: "541 Orchard Rd, Liat Towers",
		UnitNo:  "#01-01",
		Zipcode: "238881",
		Lat:     0,
		Lng:     0,
	},
}

var categories = []models.Category{
	models.Category{
		Name:        "Baking",
		Description: "A method of preparing food that uses dry heat, typically in an oven, but can also be done in hot ashes, or on hot stones",
	},
	models.Category{
		Name:        "Private tutoring",
		Description: "A tutor as a person who gives individual, or in some cases small group",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.User{}, &models.Business{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v\n", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Business{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v\n", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table %v\n", err)
		}
	}

	for i := range businesses {

		businesses[i].Geocode()

		err = db.Debug().Model(&models.Business{}).Create(&businesses[i]).Error
		if err != nil {
			log.Fatalf("cannot seed business table %v\n", err)
		}
	}

	err2 := db.Debug().DropTableIfExists(&models.Category{}).Error
	if err2 != nil {
		log.Fatalf("cannot drop table: %v\n", err)
	}

	err = db.Debug().AutoMigrate(&models.Category{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v\n", err)
	}

	for i := range categories {
		err = db.Debug().Model(&models.Category{}).Create(&categories[i]).Error
		if err != nil {
			log.Fatalf("cannot seed categories table %v\n", err)
		}
	}
}
