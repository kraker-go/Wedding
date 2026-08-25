-- +goose Up

CREATE TABLE changeuser(
    id SERIAL PRIMARY KEY,
    guest_id INTEGER NOT NULL,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down

DROP TABLE changeuser;