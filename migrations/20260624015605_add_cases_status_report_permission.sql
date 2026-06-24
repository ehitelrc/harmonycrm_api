-- +goose Up
-- +goose StatementBegin
INSERT INTO permissions (code, description) VALUES ('reports.cases_status.access', 'Acceso al reporte granular de casos cerrados, abiertos o no respondidos');
INSERT INTO role_permissions (role_id, permission_id) 
SELECT 1, id FROM permissions WHERE code = 'reports.cases_status.access';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions WHERE permission_id = (SELECT id FROM permissions WHERE code = 'reports.cases_status.access');
DELETE FROM permissions WHERE code = 'reports.cases_status.access';
-- +goose StatementEnd
