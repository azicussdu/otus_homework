package app

import (
	"context"
	"time"

	"github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type Storage interface {
	Create(event storage.Event) (uint, error)
	Update(eventID uint, event storage.Event) error
	Delete(eventID uint) error
	ListEventsByDay(startDate time.Time) ([]storage.Event, error)
	ListEventsByWeek(startDate time.Time) ([]storage.Event, error)
	ListEventsByMonth(startDate time.Time) ([]storage.Event, error)
	Close() error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(_ context.Context, _, _ string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
