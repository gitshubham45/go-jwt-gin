package middleware

import (
	helper "goJwt/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			c.Abort()
			return
		}
		c.Set("email",claims.Email)
		c.Set("first_name",claims.First_name)
		c.Set("last_name",claims.Last_name)
		c.Set("user_type",claims.User_type)
		c.Set("uid",claims.Uid)
		c.Next()
	}
}
