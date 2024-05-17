-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    today_subscribe BOOLEAN NOT NULL,
    releases_message_id INTEGER,
    releases_page_count INTEGER,
    today_releases_message_id INTEGER,
    today_releases_page_count INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
