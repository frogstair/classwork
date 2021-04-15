package api

import (
	"classwork/util"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// CreateFile creates a new file for upload
func CreateFile(c *gin.Context) {

	gcChan, ok := c.Keys["collector"].(*chan []string) // Get the garbage file collector from context
	if !ok {
		c.JSON(500, gin.H{"error": "internal error"})
		panic("no channel variable in context")
	}

	// The user sends the file as a multipart form with the files field
	form, err := c.MultipartForm() // Get the data that user sends
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid data"})
		return
	}

	// Get all files from the form
	files := form.File["files"]
	names := make([]string, len(files))

	// If no files were sent, then dont do anthing
	if len(files) == 0 {
		return
	}

	// For each file check the size, if its more than 100mb then exit
	for _, file := range files {
		if file.Size > 100000000 {
			c.JSON(400, gin.H{"error": fmt.Sprintf("File %s is too large, limit is 100MB", file.Filename)})
			return
		}
	}

	// For each file
	for i, file := range files {
		// Remove extension
		_, ext := util.SplitName(file.Filename)
		// Generate a random name for a file
		name := util.GenerateName()

		// Add _0 to the end and the extension, to mark that the file isnt
		// verified
		name = name + "_0" + ext
		fname := name

		// Get the global path of the file
		name = util.ToGlobalPath(name)

		// Save the uploaded file into the disk
		err = c.SaveUploadedFile(file, name)
		if err != nil {
			c.JSON(500, gin.H{"error": "internal error"})
			log.Printf("File upload error %s\n", err.Error())
			return
		}

		names[i] = fname
	}
	// Send the files to the garbage collector
	*gcChan <- names

	// respond to the user with the uploaded files and ways to access them
	resp := struct {
		Files []string `json:"files"`
	}{names}

	c.JSON(200, resp)
}
