-- +goose Up
-- +goose StatementBegin
-- Insert 3 testing notifications into the notifications table
INSERT INTO notifications (title, start_date, recipient)
VALUES
    ('Notification 1', '2024-12-15 08:00:00', 1),
    ('Notification 2', '2024-12-16 09:30:00', 2),
    ('Notification 3', '2024-12-17 14:45:00', 3);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the testing notifications from the notifications table
DELETE FROM notifications
WHERE title IN ('Notification 1', 'Notification 2', 'Notification 3');
-- +goose StatementEnd
