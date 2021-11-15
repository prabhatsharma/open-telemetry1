package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("otel1")

func Ping(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "Ping")
	span.SetAttributes(attribute.String("game", "ping"))
	defer span.End()

	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
