package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pradytpk/go-ecommerce/tokens"
)

// Authentication Middleware of the program
//
//	@return gin.HandlerFunc
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "no authorizatoin header provided",
			})
			c.Abort()
			return
		}
		claims, err := tokens.ValidateToken(ClientToken)
		if err == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("uid", claims.UID)
		c.Next()
	}
}
