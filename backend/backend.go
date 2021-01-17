package backend

import (
	"classwork/backend/api"
	"classwork/backend/database"
	"classwork/backend/garbage"
	m "classwork/backend/middleware"
	"classwork/backend/pages"

	"log"
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

	logoutGroup := apiGroup.Group("/logout")
	logoutGroup.POST("/", m.ValidateJWT, api.Logout)

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

	schGroup.GET("/info", m.ValidateJWT, api.GetSchool)

	schGroup.POST("/teacher", m.ValidateJWT, api.AddTeacher)
	schGroup.DELETE("/teacher", m.ValidateJWT, api.DeleteTeacher)

	schGroup.POST("/student", m.ValidateJWT, api.AddStudent)
	schGroup.DELETE("/student", m.ValidateJWT, api.DeleteStudent)

	subGroup := schGroup.Group("/subject")
	subGroup.POST("/", m.ValidateJWT, api.AddSubject)
	subGroup.DELETE("/", m.ValidateJWT, api.DeleteSubject)
	subGroup.POST("/students", m.ValidateJWT, api.AddStudentSubject)
	subGroup.DELETE("/students", m.ValidateJWT, nil)

	g.Use(garbage.AddCollectorToContext)
	go garbage.Run()

	g.Static("/files", "./files")

	fsgroup := g.Group("/files")
	fsgroup.POST("/", api.CreateFile)

	address, port := os.Getenv("ADDRESS"), os.Getenv("PORT")

	log.Printf("Running backend on %s:%s\n", address, port)

	g.Run(address + ":" + port)
}
