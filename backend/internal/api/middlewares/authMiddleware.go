package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required."})
			return
		}
		auth_header_slice := strings.Split(authHeader, " ")

		if len(auth_header_slice) > 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization Header."})
			return
		}

		token := auth_header_slice[1]
		parded_token, err := utils.ParseAndValidateToken(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Set("user_id", parded_token.UserId)
		c.Set("email_id", parded_token.Email)
		c.Next()
	}
}
