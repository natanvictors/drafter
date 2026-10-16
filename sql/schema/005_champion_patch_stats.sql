-- +goose Up 

CREATE TABLE champion_patch_stats(
    champ_id INTEGER NOT NULL REFERENCES champions(id),
    role TEXT NOT NULL,
    patch TEXT NOT NULL,
    wins INTEGER NOT NULL DEFAULT 0,
    losses INTEGER NOT NULL DEFAULT 0,
    times_picked INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (champ_id, role, patch)
);


-- +goose Down

DROP TABLE champion_patch_stats;