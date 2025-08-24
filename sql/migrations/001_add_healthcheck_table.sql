-- +goose Up
CREATE TABLE healthchecks (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE healthchecks;
