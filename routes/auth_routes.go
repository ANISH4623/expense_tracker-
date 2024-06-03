package routes

import (
	"awesomeProject1/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/signup", controllers.Signup())
	router.POST("/login", controllers.Login())
	router.POST("/logout", controllers.Logout())

}
