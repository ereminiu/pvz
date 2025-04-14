-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    user_id INT,
    action VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT(NOW()),
    error TEXT NOT NULL,
    attempts INT DEFAULT(3),
    status varchar(100) NOT NULL DEFAULT('CREATED'),
    created_at TIMESTAMP NOT NULL DEFAULT(NOW()),
    updated_at TIMESTAMP NOT NULL DEFAULT(NOW()),
    complited_at TIMESTAMP,
    processing_from TIMESTAMP NOT NULL DEFAULT(NOW())
);

INSERT INTO tasks (user_id, action, description, timestamp, error)
SELECT user_id, action, description, timestamp, error FROM logs;

TRUNCATE TABLE logs;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
-- +goose StatementEnd
