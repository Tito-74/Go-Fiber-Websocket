package database

import (
	"log"

	"github.com/Tito-74/fiber-websocket/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)



var Db *gorm.DB
var err error

type DbInstance struct {
	Db *gorm.DB
	
}

var Database DbInstance


func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("error occure while initializing database connection", err)
	}
	db.AutoMigrate(&models.Message{},
	)

	Database = DbInstance{Db: db}

}


