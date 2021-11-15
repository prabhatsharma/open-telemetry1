package main

import (
	"context"
	"fmt"
	"os"

	zerolog "github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/prabhatsharma/open-telemetry1/pkg/routes"
	// "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			zerolog.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	err := godotenv.Load()
	if err != nil {
		zerolog.Print("Error loading .env file")
	}

	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("otel1"))

	routes.SetRoutes(r) // Set up all API routes.

	// Run the server

	if os.Getenv("PORT") != "" {
		r.Run(":" + os.Getenv("PORT"))
	} else {
		r.Run(":" + "9876")
	}
}

func initTracer() *sdktrace.TracerProvider {
	// stdoutExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	OTEL_OTLP_GRPC_ENDPOINT := os.Getenv("OTEL_OTLP_GRPC_ENDPOINT")

	if OTEL_OTLP_GRPC_ENDPOINT == "" {
		OTEL_OTLP_GRPC_ENDPOINT = "localhost:4317"
	}

	OTEL_OTLP_HTTP_ENDPOINT := os.Getenv("OTEL_OTLP_HTTP_ENDPOINT")

	if OTEL_OTLP_HTTP_ENDPOINT == "" {
		OTEL_OTLP_GRPC_ENDPOINT = "localhost:55681"
	}

	otlpExporter, err := otlptracegrpc.New(context.TODO(),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(OTEL_OTLP_GRPC_ENDPOINT),
	)

	if err != nil {
		fmt.Println("Error creating GRPC OTLP exporter: ", err)
	}

	otlpHTTPExporter, err := otlptracehttp.New(context.TODO(),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(OTEL_OTLP_HTTP_ENDPOINT),
	)

	if err != nil {
		fmt.Println("Error creating HTTP OTLP exporter: ", err)
	}

	idg := xray.NewIDGenerator()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		// sdktrace.WithBatcher(stdoutExporter),
		sdktrace.WithBatcher(otlpExporter),
		sdktrace.WithBatcher(otlpHTTPExporter),
		sdktrace.WithIDGenerator(idg),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
