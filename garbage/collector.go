package garbage

import (
	"classwork/util"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// The channel which is used to transfer files to the garbage collector
var fileChannel chan []string

// Quit is the channel to quit the execution of the garbage collector
var Quit chan bool

// Init is run before the main() function
func init() {
	fileChannel = make(chan []string)
}

// GetChannel gets the pointer to the garbage collector file channel
func GetChannel() *chan []string {
	return &fileChannel
}

// AddCollectorToContext adds the garbage collector channel to the context
func AddCollectorToContext(c *gin.Context) {
	c.Set("collector", GetChannel())
	c.Next()
}

// Run runs the garbage collector
func Run() {
	log.Println("Started garbage collector") // Write out a log
	for { // Check which channel has incoming information
		select {
		case f := <-fileChannel: // If there is an incoming file
			go func(files []string) { // Run a function on a new goroutine
				defer log.Printf("Cleaned %v", files) // When it returns, write out the log with all the cleaned files
				for i := 0; i < 10; i++ { // Check each minute if the files are still unverified
					time.Sleep(1 * time.Minute)
					for _, file := range files { // For each file check if it exists
						if _, err := os.Stat(util.ToGlobalPath(file)); os.IsNotExist(err) {
							return // If a file doesnt exist that means it was verified and the goroutine can be stopped
						}
					}
				}
				// After ten minutes delete all the files
				for _, file := range files {
					// Check if the file exists
					file = util.ToGlobalPath(file)
					if _, err := os.Stat(file); os.IsNotExist(err) {
						continue
					}
					// Remove the file
					err := os.Remove(file)
					if err != nil {
						panic(err) // If an unexpected error happened then panic
					}
				}
				return
			}(f)
		// When a quit signal is received then stop the garbage collector
		case <-Quit:
			log.Printf("Stopping garbage collection")
			return
		}
	}
}
