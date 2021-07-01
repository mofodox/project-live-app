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
		Name:             "Eat My CB",
		UserID:           1,
		ShortDescription: "Homemade curry-stuffed buns",
		Description:      "",
		Address:          "175C Punggol Field",
		UnitNo:           "",
		Zipcode:          "823175",
		Lat:              0,
		Lng:              0,
		Website:          "https://eatmycb.sg/",
		Instagram:        "https://www.instagram.com/eatmycb/",
		Facebook:         "https://www.facebook.com/eatmycb/",
	},

	models.Business{
		Name:             "Makan Fix",
		UserID:           1,
		ShortDescription: "Makan Fix sells unique desserts like Baklava Rings and the Middle Eastern cheese-layered pastry, Kunafe, that you probably won’t be able to find in regular neighbourhood bakeries.",
		Description:      "Makan Fix sells unique desserts like Baklava Rings and the Middle Eastern cheese-layered pastry, Kunafe, that you probably won’t be able to find in regular neighbourhood bakeries.",
		Address:          "Gallop Kranji Farm Resort",
		UnitNo:           "",
		Zipcode:          "",
		Lat:              0,
		Lng:              0,
		Website:          "",
		Instagram:        "https://www.instagram.com/makanfix",
		Facebook:         "",
	},

	models.Business{
		Name:             "Le Vyr",
		UserID:           1,
		ShortDescription: "Locally fermented kombucha",
		Description:      "Brewed with an intense passion for our cuisine culture, our brews are inspired by and fermented with premium local ingredients.",
		Address:          "King George's Building",
		UnitNo:           "",
		Zipcode:          "208580",
		Lat:              0,
		Lng:              0,
		Website:          "https://www.levyrsg.com/",
		Instagram:        "https://www.instagram.com/levyr.sg/",
		Facebook:         "",
	},

	models.Business{
		Name:             "Hakka Chang",
		UserID:           1,
		ShortDescription: "Homemade Hakka Yong Tau Foo & Abacus Seed\nAvailable for Delivery/Self-Collection\nWhatsApp us to Pre-Order",
		Description:      "Freshly made every Sat (26 Jun)",
		Address:          "329 Sembawang Close",
		UnitNo:           "",
		Zipcode:          "750329",
		Lat:              0,
		Lng:              0,
		Website:          "",
		Instagram:        "https://www.instagram.com/hakkachang",
		Facebook:         "",
	},

	models.Business{
		Name:             "Upcakes",
		UserID:           1,
		ShortDescription: "Mao Shan Wang chocolate cakes",
		Description:      "Calling all durian lovers: Upcakes has one of the best durian chocolate cakes ever, IMO. Not only is the cake extremely moist on the inside, the flavours of the chocolate and durian are also perfectly balanced, as all things should be.",
		Address:          "Shelford Road",
		UnitNo:           "",
		Zipcode:          "288435",
		Lat:              0,
		Lng:              0,
		Website:          "https://www.upcakes.sg/",
		Instagram:        "https://www.instagram.com/upcakes.sg/",
		Facebook:         "",
	},

	models.Business{
		Name:             "Your Daily Batter",
		UserID:           1,
		ShortDescription: "Basque burnt cheesecakes",
		Description:      "DM to order!",
		Address:          "Jalan Kayu",
		UnitNo:           "",
		Zipcode:          "",
		Lat:              0,
		Lng:              0,
		Website:          "",
		Instagram:        "https://www.instagram.com/p/CJbKiqPHsnH/",
		Facebook:         "",
	},

	models.Business{
		Name:             "Okieco",
		UserID:           1,
		ShortDescription: "New York banana pudding by the pint",
		Description:      "Our Puddings are made from scratch, assembled by hand and contain only quality ingredients we approve of. Oh, and did we mention that they are made F R E S H?",
		Address:          "8 Somapah Road, SUTD, Building 2",
		UnitNo:           "#01-202A",
		Zipcode:          "487372",
		Lat:              0,
		Lng:              0,
		Website:          "https://www.okieco.sg/",
		Instagram:        "https://www.instagram.com/okieco.sg/",
		Facebook:         "",
	},

	models.Business{
		Name:             "The Munch Tray",
		UserID:           1,
		ShortDescription: "Freshly baked goodies to satisfy your sweet tooth cravings 🍪",
		Description:      "Two words: gooey cookies. If that’s not enough to pique your interest, you probably belong in the Crunchy Cookies camp. With The Munch Tray’s array of biscoff-filled and melty chocolate chip cookies to choose from, gooey cookie fans will be in for a treat.",
		Address:          "Ang Mo Kio, Singapore",
		UnitNo:           "",
		Zipcode:          "",
		Lat:              0,
		Lng:              0,
		Website:          "https://themunchtray.cococart.co",
		Instagram:        "https://www.instagram.com/themunchtray",
		Facebook:         "",
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
		Content:    "Best curry buns in Punggol area! 10/10",
	},
	{
		BusinessID: 1,
		UserID:     2,
		Content:    "I really like eating the CBs from here",
	},
	{
		BusinessID: 2,
		UserID:     1,
		Content:    "If you are a cheese lover, Kunafe is a must try",
	},
	{
		BusinessID: 3,
		UserID:     2,
		Content:    "Good brews.",
	},
	{
		BusinessID: 3,
		UserID:     2,
		Content:    "Must try for everyone.",
	},
	{
		BusinessID: 3,
		UserID:     2,
		Content:    "Good.",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.User{}, &models.Business{}, &models.Category{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v\n", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Business{}, &models.Category{}).Error
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
