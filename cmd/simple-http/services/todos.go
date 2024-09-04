package services

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/go-workshops/ppp/cmd/simple-http/models"
	sharedContext "github.com/go-workshops/ppp/pkg/context"
	"github.com/go-workshops/ppp/pkg/db"
)

type database interface {
	File(name string) db.FS
}

type Todo struct {
	db database
}

func NewTodo(db database) *Todo {
	return &Todo{
		db: db,
	}
}

func (s *Todo) CreateTodo(ctx context.Context, todo models.Todo) error {
	logger := sharedContext.Logger(ctx)

	err := json.NewEncoder(s.db.File(todo.ID + ".json")).Encode(todo)
	if err != nil {
		logger.Error("could not write todo to the db")
		return err
	}

	return nil
}

func (s *Todo) UpdateTodo(ctx context.Context, todo models.Todo) error {
	logger := sharedContext.Logger(ctx)

	var oldTodo models.Todo
	err := json.NewDecoder(s.db.File(todo.ID + ".json")).Decode(&oldTodo)
	if err != nil {
		logger.Error("could not read todo from the db", zap.Error(err))
		return err
	}

	newTodo := models.Todo{
		ID:          oldTodo.ID,
		Title:       oldTodo.Title,
		Description: oldTodo.Description,
	}
	if todo.Title != "" {
		newTodo.Title = todo.Title
	}
	if todo.Description != "" {
		newTodo.Description = todo.Description
	}

	return s.CreateTodo(ctx, newTodo)
}
