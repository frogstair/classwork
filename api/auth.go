package api

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	m "classwork/models"
)

// Register is used to register a new user
func Register(c *gin.Context) {
	db, ok := c.Keys["db"].(*gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	reguser := new(m.RegisterUser)
	err := json.NewDecoder(c.Request.Body).Decode(reguser)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, data := reguser.Register(db)
	c.JSON(code, data)
}

// Login creates a token for the user to use for future authentication
func Login(c *gin.Context) {
	db, ok := c.Keys["db"].(*gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	loginuser := new(m.LoginUser)
	err := json.NewDecoder(c.Request.Body).Decode(loginuser)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, data, tok := loginuser.Login(db)
	c.SetCookie("_tkn", tok, int(m.TokenValidity), "/api/", "", false, true)

	c.JSON(code, data)
}

// GenerateOTC creates an OTC for a user if their password is not set
func GenerateOTC(c *gin.Context) {
	db, ok := c.Keys["db"].(*gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	hasPass := new(m.OTCCreate)
	err := json.NewDecoder(c.Request.Body).Decode(hasPass)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	code, resp := hasPass.Create(db)
	c.JSON(code, resp)
}
