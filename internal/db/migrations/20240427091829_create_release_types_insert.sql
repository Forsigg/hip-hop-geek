-- +goose Up
-- +goose StatementBegin
INSERT INTO release_types (id, type)
VALUES 
    (1, "Album"),
    (2, "Single");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM release_types;
-- +goose StatementEnd
