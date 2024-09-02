package services

import (
	"context"
	"encoding/json"
	"io"

	"github.com/go-workshops/ppp/cmd/simple-http/models"
	sharedContext "github.com/go-workshops/ppp/pkg/context"
)

type database interface {
	io.Writer
}

type Todo struct {
	db database
}

func NewTodo(db database) *Todo {
	return &Todo{
		db: db,
	}
}

func (t *Todo) CreateTodo(ctx context.Context, todo models.Todo) error {
	logger := sharedContext.Logger(ctx)

	err := json.NewEncoder(t.db).Encode(todo)
	if err != nil {
		logger.Error("could not write todo to the db")
		return err
	}

	return nil
}
