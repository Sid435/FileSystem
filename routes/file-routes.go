package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sid/FileSystem/controllers"
)

var FileRoutes = func(router *gin.RouterGroup) {
	router.POST("/upload", controllers.UploadFile)
	router.GET("/get", controllers.GetSignedUrl)
	router.DELETE("/delete", controllers.DeleteFile)
}

var AuthRoutes = func(router *gin.RouterGroup) {
	router.POST("/signup", controllers.SignUp)
	router.POST("/login", controllers.LogIn)
}
