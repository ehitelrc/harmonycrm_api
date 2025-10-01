// src/app/models/company-channel-template-view.model.ts
export interface CompanyChannelTemplateView {
  company_id: number;
  channel_integration_id: number;
  webhook_url: string;
  access_token?: string;
  app_identifier?: string;
  integration_active: boolean;
  integration_created_at: string;
  integration_updated_at: string;

  channel_id: number;
  channel_code: string;
  channel_name: string;
  channel_description?: string;

  template_id?: number;
  template_name?: string;
  language?: string;
  template_active?: boolean;
}