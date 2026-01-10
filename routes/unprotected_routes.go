package routes

import (
	"github.com/EsanSamuel/Reddit_Clone/controllers"
	"github.com/gin-gonic/gin"
)

func UnProtectedRoutes(r *gin.Engine) {
	r.POST("/register", controllers.CreateUser())
	r.PATCH("/verify-user", controllers.VerifyEmail())
	r.POST("/login", controllers.Login())
	r.POST("/logout", controllers.LogoutHandler())
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
	r.POST("/posts", controllers.CreatePost())
	r.GET("/posts", controllers.GetPosts())
	r.GET("/posts/subreddit/:subreddit_id", controllers.GetSubRedditPosts())
	r.GET("/tags/posts", controllers.GetTagPosts())
	r.GET("/posts/:id", controllers.GetPostById())
	r.POST("/comments", controllers.CreateComment())
	r.GET("/comments/post/:post_id", controllers.GetPostComments())
	r.GET("/comments/parent/:parent_id)", controllers.GetParentComments())
	r.GET("/comments/:id", controllers.GetCommentById())
	r.POST("/post/upvote", controllers.UpVotePost())
	r.POST("/post/downvote", controllers.DownVotePost())
	r.GET("/summary/:post_id", controllers.ThreadsSummary())
	r.POST("/rag/:postId", controllers.SeachPostDetailsWithAI())
	//r.POST("/upload", controllers.UploadFiles())
}
