package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func LoggerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		fmt.Println(c.Request.URL)
		fmt.Printf("PATH: %v | USE TIME: %v | RESPONSE STATUS: %d\n", c.Request.Method+c.FullPath(), latency, c.Writer.Status())
	}
}
