-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_queue_city ON queue (city);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_queue_city;
-- +goose StatementEnd
