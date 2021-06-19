package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/mofodox/project-live-app/api/models"
)

var users = []models.User{
	models.User{
		Fullname: "John Doe",
		Email: "johndoe@gmail.com",
		Password: "password",
	},
	models.User{
		Fullname: "Martin Luther",
		Email: "luther@gmail.com",
		Password: "password",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v\n", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v\n", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table %v\n", err)
		}
	}
}
