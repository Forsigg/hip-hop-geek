-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS release_types (
    id INTEGER PRIMARY KEY,
    type TEXT UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS release_types;
-- +goose StatementEnd
