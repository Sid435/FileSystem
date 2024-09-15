package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/sid/FileSystem/pkg/controllers"
)

var FileRoutes = func(router *gin.RouterGroup) {
	router.POST("/upload", controllers.UploadFile)
	router.GET("/get", controllers.GetPreSignedURL) // Uncomment this to handle file retrieval by name
}
var OnboardingRoutes = func(router *mux.Router) {
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/signup", controllers.Signup).Methods("POST")
}
