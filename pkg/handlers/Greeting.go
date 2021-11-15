package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Greeting function gets all events
func Greeting(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"otel1": "A sample implementation demonstrating auto instrumenting a go app using gin gonic framework. ",
	})
}
