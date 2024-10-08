package middleware

import (
	"fmt"
	"go-product/internal/core/port"
	"time"

	"github.com/gin-gonic/gin"
)

func SpeedTestMiddleware(speedTestService port.SpeedTestService) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		if err := speedTestService.WriteLog(c.Request.Method, c.Request.URL.Path, duration); err != nil {
			fmt.Printf("Unable to write log in MongoDB: %v\n", err)
		}
		fmt.Printf("Request %s %s took %v\n", c.Request.Method, c.Request.URL.Path, duration)
	}
}
