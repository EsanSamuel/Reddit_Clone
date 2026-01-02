package middlewares

import (
	"fmt"
	"net/http"

	"github.com/EsanSamuel/Reddit_Clone/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetAuthToken(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting auth token", "details": err.Error()})
			c.Abort()
			return
		}

		if token == "" {
			fmt.Println("Token is empty")
			c.Abort()
			return
		}

		claims, err := utils.VerifyAuthToken(token)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error verifying auth token", "details": err.Error()})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)

	}
}
