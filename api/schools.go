package api

import (
	m "classwork/models"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AddSchool adds a school
func AddSchool(c *gin.Context) {
	db, ok := c.Keys["db"].(*gorm.DB) // Get database variable from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	user, ok := c.Keys["usr"].(*m.User) // Get the user to whom the school belongs
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no user variable in context")
	}

	if !user.Has(m.Headmaster) { // Teachers and students cant add schools
		c.JSON(403, gin.H{"error": "insufficient permissions"})
		return
	}

	newSchool := new(m.NewSchool) // Placeholder
	err := json.NewDecoder(c.Request.Body).Decode(newSchool)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := newSchool.Add(db, user) // Call function
	c.JSON(code, resp)
}

// GetSchool gets info about the school
func GetSchool(c *gin.Context) {
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

	schoolGetInfo := new(m.GetSchoolInfo)
	schoolGetInfo.ID = c.Query("id")

	code, resp := schoolGetInfo.GetInfo(db, user)
	c.JSON(code, resp)
}

// DeleteSchool will delete a school from the database
func DeleteSchool(c *gin.Context) {
	db, ok := c.Keys["db"].(*gorm.DB) // Database from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	user, ok := c.Keys["usr"].(*m.User) // User from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no user variable in context")
	}

	if !user.Has(m.Headmaster) { // Teachers and students cant delete schools
		c.JSON(403, gin.H{"error": "insufficient permissions"})
		return
	}

	deleteSchool := new(m.DeleteSchool) // Placeholder
	deleteSchool.ID = c.Query("id")     // Get schoold ID from query

	code, resp := deleteSchool.Delete(db, user) // Call function
	c.JSON(code, resp)
}

// GetStudents will get the students from a school
func GetStudents(c *gin.Context) {
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

	if user.Has(m.Student) {
		c.JSON(403, gin.H{"error": "insufficient permissions"})
		return
	}

	getStudents := new(m.GetStudents)
	getStudents.ID = c.Query("id")

	code, res := getStudents.Get(db)
	c.JSON(code, res)
}
