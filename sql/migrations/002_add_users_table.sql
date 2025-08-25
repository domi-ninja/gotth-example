-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,

    display_name TEXT
);

INSERT INTO users ( id, email, password_hash, display_name ) VALUES(
    "fakeuser-id-123",
    "fake@fake.fake",
    "fake-password-hash",
    "Fakey Mc Fakeface"
);

-- +goose Down
DROP TABLE users;
