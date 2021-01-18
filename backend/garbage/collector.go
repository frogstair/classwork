package garbage

import (
	"classwork/backend/util"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var fileChannel chan []string

// Quit is the channel to quit the execution of the garbage collector
var Quit chan bool

func init() {
	fileChannel = make(chan []string)
}

// GetChannel gets the pointer to the file channel to notify the garbage collector
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
	for {
		select {
		case f := <-fileChannel:
			go func(files []string) {
				for i := 0; i < 10; i++ {
					time.Sleep(1 * time.Minute)
					for _, file := range files {
						if _, err := os.Stat(util.ToRelativeFPath(file)); os.IsNotExist(err) {
							return
						}
					}
				}
				for _, file := range files {
					file = util.ToRelativeFPath(file)
					if _, err := os.Stat(file); os.IsNotExist(err) {
						continue
					}
					err := os.Remove(file)
					if err != nil {
						panic(err)
					}
				}
				return
			}(f)
		case <-Quit:
			return
		}
	}
}
