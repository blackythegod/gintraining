package database

import (
	"gintraining/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Database struct {
	*gorm.DB
}

func (db *Database) CheckCreds(user *models.User) bool {
	db.Find(&user, "email = ? AND password = ?", user.Email, user.Password)
	if user.Username == "" {
		return false
	} else {
		return true
	}
}
func (db *Database) CheckCredsForExisting(user *models.User) bool {
	db.Find(&user, "username = ?", user.Username)
	return user.Username != ""
}

func InitDB() *Database {
	db, err := gorm.Open(postgres.Open("postgresql://postgres:postgrespw@db:5432/gindb"), &gorm.Config{})
	if err != nil {
		log.Fatalf("server couldn't be started: %s", err)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("couldn't migrate tables: %s", err)
	}
	return &Database{db}
}
