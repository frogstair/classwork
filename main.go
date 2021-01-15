package main

import (
	"log"
	"os"

	"classwork/api"
	"classwork/database"
	m "classwork/middleware"
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
	g.Use(m.Postgres)
	g.NoRoute(pages.NotFound)
	g.NoMethod(pages.NoMethod)

	apiGroup := g.Group("/api")

	logGroup := apiGroup.Group("/login")
	logGroup.POST("/", api.Login)
	logGroup.GET("/pass", api.GenerateOTC)

	regGroup := apiGroup.Group("/register")
	regGroup.POST("/", api.Register)
	regGroup.GET("/email", api.EmailValid)

	dashboardGroup := apiGroup.Group("/dashboard")
	dashboardGroup.GET("/", m.ValidateJWT, api.GetDashboard)

	schGroup := apiGroup.Group("/school")
	schGroup.POST("/", m.ValidateJWT, api.AddSchool)
	schGroup.DELETE("/", m.ValidateJWT, api.DeleteSchool)
	schGroup.POST("/teacher", m.ValidateJWT, api.AddTeacher)

	g.Run(os.Getenv("ADDRESS") + ":" + os.Getenv("PORT"))
}
