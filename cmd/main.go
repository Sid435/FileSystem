package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/sid/FileSystem/pkg/config"
	"github.com/sid/FileSystem/pkg/controllers"
	"github.com/sid/FileSystem/pkg/routes"
)

func main() {
	fmt.Println("this is me")
	r := mux.NewRouter() // getting the route
	router := gin.Default()
	config.InitRedis()

	// Apply JWTAuthMiddleware to all /files routes
	fileRouter := router.Group("/files")
	fileRouter.Use(controllers.JWTAuthMiddleware)
	{
		routes.FileRoutes(fileRouter) // Assuming this is where file routes are defined
	}
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
	onBoard := r.PathPrefix("/auth").Subrouter()
	routes.OnboardingRoutes(onBoard)

	log.Fatal(http.ListenAndServe("localhost:9010", r))
}
