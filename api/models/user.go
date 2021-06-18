package models

import (
	"html"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email"`
	Fullname string `json:"full_name"`
	Password string `json:"password"`
}

/**
	To generate hashed password
	* @param password gets the password from the user
*/
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

/**
	* Before saving, we get user password as string
*/
func (user *User) BeforeSave() error {
	hashedPassword, err := Hash(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return nil
}

/**
	* To verify the password with hashed password is equal
	* @param hashedPassword get hashed password
	* @param password get string user password
*/
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}


func (user *User) Prepare() {
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.Fullname = html.EscapeString(strings.TrimSpace(user.Fullname))
}

/**
	* Create user and save it in the database
	* @param db pointer the DB
*/
func (user *User) SaveUser(db *gorm.DB) (*User, error) {
	err := db.Debug().Create(*user).Error
	if err != nil {
		return &User{}, err
	}

	return user, nil
}

/**
	* Find all users
	* @param db pointer the DB
*/
func (user *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	
	if err := db.Debug().Model(&User{}).Find(&users).Error; err != nil {
		return &[]User{}, err
	}

	return &users, nil
}

/**
	* Find a user by ID
	* @param db pointer the DB
	* @param uid userId
*/
func (user *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	err := db.Debug().Model(&User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		return &User{}, err
	}

	return user, nil
}

/**
	* Update a user by ID
	* @param db pointer the DB
	* @param uid userId
*/
func (user *User) UpdateUserByID(db *gorm.DB, uid uint32) (*User, error) {
	err := user.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{} {
			"password": user.Password,
			"email": user.Email,
			"fullname": user.Fullname,
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}

	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		return &User{}, err
	}

	return user, nil
}

/**
	* Delete a user by ID
	* @param db pointer the DB
	* @param uid userId
*/
func (user *User) DeleteUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
