-- +goose Up
-- +goose StatementBegin
INSERT INTO permissions (code, description) VALUES ('MANAGE_TAGS', 'Gestión de tags globales para casos');
INSERT INTO role_permissions (role_id, permission_id) 
SELECT 1, id FROM permissions WHERE code = 'MANAGE_TAGS';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions WHERE permission_id = (SELECT id FROM permissions WHERE code = 'MANAGE_TAGS');
DELETE FROM permissions WHERE code = 'MANAGE_TAGS';
-- +goose StatementEnd
