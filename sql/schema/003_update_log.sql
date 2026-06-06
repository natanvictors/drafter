-- +goose Up

CREATE TABLE update_log (
    region TEXT PRIMARY KEY,
    last_updated_at TIMESTAMPTZ NOT NULL
);

-- +goose Down

DROP TABLE update_log;