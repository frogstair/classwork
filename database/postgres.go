package database

import (
	"fmt"
	"log"
	"os"

	m "classwork/models"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB       // Global variable that keeps the database connection
var connected = false // Flag if the connection is established

// GetPostgres initializes the connection to a postgres database or returns an existing connection
func GetPostgres() *gorm.DB {

	if connected { // If the connection is already established then return the connection
		return db
	}

	host := os.Getenv("DB_ADDR") // Get variables from the environment
	port := os.Getenv("DB_PORT")
	role := os.Getenv("DB_ROLE")
	name := os.Getenv("DB_NAME")
	var err error // Variable to contain the error

	cstring := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable port=%s",
		host, role, name, port) // Connection string, formed using a format string

	db, err = gorm.Open("postgres", cstring) // Open a "postgres" database connection
	// If an error occurred then panic and exit the program
	if err != nil {
		panic(fmt.Sprintf("\n===========\ncannot establish database connection: \n%s\n===========", err))
	}
	// Set the flag
	connected = true
	// Remove unnecessary database logs
	db.LogMode(false)

	// Print a success log
	log.Println("Connected to database")

	// Synchronize all the tables in the code with the database
	db.AutoMigrate(&m.User{}, &m.School{}, &m.Subject{}, &m.Assignment{}, &m.Request{}, &m.AssignmentFile{}, &m.RequestUpload{})
	log.Println("Migrated tables")

	// Return the newly established connection
	return db
}

// Disconnect closes the database connection
func Disconnect() {
	// Close the connection to the database
	db.Close()
	// Set the flag to false
	connected = false
}
