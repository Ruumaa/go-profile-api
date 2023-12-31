package main

import (
	"go-web/controllers"
	"go-web/database"
	"go-web/helpers"
	"go-web/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	helpers.LoadEnvVariables()
	database.ConnectDB()
}

func main() {
	router := gin.Default()

	router.GET("/user/login/", controllers.GetAllUsers)
	router.POST("/user/register/", controllers.Register)
	router.POST("/user/login/", controllers.Login)
	router.PUT("/user/:id/", controllers.EditUser)
	router.DELETE("/user/:id/", controllers.DeleteUser)

	router.GET("/photos/", controllers.GetAllPhotos)
	router.POST("/photos/", controllers.CreatePhoto)
	router.PUT("/:photoId/", middleware.RequireAuth, controllers.EditPhoto)
	router.DELETE("/:photoId/", middleware.RequireAuth, controllers.DeletePhoto)

	router.GET("/validate/", middleware.RequireAuth, controllers.Validate)
	router.GET("/user/logout/", middleware.RequireAuth, controllers.Logout)

	router.Run()
}
