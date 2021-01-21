package main

import (
	"classwork/api"
	"classwork/database"
	"classwork/garbage"
	m "classwork/middleware"
	"classwork/pages"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Could not find .env file!")
	}

	rand.Seed(time.Now().UnixNano())

	wg := sync.WaitGroup{}
	wg.Add(2)

	go run(&wg)

	wg.Wait()
}

func run(wg *sync.WaitGroup) {
	db := database.GetPostgres()
	defer db.Close()
	defer wg.Done()
	defer func() { garbage.Quit <- true }()

	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(m.Postgres)
	g.NoRoute(pages.NotFound)
	g.NoMethod(pages.NoMethod)

	g.GET("/register", pages.ServeRegister)
	g.GET("/login", pages.ServeLogin)
	g.GET("/login/pass", pages.ServeLoginPassword)
	g.GET("/dashboard", pages.ServeDashboard)
	g.Static("/static", "./web/static")

	apiGroup := g.Group("/api")

	logoutGroup := apiGroup.Group("/logout")
	logoutGroup.POST("/", m.ValidateJWT, api.Logout)

	logGroup := apiGroup.Group("/login")
	logGroup.POST("/", api.Login)
	logGroup.POST("/pass", api.GenerateOTC)

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

	assgnGroup := subGroup.Group("/assignment")
	assgnGroup.POST("/", m.ValidateJWT, api.NewAssignment)
	assgnGroup.POST("/complete", m.ValidateJWT, api.CompleteAssignment)

	g.Use(garbage.AddCollectorToContext)
	go garbage.Run()

	g.Static("/files", "./files")

	fsgroup := g.Group("/files")
	fsgroup.POST("/", api.CreateFile)

	address, port := os.Getenv("ADDRESS"), os.Getenv("PORT")

	log.Printf("Running on %s:%s\n", address, port)

	g.Run(address + ":" + port)
}
