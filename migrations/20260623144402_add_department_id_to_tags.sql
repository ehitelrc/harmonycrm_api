-- +goose Up
-- +goose StatementBegin
ALTER TABLE tags ADD COLUMN IF NOT EXISTS department_id INTEGER;
ALTER TABLE tags ADD CONSTRAINT fk_tag_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tags DROP CONSTRAINT IF EXISTS fk_tag_department;
ALTER TABLE tags DROP COLUMN IF EXISTS department_id;
-- +goose StatementEnd
