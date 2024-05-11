-- +goose Up
-- +goose StatementBegin
CREATE TABLE queue (
    id SERIAL PRIMARY KEY,
    chatid INTEGER NOT NULL,
    city VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE queue;
-- +goose StatementEnd
