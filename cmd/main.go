package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sid/FileSystem/pkg/config"
	"github.com/sid/FileSystem/pkg/controllers"
	"github.com/sid/FileSystem/pkg/routes"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
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

	log.Fatal(http.ListenAndServe("localhost:9010", r))
}
