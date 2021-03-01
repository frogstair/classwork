package middleware

import (
	"classwork/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" //Required by GORM
)

// Postgres attaches a database variable to a given context
func Postgres(c *gin.Context) {
	db := database.GetPostgres() // Get existing connection/establish new connection
	c.Set("db", db)              // Set context variable
	c.Next()
}
