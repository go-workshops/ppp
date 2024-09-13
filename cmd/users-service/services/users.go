package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-workshops/ppp/pkg/tracing"
)

func NewUsers() *Users {
	return &Users{}
}

type Users struct {
}

func (s *Users) Register(ctx context.Context) (string, error) {
	// omitting the data layer in here for brevity
	_, span := tracing.StartPostgres(ctx, "register_user_txn")
	defer span.End()

	userID := uuid.New().String()
	attributes := trace.WithAttributes(attribute.String("user_id", userID))
	span.AddEvent("create user identity", attributes)
	time.Sleep(200 * time.Millisecond)
	// code for creating the user identity

	span.AddEvent("create user profile", attributes)
	time.Sleep(300 * time.Millisecond)
	// code for creating a user profile

	return userID, nil
}
