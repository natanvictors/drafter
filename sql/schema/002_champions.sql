-- +goose Up

CREATE TABLE champions (
    ID SERIAL PRIMARY KEY,
    champ_name TEXT NOT NULL,
    times_picked INTEGER DEFAULT 0,
    picked_top INTEGER DEFAULT 0,
    picked_jungle INTEGER DEFAULT 0,
    picked_mid INTEGER DEFAULT 0,
    picked_adc INTEGER DEFAULT 0,
    picked_sup INTEGER DEFAULT 0,
    wins_top INTEGER DEFAULT 0,
    losses_top INTEGER DEFAULT 0,
    wins_jg INTEGER DEFAULT 0,
    losses_jg INTEGER DEFAULT 0,
    wins_mid INTEGER DEFAULT 0,
    losses_mid INTEGER DEFAULT 0,
    wins_adc INTEGER DEFAULT 0,
    losses_adc INTEGER DEFAULT 0,
    wins_sup INTEGER DEFAULT 0,
    losses_sup INTEGER DEFAULT 0,
    ingame_role TEXT
);

-- +goose Down

DROP TABLE champions;