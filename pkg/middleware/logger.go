package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger manages the security by validating the token
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		verb := c.Request.Method
		path := c.Request.RequestURI

		// Process request
		c.Next()

		var size int
		if c.Writer != nil {
			size = c.Writer.Size()
		}
		elapsed := time.Since(t)

		fmt.Printf("time: %v\npath: %s\nverb: %s\nresponse size: %d\nelapsed time: %v", t, path, verb, size, elapsed)
	}
}
