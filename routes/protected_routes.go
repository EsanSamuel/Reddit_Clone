package routes

import (
	"github.com/EsanSamuel/Reddit_Clone/middlewares"
	"github.com/gin-gonic/gin"
)

func ProtectedRoutes(r *gin.Engine) {
	r.Use(middlewares.AuthMiddleware())

}
