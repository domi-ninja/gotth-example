-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    title TEXT NOT NULL,
    body TEXT NOT NULL,

    author UUID REFERENCES users
);

INSERT INTO posts (id, created_at, updated_at, title, body, author) VALUES (
    '123e4567-whatever',
    '2021-01-01 00:00:00',
    '2021-01-01 00:00:00',
    'Hello, world!',
    'This is a test post.',
    null);   


-- +goose Down
DROP TABLE posts;
