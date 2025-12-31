package main

import (
	"fmt"
	"net/http"

	controllers "github.com/EsanSamuel/Reddit_Clone/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to reddit_clone api"})
	})
	r.POST("/register", controllers.CreateUser())
	r.PATCH("/verify-user", controllers.VerifyEmail())
	r.POST("/login", controllers.Login())

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Error starting server")
	}
}
