-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS releases (
    release_id INTEGER PRIMARY KEY,
    artist_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    release_type INTEGER NOT NULL,
    out_year INTEGER NOT NULL,
    out_month INTEGER NOT NULL,
    out_day INTEGER,
    cover_url TEXT,
    FOREIGN KEY (artist_id)
        REFERENCES artists (id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    FOREIGN KEY (release_type)
        REFERENCES release_types (id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS releases;
-- +goose StatementEnd
