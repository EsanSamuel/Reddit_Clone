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
	r.PATCH("/avatar/:userId", controllers.UploadAvatar())
	r.POST("/subreddit", controllers.CreateSubreddit())
	r.POST("/subreddit/member", controllers.JoinSubreddit())
	r.POST("subreddit/moderator", controllers.AddModerators())
	//r.POST("/upload", controllers.UploadFiles())
}
