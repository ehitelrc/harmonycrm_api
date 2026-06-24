-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.messages ADD COLUMN template_id integer;
ALTER TABLE public.messages ADD CONSTRAINT fk_messages_template FOREIGN KEY (template_id) REFERENCES public.message_templates(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.messages DROP CONSTRAINT IF EXISTS fk_messages_template;
ALTER TABLE public.messages DROP COLUMN IF EXISTS template_id;
-- +goose StatementEnd
