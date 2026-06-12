package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware intercepts, reads, and restores the request body stream.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var bodyBytes []byte
		if c.Request.Body != nil {
			var err error
			bodyBytes, err = io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
				return
			}
			// Restore the request body for downstream handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Example authentication logic
		if !isValidAuth(bodyBytes, c.Request.Header) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}

func isValidAuth(body []byte, headers http.Header) bool {
	// Simple validation logic for demonstration/testing.
	if headers.Get("Authorization") == "Bearer invalid-token" {
		return false
	}
	return true
}

type TestPayload struct {
	Message string `json:"message"`
}

func main() {
	r := gin.Default()

	r.Use(AuthMiddleware())

	r.POST("/test", func(c *gin.Context) {
		var payload TestPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"received": payload.Message})
	})

	fmt.Println("Starting server on :8080")
	r.Run(":8080")
}
