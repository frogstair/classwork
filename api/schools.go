package api

import (
	m "classwork/models"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AddSchool adds a school
func AddSchool(c *gin.Context) {
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

	newSchool := new(m.NewSchool)
	err := json.NewDecoder(c.Request.Body).Decode(newSchool)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := newSchool.Add(db, user)
	c.JSON(code, resp)
}

// DeleteSchool will delete a school from the database
func DeleteSchool(c *gin.Context) {
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

	deleteSchool := new(m.DeleteSchool)
	err := json.NewDecoder(c.Request.Body).Decode(deleteSchool)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := deleteSchool.Delete(db, user)
	c.JSON(code, resp)
}

// AddTeacher adds a teacher to the school
func AddTeacher(c *gin.Context) {
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

	newTeacher := new(m.NewTeacher)
	err := json.NewDecoder(c.Request.Body).Decode(newTeacher)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := newTeacher.Add(db)
	c.JSON(code, resp)
}
