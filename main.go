package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sid/FileSystem/config"
	"github.com/sid/FileSystem/controllers"
	"github.com/sid/FileSystem/routes"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	config.InitRedis()
	fileRouters := r.Group("/files")
	fileRouters.Use(controllers.JwtAuthMiddleware)
	{
		routes.FileRoutes(fileRouters)
	}

	fileRouters = r.Group("/auth")
	routes.AuthRoutes(fileRouters)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gin server failed")
	}
}
