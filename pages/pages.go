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

// ServeLogin serves the login page
func ServeLogin(c *gin.Context) {
	data, _ := ioutil.ReadFile("./web/login.html")
	c.Data(200, "text/html; charset=utf-8", data)
}

// ServeLoginPassword serves the password page
func ServeLoginPassword(c *gin.Context) {
	data, _ := ioutil.ReadFile("./web/password.html")
	c.Data(200, "text/html; charset=utf-8", data)
}

// ServeDashboard serves the dashboard page
func ServeDashboard(c *gin.Context) {
	data, _ := ioutil.ReadFile("./web/dashboard.html")
	c.Data(200, "text/html; charset=utf-8", data)
}

// ServeSchool serves the school page
func ServeSchool(c *gin.Context) {
	data, _ := ioutil.ReadFile("./web/school.html")
	c.Data(200, "text/html; charset=utf-8", data)
}

// ServeSubject serves the school page
func ServeSubject(c *gin.Context) {
	data, _ := ioutil.ReadFile("./web/subject.html")
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
