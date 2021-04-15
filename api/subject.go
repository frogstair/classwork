package api

import (
	m "classwork/models"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AddSubject adds a new subject
func AddSubject(c *gin.Context) {
	db, ok := c.Keys["db"].(*gorm.DB) // Database variable from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	user, ok := c.Keys["usr"].(*m.User) // Get the user from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no user variable in context")
	}

	if !user.Has(m.Teacher) { // Only the teacher can add subjects
		c.JSON(403, gin.H{"error": "insufficient permissions"})
		return
	}

	newSubject := new(m.NewSubject) // Create a model
	err := json.NewDecoder(c.Request.Body).Decode(newSubject)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := newSubject.Add(db, user) // Run function
	c.JSON(code, resp)
}

// DeleteSubject deletes a subject
func DeleteSubject(c *gin.Context) { // Database variable from context
	db, ok := c.Keys["db"].(*gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	user, ok := c.Keys["usr"].(*m.User) // User variable from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no user variable in context")
	}

	if !user.Has(m.Teacher) { // Only the teacher can add subjects
		c.JSON(403, gin.H{"error": "insufficient permissions"})
		return
	}

	deleteSubject := new(m.DeleteSubject) // Create model
	deleteSubject.ID = c.Query("id")

	code, resp := deleteSubject.Delete(db, user) // Call function
	c.JSON(code, resp)
}

// GetSubject gets information about a subject
func GetSubject(c *gin.Context) {
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

	getSubject := new(m.GetSubjectInfo)
	getSubject.ID = c.Query("id") // Get the subject ID from query
	getSubject.SID = c.Query("sid") // Get the school ID from query
	code, resp := getSubject.Get(db, user)
	c.JSON(code, resp)
}

// AddStudentSubject adds a new student to a subject
func AddStudentSubject(c *gin.Context) {
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

	if !user.Has(m.Teacher) && !user.Has(m.Headmaster) { // Students cant add others to a subject
		c.JSON(403, gin.H{"error": "fobidden"})
		return
	}

	newStudentSubject := new(m.NewSubjectStudent) // Create model
	err := json.NewDecoder(c.Request.Body).Decode(newStudentSubject)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := newStudentSubject.Add(db, user) // Add the student
	c.JSON(code, resp)
}
