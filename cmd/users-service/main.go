package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/go-workshops/ppp/cmd/users-service/clients"
	"github.com/go-workshops/ppp/cmd/users-service/routes"
	"github.com/go-workshops/ppp/cmd/users-service/services"
	"github.com/go-workshops/ppp/pkg/logging"
	"github.com/go-workshops/ppp/pkg/tracing"
)

func main() {
	ctx := context.Background()

	err := logging.Init(logging.Config{
		LoggingLevel:  "debug",
		LoggingOutput: []string{"stdout", "app.log"},
	})
	if err != nil {
		log.Fatalf("could not initialize logger: %v", err)
	}

	se, err := tracing.NewOTLPExporter("localhost:4317", 5*time.Second)
	if err != nil {
		log.Fatalf("could not initialize OTLP exporter: %v", err)
	}
	cfg := tracing.TracerProviderConfig{
		TracingEnabled: true,
		SpanExporter:   se,
		ServiceName:    "users-service",
		BatchTimeout:   30 * time.Second,
		ExportTimeout:  5 * time.Second,
		MaxBatchSize:   512,
		MaxQueueSize:   2048,
	}
	provider, err := tracing.NewTracerProvider(cfg)
	if err != nil {
		log.Fatalf("could not initialize tracing provider: %v", err)
	}
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(tracing.NewTextMapPropagator(ctx))

	notificationClient := clients.NewNotification("http://localhost:8002")
	usersService := services.NewUsers()
	routerCfg := routes.Config{
		UsersService:       usersService,
		NotificationClient: notificationClient,
	}

	log.Fatalln(http.ListenAndServe(":8001", routes.NewRouter(routerCfg)))
}
