package api

import (
	m "classwork/backend/models"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AddSubject adds a new subject
func AddSubject(c *gin.Context) {
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

	newSubject := new(m.NewSubject)
	err := json.NewDecoder(c.Request.Body).Decode(newSubject)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := newSubject.Add(db, user)
	c.JSON(code, resp)
}
