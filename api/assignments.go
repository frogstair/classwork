package api

import (
	m "classwork/models"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// NewAssignment creates a new assignment
func NewAssignment(c *gin.Context) {
	db, ok := c.Keys["db"].(*gorm.DB) // Get database from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	user, ok := c.Keys["usr"].(*m.User) // Get user from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no user variable in context")
	}

	if !user.Has(m.Teacher) { // Only teachers can create assignments
		c.JSON(403, gin.H{"error": "insufficient permissions"})
		return
	}

	newAssignment := new(m.NewAssignment) // Crearte a model
	err := json.NewDecoder(c.Request.Body).Decode(newAssignment)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := newAssignment.Create(db, user) // Call model function
	c.JSON(code, resp)
}

// CompleteAssignment completes an assignment
func CompleteAssignment(c *gin.Context) {
	db, ok := c.Keys["db"].(*gorm.DB) // Database variable from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	user, ok := c.Keys["usr"].(*m.User) // Get user from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no user variable in context")
	}

	if !user.Has(m.Student) { // Only students can complete their assignments
		c.JSON(403, gin.H{"error": "insufficient permissions"})
		return
	}

	newComplete := new(m.NewRequestComplete) // Create a model and run the method
	err := json.NewDecoder(c.Request.Body).Decode(newComplete)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := newComplete.Complete(db, user)
	c.JSON(code, resp)
}

// GetAssignment gets information about the assignment
func GetAssignment(c *gin.Context) {
	db, ok := c.Keys["db"].(*gorm.DB) // Get database from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	user, ok := c.Keys["usr"].(*m.User) // Get user from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no user variable in context")
	}

	gasn := new(m.GetAssignment) // Create models and run functions
	gasn.ID = c.Query("id")
	code, resp := gasn.Get(db, user)
	c.JSON(code, resp)
}
