package api

import (
	m "classwork/models"
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

// DeleteSubject deletes a subject
func DeleteSubject(c *gin.Context) {
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

	deleteSubject := new(m.DeleteSubject)
	err := json.NewDecoder(c.Request.Body).Decode(deleteSubject)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := deleteSubject.Delete(db, user)
	c.JSON(code, resp)
}

// AddStudentSubject adds a new student to a subject
func AddStudentSubject(c *gin.Context) {
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

	if !user.Has(m.Teacher) && !user.Has(m.Headmaster) {
		c.JSON(403, gin.H{"error": "fobidden"})
		return
	}

	newStudentSubject := new(m.NewSubjectStudent)
	err := json.NewDecoder(c.Request.Body).Decode(newStudentSubject)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := newStudentSubject.Add(db, user)
	c.JSON(code, resp)
}