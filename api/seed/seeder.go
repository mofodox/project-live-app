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

// https://thesmartlocal.com/read/home-based-businesses-singapore/
var businesses = []models.Business{
	models.Business{
		Name:        "Eat My CB",
		UserID:      1,
		Description: "Homemade curry-stuffed buns",
		Address:     "175C Punggol Field",
		UnitNo:      "",
		Zipcode:     "823175",
		Lat:         0,
		Lng:         0,
		Website:     "https://eatmycb.sg/",
		Instagram:   "https://www.instagram.com/eatmycb/",
		Facebook:    "https://www.facebook.com/eatmycb/",
	},

	models.Business{
		Name:        "Bekel Mama",
		UserID:      1,
		Description: "Homemade Malay food. Quality Malay food delights. $9 island wide delivery for lunch or dinner. Click on this link to ORDER: <a href=\"wa.me/6598266648\">https://wa.me/6598266648</a>",
		Address:     "",
		UnitNo:      "",
		Zipcode:     "",
		Lat:         0,
		Lng:         0,
		Website:     "",
		Instagram:   "https://www.instagram.com/bekalmama.sg/",
		Facebook:    "",
	},

	models.Business{
		Name:        "Le Vyr - Locally fermented kombucha",
		UserID:      1,
		Description: "Brewed with an intense passion for our cuisine culture, our brews are inspired by and fermented with premium local ingredients.",
		Address:     "King George's Building",
		UnitNo:      "",
		Zipcode:     "208580",
		Lat:         0,
		Lng:         0,
		Website:     "https://www.levyrsg.com/",
		Instagram:   "https://www.instagram.com/levyr.sg/",
		Facebook:    "",
	},

	models.Business{
		Name:        "Ms Chili SG – freshly made Indonesian keropok and chilli",
		UserID:      1,
		Description: "Bring Indonesian hometown flavours to you ❤️<br/>🌶️ Freshly made to order<br/>🌶️ Pre-order every Sunday<br/>🌶️ No pork / No lard<br/>Order here 👇🏻<br/><a href=\"https://mschili.cococart.co/\">https://mschili.cococart.co/</a>",
		Address:     "",
		UnitNo:      "",
		Zipcode:     "",
		Lat:         0,
		Lng:         0,
		Website:     "https://mschili.cococart.co/",
		Instagram:   "https://www.instagram.com/ms.chili_sg/",
		Facebook:    "",
	},

	models.Business{
		Name:        "Upcakes – Mao Shan Wang chocolate cakes",
		UserID:      1,
		Description: "Calling all durian lovers: Upcakes has one of the best durian chocolate cakes ever, IMO. Not only is the cake extremely moist on the inside, the flavours of the chocolate and durian are also perfectly balanced, as all things should be.",
		Address:     "Shelford Road",
		UnitNo:      "",
		Zipcode:     "288435",
		Lat:         0,
		Lng:         0,
		Website:     "https://www.upcakes.sg/",
		Instagram:   "https://www.instagram.com/upcakes.sg/",
		Facebook:    "",
	},

	models.Business{
		Name:        "Your Daily Batter - basque burnt cheesecakes",
		UserID:      1,
		Description: "DM to order!",
		Address:     "Jalan Kayu",
		UnitNo:      "",
		Zipcode:     "",
		Lat:         0,
		Lng:         0,
		Website:     "",
		Instagram:   "https://www.instagram.com/p/CJbKiqPHsnH/",
		Facebook:    "",
	},

	models.Business{
		Name:        "Okieco - New York banana pudding by the pint",
		UserID:      1,
		Description: "Our Puddings are made from scratch, assembled by hand and contain only quality ingredients we approve of. Oh, and did we mention that they are made F R E S H?",
		Address:     "8 Somapah Road, SUTD, Building 2",
		UnitNo:      "#01-202A",
		Zipcode:     "487372",
		Lat:         0,
		Lng:         0,
		Website:     "https://www.okieco.sg/",
		Instagram:   "https://www.instagram.com/okieco.sg/",
		Facebook:    "",
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

var comments = []models.Comment{
	{
		BusinessID: 1,
		UserID:     1,
		Content:    "Test comment 1",
	},
	{
		BusinessID: 1,
		UserID:     2,
		Content:    "Test comment 2",
	},
	{
		BusinessID: 2,
		UserID:     1,
		Content:    "Test comment 3",
	},
	{
		BusinessID: 3,
		UserID:     2,
		Content:    "Test comment 4",
	},
	{
		BusinessID: 3,
		UserID:     2,
		Content:    "Test comment 5",
	},
	{
		BusinessID: 3,
		UserID:     2,
		Content:    "Test comment 6",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.User{}, &models.Business{}, &models.File{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v\n", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Business{}, &models.File{}).Error
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

	err = db.Debug().DropTableIfExists(&models.Comment{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v\n", err)
	}
	err = db.Debug().AutoMigrate(&models.Comment{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v\n", err)
	}
	for i := range comments {
		err = db.Debug().Model(&models.Comment{}).Create(&comments[i]).Error
		if err != nil {
			log.Fatalf("cannot seed coments table %v\n", err)
		}
	}
}
