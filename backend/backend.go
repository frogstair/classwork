package backend

import (
	"classwork/backend/api"
	"classwork/backend/database"
	m "classwork/backend/middleware"
	"classwork/backend/pages"
	"sync"

	"os"

	"github.com/gin-gonic/gin"
)

// Run runs the backend server
func Run(wg *sync.WaitGroup) {
	db := database.GetPostgres()
	defer db.Close()
	defer wg.Done()

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
	schGroup.DELETE("/teacher", m.ValidateJWT, api.DeleteTeacher)
	schGroup.POST("/student", m.ValidateJWT, api.AddStudent)

	g.Run(os.Getenv("ADDRESS") + ":" + os.Getenv("PORT"))
}
