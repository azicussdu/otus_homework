package sqlstorage

import (
	"errors"
	"fmt"
	"github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for sqlx
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"time"
)

type SQLStorage struct {
	db *sqlx.DB
}

func New(dataSource string, migrateDir string) (*SQLStorage, error) {
	// dataSource := "postgres://username:password@localhost:5432/mydatabase"

	// Initialize sqlx with pgx/stdlib driver (high level connection)
	db, err := sqlx.Connect("pgx", dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sqlx database: %w", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	sqlStorage := &SQLStorage{db: db}

	// Run migrations
	if err = sqlStorage.runMigrations(migrateDir); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return sqlStorage, nil
}

func (s *SQLStorage) runMigrations(migrateDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(s.db.DB, migrateDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}

func (s *SQLStorage) Close() error {
	// Close the sqlx connection
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close sqlx connection: %w", err)
	}
	return nil
}

func (s *SQLStorage) Create(event storage.Event) (uint, error) {
	query := `
		INSERT INTO events (title, start_time, duration, description, user_id, notify_before)
		VALUES (:title, :start_time, :duration, :description, :user_id, :notify_before)
		RETURNING id`

	var eventID uint
	err := s.db.QueryRowx(query, event).Scan(&eventID)
	if err != nil {
		if isUniqueViolation(err) {
			return 0, storage.ErrDateBusy
		}
		return 0, err
	}
	return eventID, nil
}

func (s *SQLStorage) Update(eventID uint, event storage.Event) error {
	query := `
		UPDATE events
		SET title = :title, start_time = :start_time, duration = :duration, 
		    description = :description, user_id = :user_id, notify_before = :notify_before
		WHERE id = :id`
	event.ID = eventID
	_, err := s.db.NamedExec(query, event)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLStorage) Delete(eventID uint) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := s.db.Exec(query, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLStorage) ListEventsByDay(date time.Time) ([]storage.Event, error) {
	query := `
		SELECT id, title, start_time, duration, description, user_id, notify_before
		FROM events
		WHERE DATE(start_time) = $1`
	var events []storage.Event
	err := s.db.Select(&events, query, date) //dateOnly := date.Format("2006-01-02") ??
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *SQLStorage) ListEventsByWeek(startDate time.Time) ([]storage.Event, error) {
	query := `
		SELECT id, title, start_time, duration, description, user_id, notify_before
		FROM events
		WHERE start_time >= $1 AND start_time < $2`
	endDate := startDate.AddDate(0, 0, 7)
	var events []storage.Event
	err := s.db.Select(&events, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *SQLStorage) ListEventsByMonth(startDate time.Time) ([]storage.Event, error) {
	query := `
		SELECT id, title, start_time, duration, description, user_id, notify_before
		FROM events
		WHERE start_time >= $1 AND start_time < $2`
	endDate := startDate.AddDate(0, 1, 0)
	var events []storage.Event
	err := s.db.Select(&events, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// Utility function to check for unique constraint violations
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // PostgreSQL unique constraint violation code
	}
	return false
}
