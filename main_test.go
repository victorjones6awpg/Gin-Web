package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware_PreservesBody(t *testing.T) {
	// Setup Gin router
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Apply the authentication middleware
	r.Use(AuthMiddleware())

	// Define a handler that binds the JSON body
	r.POST("/test", func(c *gin.Context) {
		var payload TestPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"received": payload.Message})
	})

	// Perform request
	jsonBody := `{"message":"hello world"}`
	req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Response: %s", w.Code, w.Body.String())
	}

	expectedResponse := `{"received":"hello world"}`
	if !strings.Contains(w.Body.String(), expectedResponse) {
		t.Errorf("Expected response to contain %q, got %q", expectedResponse, w.Body.String())
	}
}
