package storage

import (
	"errors"
	"time"
)

//nolint:tagliatelle
type Event struct {
	ID           uint           `json:"id"`
	Title        string         `json:"title"`
	StartTime    time.Time      `json:"start_time"`
	Duration     time.Duration  `json:"duration"`
	Description  *string        `json:"description,omitempty"`
	UserID       uint           `json:"user_id"`
	NotifyBefore *time.Duration `json:"notify_before,omitempty"`
}

func (e *Event) SetDescription(description string) {
	e.Description = &description
}

func (e *Event) SetNotifyBefore(duration time.Duration) {
	e.NotifyBefore = &duration
}

var (
	ErrDateBusy      = errors.New("this date is busy with another event")
	ErrEventNotFound = errors.New("requested event is not found")
)
