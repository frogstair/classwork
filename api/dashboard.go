package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	m "classwork/models"
)

// GetDashboard gets the dashboard for a user
func GetDashboard(c *gin.Context) {
	db, ok := c.Keys["db"].(*gorm.DB) // Get database from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	user, ok := c.Keys["usr"].(*m.User) // Get the user from the context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no user variable in context")
	}

	code, resp := user.GetDashboard(db) // Get the dashboard for the user
	c.JSON(code, resp)
}
