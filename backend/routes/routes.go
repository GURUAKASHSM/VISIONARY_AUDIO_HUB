package routes

import (
	"audiohub/controller"

	"github.com/gin-gonic/gin"
)

// Router creates and configures the Gin router.
func Router() *gin.Engine {
	router := gin.Default()
	router.Use(controller.CorsMiddleware())
	router.POST("/login", controller.Login)
	router.POST("/register", controller.CreateProfile)
	router.POST("/verify", controller.VerifyEmail)
	router.POST("/forgetpassword", controller.ForgetPassword)
	router.POST("/passwordchange", controller.PasswordChange)
	router.POST("/getuserbyid", controller.LisUserDetails)

	return router

}
