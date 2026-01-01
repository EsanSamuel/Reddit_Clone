package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	controllers "github.com/EsanSamuel/Reddit_Clone/controllers"
	"github.com/EsanSamuel/Reddit_Clone/jobs/workers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	go workers.EmailWorker()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to reddit_clone api"})
	})
	r.POST("/register", controllers.CreateUser())
	r.PATCH("/verify-user", controllers.VerifyEmail())
	r.POST("/login", controllers.Login())
	r.PATCH("/reset-password-request", controllers.ResetPasswordRequest())
	r.PATCH("/reset-password", controllers.ResetPassword())

	// This is needed to gracefully shutdown the application with ctrl^c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	workers.StopEmailWorker()

	go func() {
		if err := r.Run(":8080"); err != nil {
			fmt.Println("Error starting server")
		}
	}()
}
