package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sid/FileSystem/pkg/config"
	"github.com/sid/FileSystem/pkg/controllers"
	"github.com/sid/FileSystem/pkg/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	fmt.Println("Starting the servers...")
	router := gin.Default()
	config.InitRedis()

	fileRouter := router.Group("/files")
	fileRouter.Use(controllers.JWTAuthMiddleware)
	{
		routes.FileRoutes(fileRouter)
	}

	go func() {
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("Gin server failed: %v", err)
		}
	}()

	r := mux.NewRouter()
	onBoard := r.PathPrefix("/auth").Subrouter()
	routes.OnboardingRoutes(onBoard)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9010"
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
