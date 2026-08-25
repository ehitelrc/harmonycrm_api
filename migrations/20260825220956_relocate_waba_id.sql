-- +goose Up
-- 1. Ensure meta_waba_id column exists in channel_integrations
ALTER TABLE public.channel_integrations ADD COLUMN IF NOT EXISTS meta_waba_id text;

-- 2. Update vw_channel_integrations view to include meta_waba_id
CREATE OR REPLACE VIEW public.vw_channel_integrations AS
 SELECT ci.id AS channel_integration_id,
    ci.company_id,
    ci.channel_id,
    ci.webhook_url,
    ci.access_token,
    ci.app_identifier,
    ci.meta_waba_id,
    ci.is_active,
    ci.created_at,
    ci.updated_at,
    ci.is_non_commercial,
    ci.integration_name,
    ci.department_id,
    d.name AS department_name,
    c.code AS channel_code,
    ci.analyze_incoming_images
   FROM ((public.channel_integrations ci
     JOIN public.channels c ON ((c.id = ci.channel_id)))
     LEFT JOIN public.departments d ON ((ci.department_id = d.id)));

-- 3. Update vw_case_channel_integration view to include meta_waba_id
CREATE OR REPLACE VIEW public.vw_case_channel_integration AS
 SELECT c.id AS case_id,
    c.company_id,
    c.channel_id,
    c.sender_id,
    c.status,
    c.started_at,
    c.closed_at,
    c.channel_integration_id,
    ch.name AS channel_name,
    ch.code AS channel_code,
    ci.webhook_url,
    ci.access_token,
    ci.app_identifier,
    ci.meta_waba_id,
    ci.is_active AS integration_is_active,
    ci.updated_at AS integration_updated_at
   FROM ((public.cases c
     JOIN public.channel_integrations ci ON ((ci.id = c.channel_integration_id)))
     JOIN public.channels ch ON ((ch.id = ci.channel_id)));

-- 4. Copy existing meta_waba_id values from channels to channel_integrations (Migrate data)
UPDATE public.channel_integrations ci
SET meta_waba_id = c.meta_waba_id
FROM public.channels c
WHERE ci.channel_id = c.id
  AND c.meta_waba_id IS NOT NULL
  AND ci.meta_waba_id IS NULL;


-- +goose Down
-- Revert view vw_case_channel_integration
CREATE OR REPLACE VIEW public.vw_case_channel_integration AS
 SELECT c.id AS case_id,
    c.company_id,
    c.channel_id,
    c.sender_id,
    c.status,
    c.started_at,
    c.closed_at,
    c.channel_integration_id,
    ch.name AS channel_name,
    ch.code AS channel_code,
    ci.webhook_url,
    ci.access_token,
    ci.app_identifier,
    ci.is_active AS integration_is_active,
    ci.updated_at AS integration_updated_at
   FROM ((public.cases c
     JOIN public.channel_integrations ci ON ((ci.id = c.channel_integration_id)))
     JOIN public.channels ch ON ((ch.id = ci.channel_id)));

-- Revert view vw_channel_integrations
CREATE OR REPLACE VIEW public.vw_channel_integrations AS
 SELECT ci.id AS channel_integration_id,
    ci.company_id,
    ci.channel_id,
    ci.webhook_url,
    ci.access_token,
    ci.app_identifier,
    ci.is_active,
    ci.created_at,
    ci.updated_at,
    ci.is_non_commercial,
    ci.integration_name,
    ci.department_id,
    d.name AS department_name,
    c.code AS channel_code,
    ci.analyze_incoming_images
   FROM ((public.channel_integrations ci
     JOIN public.channels c ON ((c.id = ci.channel_id)))
     LEFT JOIN public.departments d ON ((ci.department_id = d.id)));

-- Optionally, we keep the column in case of rollback, or we could drop it:
-- ALTER TABLE public.channel_integrations DROP COLUMN IF EXISTS meta_waba_id;
