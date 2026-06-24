-- +goose Up
-- +goose StatementBegin
INSERT INTO permissions (code, description) VALUES ('reports.templates.access', 'Acceso al reporte de envíos de plantillas (templates)');
INSERT INTO role_permissions (role_id, permission_id) 
SELECT 1, id FROM permissions WHERE code = 'reports.templates.access';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions WHERE permission_id = (SELECT id FROM permissions WHERE code = 'reports.templates.access');
DELETE FROM permissions WHERE code = 'reports.templates.access';
-- +goose StatementEnd
