-- +goose Up

CREATE TABLE champions (
    id SERIAL PRIMARY KEY NOT NULL,
    champ_name TEXT NOT NULL UNIQUE
);

-- +goose Down

DROP TABLE champions;