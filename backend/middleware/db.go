package middleware

import (
	"classwork/backend/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" //Required by GORM
)

// Postgres attaches a database variable to a given context
func Postgres(c *gin.Context) {
	db := database.GetPostgres()
	c.Set("db", db)
	c.Next()
}
