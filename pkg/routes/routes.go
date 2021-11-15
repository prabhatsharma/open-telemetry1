package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/open-telemetry1/pkg/handlers"
)

// var tracer = otel.Tracer("otel1")

// SetRoutes sets up all gi HTTP API endpoints that can be called by front end
func SetRoutes(r *gin.Engine) {

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// meta service - healthz
	r.GET("/healthz", handlers.GetHealthz)
	r.GET("/", handlers.Greeting)

	r.GET("/api/ping", handlers.Ping)

}
