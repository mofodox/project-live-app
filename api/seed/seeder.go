package seed

import (
	"log"
	//"project-live-app/api/models"

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
		Name:    "John Doe Pte Ltd",
		Address: "Some address",
		UnitNo:  "12345",
		Zipcode: "S123456",
		Lat:     0,
		Lng:     0,
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
		err = db.Debug().Model(&models.Business{}).Create(&businesses[i]).Error
		if err != nil {
			log.Fatalf("cannot seed business table %v\n", err)
		}
	}
}
