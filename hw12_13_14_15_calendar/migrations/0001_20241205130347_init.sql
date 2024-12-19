-- +goose Up
-- +goose StatementBegin
-- Create the events table
CREATE TABLE events (
                        id SERIAL PRIMARY KEY,                 -- Auto-incrementing primary key
                        title VARCHAR(255) NOT NULL,           -- Title of the event (required)
                        start_time TIMESTAMP NOT NULL,         -- Start time of the event
                        duration INTERVAL NOT NULL,            -- Duration of the event
                        description TEXT,                      -- Optional description
                        user_id INTEGER NOT NULL,              -- ID of the user (required)
                        notify_before INTERVAL,                -- Optional notification interval
                        created_at TIMESTAMP DEFAULT now(),    -- Automatically set creation timestamp
                        updated_at TIMESTAMP DEFAULT now()     -- Automatically set update timestamp
);

-- Create the notifications table
CREATE TABLE notifications (
                               id SERIAL PRIMARY KEY,                 -- Auto-incrementing primary key
                               title VARCHAR(255) NOT NULL,           -- Title of the notification (required)
                               start_date TIMESTAMP NOT NULL,         -- Start date of the notification
                               recipient INTEGER NOT NULL,            -- Recipient's user ID (required)
                               created_at TIMESTAMP DEFAULT now(),    -- Automatically set creation timestamp
                               updated_at TIMESTAMP DEFAULT now()     -- Automatically set update timestamp
);

-- Add triggers for updated_at on events
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add the trigger for updated_at column
CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON events
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add the trigger for updated_at column
CREATE TRIGGER set_notifications_updated_at
    BEFORE UPDATE ON notifications
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the events table
DROP TABLE IF EXISTS events;
-- Drop the notifications table
DROP TABLE IF EXISTS notifications;

-- Drop the trigger function (if applicable)
DROP FUNCTION IF EXISTS update_updated_at_column;
-- +goose StatementEnd
