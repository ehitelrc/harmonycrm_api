-- +goose Up
-- +goose StatementBegin
ALTER TABLE cases ADD COLUMN IF NOT EXISTS closed_in_bulk BOOLEAN DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cases DROP COLUMN closed_in_bulk;
-- +goose StatementEnd
