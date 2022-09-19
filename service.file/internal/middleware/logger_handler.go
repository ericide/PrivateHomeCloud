package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"time"
)

func LoggerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtString := c.GetHeader("authorization")

		accessSecret := os.Getenv("ACCESS_SECRET")

		token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
			return []byte(accessSecret), nil
		})

		fmt.Println("token", err, token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "not valid",
			})
			c.Abort()
			return
		}

		if token.Valid != true {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "not valid",
			})
			c.Abort()
			return
		}

		t := time.Now()
		c.Next()
		latency := time.Since(t)
		fmt.Println(c.Request.URL)
		fmt.Printf("PATH: %v | USE TIME: %v | RESPONSE STATUS: %d\n", c.Request.Method+c.FullPath(), latency, c.Writer.Status())
	}
}
