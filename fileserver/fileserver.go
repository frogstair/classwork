package fileserver

import (
	"classwork/fileserver/api"
	"classwork/fileserver/garbage"

	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

// Run runs the fileserver
func Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() { garbage.Quit <- true }()

	go garbage.Run()

	gin.SetMode(gin.ReleaseMode)
	g := gin.New()

	g.Use(gin.Recovery())
	g.Use(garbage.AddCollectorToContext)

	g.Static("/files", "./files")

	fsgroup := g.Group("/files")
	fsgroup.POST("/", api.CreateFile)

	address, port := os.Getenv("ADDRESS"), os.Getenv("PORT")
	log.Printf("Running fileserver on %s:%s\n", address, port)
	g.Run(address + ":" + port)
}
