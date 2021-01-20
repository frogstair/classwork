package database

import (
	"fmt"
	"log"
	"os"

	m "classwork/models"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var connected = false

// GetPostgres initializes the connection to a postgres database or returns an existing connection
func GetPostgres() *gorm.DB {

	if connected {
		return db
	}

	host := os.Getenv("DB_ADDR")
	port := os.Getenv("DB_PORT")
	role := os.Getenv("DB_ROLE")
	name := os.Getenv("DB_NAME")
	var err error

	cstring := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable port=%s",
		host, role, name, port)

	db, err = gorm.Open("postgres", cstring)
	if err != nil {
		panic(fmt.Sprintf("===========\ncannot establish database connection: \n%s\n===========", err))
	}
	connected = true
	//db.LogMode(false)

	log.Println("Connected to database")

	db.AutoMigrate(&m.User{}, &m.School{}, &m.Subject{}, &m.Assignment{}, &m.Request{}, &m.AssignmentFile{}, &m.RequestUpload{})
	log.Println("Migrated tables")

	return db
}
