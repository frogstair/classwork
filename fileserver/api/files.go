package api

import (
	"classwork/fileserver/util"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// CreateFile creates a new file for upload
func CreateFile(c *gin.Context) {

	gcChan, ok := c.Keys["collector"].(*chan []string)
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no channel variable in context")
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	files := form.File["files"]
	names := make([]string, len(files))

	for _, file := range files {
		if file.Size > 100000000 {
			c.JSON(400, gin.H{"error": fmt.Sprintf("File %s is too large, limit is 100MB", file.Filename)})
			return
		}
	}

	for i, file := range files {
		_, ext := util.SplitName(file.Filename)
		name := util.GenerateName()

		name = name + "_0" + ext

		name = util.ToRelativeFPath(name)

		err = c.SaveUploadedFile(file, name)
		if err != nil {
			c.JSON(500, gin.H{"error": "internal error"})
			log.Printf("File upload error %s\n", err.Error())
			return
		}

		names[i] = name
	}

	resp := struct {
		Files []string `json:"files"`
	}{names}
	*gcChan <- names

	c.JSON(200, resp)
}
