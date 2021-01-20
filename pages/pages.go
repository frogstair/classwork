package pages

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// ServeRegister serves the register page
func ServeRegister(c *gin.Context) {
	data, _ := ioutil.ReadFile("./web/register.html")
	c.Data(200, "text/html; charset=utf-8", data)
}

// NotFound serves the not found response
func NotFound(c *gin.Context) {
	c.JSON(404, gin.H{"error": "not found"})
}

// NoMethod is thrown when no method is allowed
func NoMethod(c *gin.Context) {
	c.JSON(405, gin.H{"error": "method not allowed"})
}
