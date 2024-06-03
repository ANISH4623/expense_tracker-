package middleware

import (
	"awesomeProject1/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var TokenController *helpers.JWTToken

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
			c.Abort()
			return
		}

		tokenSplit := strings.Split(token, " ")

		if len(tokenSplit) != 2 || strings.ToLower(tokenSplit[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token, expects bearer token"})
			c.Abort()
			return
		}

		userId, err := TokenController.VerifyToken(tokenSplit[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", userId)
	}
}
