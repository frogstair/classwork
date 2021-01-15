package api

import (
	"encoding/json"

	m "classwork/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AddStudent adds a new student
func AddStudent(c *gin.Context) {
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

	if !user.Has(m.Headmaster) {
		c.JSON(403, gin.H{"error": "insufficient permissions"})
		return
	}

	newStudent := new(m.NewStudent)
	err := json.NewDecoder(c.Request.Body).Decode(newStudent)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := newStudent.Add(db)
	c.JSON(code, resp)
}
