package storage

import "time"

//nolint:tagliatelle
type Notification struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"event_date"`
	Recipient uint      `json:"recipient"`
}
