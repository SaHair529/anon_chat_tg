-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY,
    user1_chatidid INT NOT NULL,
    user2_chatid INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE conversations;
-- +goose StatementEnd
