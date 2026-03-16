-- +goose Up
-- +goose StatementBegin
CREATE TABLE todos (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    created INTEGER NOT NULL,
    completed INTEGER,
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE todos;
-- +goose StatementEnd
