-- +goose Up
-- +goose StatementBegin
-- Insert 3 testing events into the events table
INSERT INTO events (title, start_time, duration, description, user_id, notify_before)
VALUES
    ('Test Event 1', '2024-12-10 10:00:00', '2 hours', 'Description for Event 1', 1, '15 minutes'),
    ('Test Event 2', '2024-12-11 14:00:00', '1 hour', 'Description for Event 2', 2, '30 minutes'),
    ('Test Event 3', '2024-12-12 18:00:00', '3 hours', NULL, 3, NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the testing events from the events table
DELETE FROM events
WHERE title IN ('Test Event 1', 'Test Event 2', 'Test Event 3');
-- +goose StatementEnd
