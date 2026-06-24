-- +goose Up
-- +goose StatementBegin
INSERT INTO permissions (code, description) VALUES ('reports.cases.access', 'Acceso al reporte de gestiones generales (casos)');
INSERT INTO role_permissions (role_id, permission_id) 
SELECT 1, id FROM permissions WHERE code = 'reports.cases.access';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions WHERE permission_id = (SELECT id FROM permissions WHERE code = 'reports.cases.access');
DELETE FROM permissions WHERE code = 'reports.cases.access';
-- +goose StatementEnd
