-- +goose Up
-- +goose StatementBegin
ALTER TABLE conversations
ALTER COLUMN user1_chatidid TYPE BIGINT;
ALTER TABLE conversations
ALTER COLUMN user2_chatid TYPE BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE conversations
ALTER COLUMN user1_chatidid TYPE INT;
ALTER TABLE conversations
ALTER COLUMN user2_chatid TYPE INT;
-- +goose StatementEnd
