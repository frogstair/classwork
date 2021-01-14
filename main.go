package main

import (
	"log"
	"os"

	"classwork/api"
	"classwork/database"
	"classwork/middleware"
	"classwork/pages"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Could not find .env file!")
	}

	db := database.GetPostgres()
	defer db.Close()

	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(middleware.Postgres)
	g.NoRoute(pages.NotFound)

	g.GET("/:fname", pages.Serve)

	apiGroup := g.Group("/api")

	regGroup := apiGroup.Group("/register")
	regGroup.POST("/", api.Register)

	apiGroup.POST("/login", api.Login)

	g.Run(os.Getenv("ADDRESS") + ":" + os.Getenv("PORT"))
}
