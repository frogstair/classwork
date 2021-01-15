package middleware

import (
	m "classwork/models"
	"classwork/util"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
		return
	}

	tok := &m.Token{}
	token, err := jwt.ParseWithClaims(tokstr, tok,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

	if !token.Valid {
		if err != nil {
			c.JSON(401, gin.H{"error": "token is invalid"})
			return
		}
	}

	user := new(m.User)
	err = db.Where("id = ?", tok.ID).First(user).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			c.JSON(401, gin.H{"error": "token does not correspond to user"})
			return
		}
		log.Printf("Database error: %s\n", err.Error())
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.Set("usr", user)
	c.Next()
}
