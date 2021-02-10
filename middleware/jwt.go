package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	m "classwork/models"
)

// ValidateJWT validates the JWT token and places the user it belongs to in the contetx
func ValidateJWT(c *gin.Context) {

	db, ok := c.Keys["db"].(*gorm.DB)
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no database variable in context")
	}

	tokstr, err := c.Cookie("_tkn")
	if err != nil {
		c.JSON(400, gin.H{"error": "no token specified"})
		c.Abort()
		return
	}

	code, data := m.ParseToken(tokstr, db)
	if code != 200 {
		c.JSON(code, gin.H{"error": data.Error})
		c.Abort()
		return
	}

	c.Set("usr", data.Data)
	c.Next()
}
