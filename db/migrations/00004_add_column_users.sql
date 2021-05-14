-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN etc jsonb NOT NULL DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS etc;
-- +goose StatementEnd