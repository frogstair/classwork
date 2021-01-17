package api

import (
	m "classwork/backend/models"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// NewAssignment creates a new assignment
func NewAssignment(c *gin.Context) {
	db, ok := c.Keys["db"].(*gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	user, ok := c.Keys["usr"].(*m.User)
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no user variable in context")
	}

	if !user.Has(m.Teacher) {
		c.JSON(403, gin.H{"error": "insufficient permissions"})
		return
	}

	newAssignment := new(m.NewAssignment)
	err := json.NewDecoder(c.Request.Body).Decode(newAssignment)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := newAssignment.Create(db, user)
	c.JSON(code, resp)
}