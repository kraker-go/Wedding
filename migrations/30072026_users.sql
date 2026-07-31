-- +goose Up
CREATE TABLE guests (
                        id           SERIAL PRIMARY KEY,
                        firstname    TEXT NOT NULL,
                        lastname   TEXT NOT NULL,
                        created_at   TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE guests;

SELECT * FROM guests;

DELETE FROM guests WHERE id = 2;