-- +goose Up
-- +goose StatementBegin
ALTER TABLE queue
ALTER COLUMN chatid TYPE BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE queue
ALTER COLUMN chatid TYPE INT;
-- +goose StatementEnd
