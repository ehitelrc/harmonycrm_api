-- +goose Up
-- +goose StatementBegin
INSERT INTO permissions (code, description) VALUES ('reports.payment_validations.access', 'Acceso al reporte de validaciones de pagos');
INSERT INTO permissions (code, description) VALUES ('reports.unreconciled_payments.access', 'Acceso al reporte de pagos ERP sin vincular');
INSERT INTO permissions (code, description) VALUES ('reports.convertions.access', 'Acceso al reporte de conversiones');

INSERT INTO role_permissions (role_id, permission_id) 
SELECT 1, id FROM permissions WHERE code IN ('reports.payment_validations.access', 'reports.unreconciled_payments.access', 'reports.convertions.access');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions WHERE permission_id IN (SELECT id FROM permissions WHERE code IN ('reports.payment_validations.access', 'reports.unreconciled_payments.access', 'reports.convertions.access'));
DELETE FROM permissions WHERE code IN ('reports.payment_validations.access', 'reports.unreconciled_payments.access', 'reports.convertions.access');
-- +goose StatementEnd
