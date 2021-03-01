package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	m "classwork/models"
)

// ValidateJWT validates the JWT token and places the user it belongs to in the context
func ValidateJWT(c *gin.Context) {

	db, ok := c.Keys["db"].(*gorm.DB) // Get database variable from context
	if !ok {                          // If not found
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	tokstr, err := c.Cookie("_tkn") // Get `_tkn` from cookie
	if err != nil {                 // If there is an error
		c.JSON(400, gin.H{"error": "no token specified"})
		c.Abort() // Doesnt run the next function in the chain
		return
	}

	code, data := m.ParseToken(tokstr, db) // Parse token
	if code != 200 {
		c.JSON(code, gin.H{"error": data.Error})
		c.Abort()
		return
	}

	c.Set("usr", data.Data) // Set user context variable
	c.Next()
}
