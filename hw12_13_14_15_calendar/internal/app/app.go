package app

import (
	"context"
)

type App struct {
	// TODO
}

type Logger interface { // TODO
}

type Storage interface { // TODO
}

func New(logger Logger, storage Storage) *App { //nolint:revive
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error { //nolint:revive
	// TODO
	return nil
}

// TODO
