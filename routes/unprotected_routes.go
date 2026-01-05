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
	r.GET("/subreddits", controllers.GetSubReddit())
	r.GET("/subreddits/user/:user_id", controllers.GetSubRedditUserJoined())
	r.GET("/subreddits/:id", controllers.GetSubRedditById())
	//r.POST("/upload", controllers.UploadFiles())
}
