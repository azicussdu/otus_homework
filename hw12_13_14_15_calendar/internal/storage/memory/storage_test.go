package memorystorage

import (
	storage2 "github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestMemoryStorage_CreateAndGet(t *testing.T) {
	storage := New()

	event := storage2.Event{
		Title:     "Test Event 1",
		StartTime: time.Now(),
		Duration:  time.Hour,
		UserID:    1,
	}

	id, err := storage.Create(event)
	require.NoError(t, err)

	retrieved, ok := storage.events[id]
	require.True(t, ok)
	require.Equal(t, uint(1), id)
	assert.Equal(t, event.Title, retrieved.Title)
	assert.Equal(t, event.UserID, retrieved.UserID)
}

func TestMemoryStorage_CreateDuplicate(T *testing.T) {
	storage := New()

	event := storage2.Event{
		Title:     "Test Event 1",
		StartTime: time.Now(),
		Duration:  time.Hour,
		UserID:    1,
	}

	_, err := storage.Create(event)
	require.NoError(T, err)

	event.StartTime = event.StartTime.Add(time.Minute * 30)
	_, err = storage.Create(event)
	require.ErrorIs(T, err, storage2.ErrDateBusy)
}

func TestMemoryStorage_Delete(T *testing.T) {
	storage := New()

	event := storage2.Event{
		Title:     "Test Event 1",
		StartTime: time.Now(),
		Duration:  time.Hour,
		UserID:    1,
	}

	id, err := storage.Create(event)
	require.NoError(T, err)

	err = storage.Delete(id)
	require.NoError(T, err)

	_, ok := storage.events[id]
	require.False(T, ok)
}

func TestMemoryStorage_ListEventsByDay(t *testing.T) {
	storage := New()

	event1 := storage2.Event{
		Title:     "Test Event 1",
		StartTime: time.Now(),
		Duration:  time.Minute,
		UserID:    1,
	}
	event2 := storage2.Event{
		Title:     "Test Event 2",
		StartTime: time.Now().Add(time.Minute * 2),
		Duration:  time.Minute,
		UserID:    1,
	}
	event3 := storage2.Event{
		Title:     "Test Event 3",
		StartTime: time.Now().Add(2 * 24 * time.Hour),
		Duration:  time.Minute,
		UserID:    1,
	}

	_, err := storage.Create(event1)
	require.NoError(t, err)

	_, err = storage.Create(event2)
	require.NoError(t, err)

	_, err = storage.Create(event3)
	require.NoError(t, err)

	events, err := storage.ListEventsByDay(time.Now())
	require.NoError(t, err)
	require.Len(t, events, 2)
	assert.Equal(t, event1.Title, events[0].Title)
	assert.Equal(t, event2.Title, events[1].Title)
}

func TestMemoryStorage_ConcurrentAccess(t *testing.T) {
	storage := New()
	event := storage2.Event{
		Title:     "Test Event",
		StartTime: time.Now(),
		Duration:  time.Hour,
		UserID:    1,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, _ = storage.Create(event)
	}()

	go func() {
		defer wg.Done()
		_, _ = storage.ListEventsByDay(time.Now())
	}()

	wg.Wait()
}
