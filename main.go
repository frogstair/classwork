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
	err := godotenv.Load() // Load all values from the .env file
	if err != nil {        // If an error occurred then exit
		log.Fatalln("Could not find .env file!")
	}

	rand.Seed(time.Now().UnixNano()) // Seed random number generator

	wg := sync.WaitGroup{} // Create a waitgroup to run everything asynchronously
	wg.Add(1)              // Only one function will be running asynchronously

	go run(&wg) // Run the function

	wg.Wait() // Wait for the program to exit
}

func run(wg *sync.WaitGroup) {
	db := database.GetPostgres()            // Get a database connection
	defer db.Close()                        // Close the database connection when the function returns
	defer wg.Done()                         // Mark the function as done when the function returns
	defer func() { garbage.Quit <- true }() // Quit the GC once the program exits

	gin.SetMode(gin.ReleaseMode) // Set gin's mode to release to remove any logs
	g := gin.New()               // Create a new router
	g.Use(gin.Recovery())        // Use the recovery middleware to recover from functions that may have crashed
	g.Use(m.Postgres)            // Use the postgres middleware to inject the database connection into every function
	g.NoRoute(pages.NotFound)
	g.NoMethod(pages.NoMethod)

	// Create routes to serve all the html pages
	// Could have made a smarter system than that
	// but it created router conflict
	g.GET("/", pages.ServeLogin)
	g.GET("/login", pages.ServeLogin)
	g.GET("/register", pages.ServeRegister)
	g.GET("/login/pass", pages.ServeLoginPassword)
	g.GET("/dashboard", pages.ServeDashboard)
	g.GET("/school", pages.ServeSchool)
	g.GET("/subject", pages.ServeSubject)
	g.GET("/assignment", pages.ServeAssignment)
	g.Static("/static", "./web/static")

	// All the API routes and their handlers
	apiGroup := g.Group("/api")

	logoutGroup := apiGroup.Group("/logout")
	logoutGroup.POST("/", m.ValidateJWT, api.Logout)

	logGroup := apiGroup.Group("/login")
	logGroup.POST("/", api.Login)
	logGroup.POST("/pass", api.GenerateOTC)

	regGroup := apiGroup.Group("/register")
	regGroup.POST("/", api.Register)
	regGroup.GET("/email", api.EmailValid)

	dbdGroup := apiGroup.Group("/dashboard")
	dbdGroup.GET("/", m.ValidateJWT, api.GetDashboard)

	schGroup := apiGroup.Group("/school")
	schGroup.POST("/", m.ValidateJWT, api.AddSchool)
	schGroup.DELETE("/", m.ValidateJWT, api.DeleteSchool)
	schGroup.GET("/", m.ValidateJWT, api.GetSchool)
	schGroup.GET("/student", m.ValidateJWT, api.GetStudents)

	schGroup.POST("/teacher", m.ValidateJWT, api.AddTeacher)
	schGroup.DELETE("/teacher", m.ValidateJWT, api.DeleteTeacher)

	schGroup.POST("/student", m.ValidateJWT, api.AddStudent)
	schGroup.DELETE("/student", m.ValidateJWT, api.DeleteStudent)

	subGroup := schGroup.Group("/subject")
	subGroup.POST("/", m.ValidateJWT, api.AddSubject)
	subGroup.DELETE("/", m.ValidateJWT, api.DeleteSubject)
	subGroup.GET("/", m.ValidateJWT, api.GetSubject)
	subGroup.POST("/students", m.ValidateJWT, api.AddStudentSubject)

	assgnGroup := subGroup.Group("/assignment")
	assgnGroup.POST("/", m.ValidateJWT, api.NewAssignment)
	assgnGroup.GET("/", m.ValidateJWT, api.GetAssignment)
	assgnGroup.POST("/complete", m.ValidateJWT, api.CompleteAssignment)

	// Garbage collection
	g.Use(garbage.AddCollectorToContext)
	go garbage.Run()

	g.Static("/files", "./files")
	fsgroup := g.Group("/files")
	fsgroup.POST("/", api.CreateFile)

	// Get the address and port from the environment
	address, port := os.Getenv("ADDRESS"), os.Getenv("PORT")

	// Print a log
	log.Printf("Running on %s:%s\n", address, port)

	// Run the server
	g.Run(address + ":" + port)
}
