package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	m "classwork/models"
)

// EmailValid checks if the supplied Email is valid
func EmailValid(c *gin.Context) {
	db, ok := c.Keys["db"].(*gorm.DB) // Get database variable from context

	if !ok { // If an error occured
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	email := new(m.Email)
	email.Email = c.Query("email") // Get email from query
	code, data := email.Valid(db)  // Check if its valid
	c.JSON(code, data)             // Respond
}
