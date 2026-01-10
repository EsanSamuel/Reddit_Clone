package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/EsanSamuel/Reddit_Clone/jobs/workers"
	"github.com/EsanSamuel/Reddit_Clone/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	go workers.EmailWorker()
	go workers.AISummaryWorker()
	go workers.AIEmbeddingWorker()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to reddit_clone api"})
	})
	routes.UnProtectedRoutes(r)
	//routes.ProtectedRoutes(r)

	go func() {
		if err := r.Run(":8080"); err != nil {
			fmt.Println("Error starting server")
		}
	}()

	// This is needed to gracefully shutdown the application with ctrl^c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	workers.StopEmailWorker()
	workers.StopAISummaryWorker()
	workers.StopAIEmbeddingWorker()
}
