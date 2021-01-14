package pages

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
)

// Serve serves a file from the web directory
func Serve(c *gin.Context) {
	name := c.Param("fname")
	reg := regexp.MustCompile("\\.{2,}")
	name = reg.ReplaceAllString(name, "")

	var filename string

	if filepath.Ext(name) != "" {
		filename = fmt.Sprintf("./web/%s", name)
	} else {
		filename = fmt.Sprintf("./web/%s.html", name)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		NotFound(c)
		return
	}

	c.File(filename)
}

// NotFound serves the not found response
func NotFound(c *gin.Context) {
	c.JSON(404, gin.H{"error": "not found"})
}
