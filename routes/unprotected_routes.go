package routes

import (
	"github.com/EsanSamuel/Reddit_Clone/controllers"
	"github.com/gin-gonic/gin"
)

func UnProtectedRoutes(r *gin.Engine) {
	r.POST("/register", controllers.CreateUser())
	r.PATCH("/verify-user", controllers.VerifyEmail())
	r.POST("/login", controllers.Login())
	r.PATCH("/reset-password", controllers.ResetPassword())
	r.PATCH("/reset-password-request", controllers.ResetPasswordRequest())
	r.GET("/users", controllers.GetAllUsers())
	r.GET("/users/:userId", controllers.GetUser())
	r.POST("/upload", controllers.UploadFiles())
}
