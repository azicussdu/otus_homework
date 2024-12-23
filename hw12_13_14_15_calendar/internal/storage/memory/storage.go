package memorystorage

import (
	"sync"
	"time"

	"github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type MemoryStorage struct {
	mu        sync.RWMutex
	events    map[uint]storage.Event
	idCounter uint
}

func New() *MemoryStorage {
	return &MemoryStorage{
		events: make(map[uint]storage.Event),
	}
}

func (s *MemoryStorage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear the events map
	s.events = nil
	// Reset the ID counter
	s.idCounter = 0
	return nil
}

func isDateBusy(events map[uint]storage.Event, event storage.Event) bool {
	for _, e := range events {
		if event.StartTime.After(e.StartTime) &&
			event.StartTime.Before(e.StartTime.Add(e.Duration)) &&
			e.UserID == event.UserID {
			return true
		}
	}

	return false
}

func (s *MemoryStorage) Create(event storage.Event) (uint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if isDateBusy(s.events, event) {
		return 0, storage.ErrDateBusy
	}

	s.idCounter++
	s.events[s.idCounter] = event

	return s.idCounter, nil
}

func (s *MemoryStorage) Update(eventID uint, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if isDateBusy(s.events, event) {
		return storage.ErrDateBusy
	}

	if _, ok := s.events[eventID]; !ok {
		return storage.ErrEventNotFound
	}

	s.events[eventID] = event
	return nil
}

func (s *MemoryStorage) Delete(eventID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[eventID]; !ok {
		return storage.ErrEventNotFound
	}

	delete(s.events, eventID)
	return nil
}

func (s *MemoryStorage) ListEventsByDay(startDate time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []storage.Event
	for _, event := range s.events {
		if event.StartTime.Year() == startDate.Year() &&
			event.StartTime.Month() == startDate.Month() &&
			event.StartTime.Day() == startDate.Day() {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *MemoryStorage) ListEventsByWeek(startDate time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startOfWeek := startDate.Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 7).Add(-time.Nanosecond)

	var events []storage.Event
	for _, event := range s.events {
		if event.StartTime.After(startOfWeek) &&
			event.StartTime.Before(endOfWeek) ||
			event.StartTime.Equal(startOfWeek) {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *MemoryStorage) ListEventsByMonth(startDate time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startOfMonth := startDate.Truncate(24 * time.Hour)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var events []storage.Event
	for _, event := range s.events {
		if event.StartTime.After(startOfMonth) &&
			event.StartTime.Before(endOfMonth) ||
			event.StartTime.Equal(startOfMonth) {
			events = append(events, event)
		}
	}

	return events, nil
}
