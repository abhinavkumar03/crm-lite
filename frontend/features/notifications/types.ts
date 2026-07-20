// Types for the Enterprise Notification Center.

export type NotificationChannel = "email" | "whatsapp";

export type NotificationStatus =
  | "draft"
  | "scheduled"
  | "queued"
  | "processing"
  | "sent"
  | "delivered"
  | "opened"
  | "read"
  | "failed"
  | "retrying"
  | "cancelled";

export type ComposeMode = "draft" | "send" | "schedule";

export type TemplateCategory =
  | "sales"
  | "follow_up"
  | "welcome"
  | "proposal"
  | "invoice"
  | "reminder"
  | "quotation"
  | "marketing"
  | "support"
  | "custom";

export interface DeliveryEvent {
  id: string;
  event: string;
  provider?: string | null;
  payload?: Record<string, unknown>;
  created_at: string;
}

export interface Notification {
  id: string;
  channel: NotificationChannel;
  recipient: string;
  cc?: string[];
  bcc?: string[];
  subject: string | null;
  body: string;
  body_html?: string | null;
  template: string | null;
  template_id?: string | null;
  data: Record<string, unknown>;
  variables_used?: Record<string, unknown>;
  status: NotificationStatus;
  provider: string | null;
  error: string | null;
  last_error?: string | null;
  provider_response?: Record<string, unknown>;
  entity_type: string | null;
  entity_id: string | null;
  module_id?: string | null;
  attachment_ids?: string[];
  retry_count: number;
  max_retries: number;
  created_by: string | null;
  scheduled_at?: string | null;
  cancelled_at?: string | null;
  queued_at?: string | null;
  processing_at?: string | null;
  sent_at: string | null;
  delivered_at?: string | null;
  opened_at?: string | null;
  read_at?: string | null;
  created_at: string;
  updated_at: string;
  events?: DeliveryEvent[];
}

export interface ComposeNotificationRequest {
  mode?: ComposeMode;
  channel: NotificationChannel;
  to: string;
  cc?: string[];
  bcc?: string[];
  subject?: string;
  body: string;
  body_html?: string;
  template?: string;
  template_id?: string;
  data?: Record<string, unknown>;
  entity_type?: string;
  entity_id?: string;
  module_id?: string;
  attachment_ids?: string[];
  scheduled_at?: string;
  max_retries?: number;
}

/** @deprecated Use ComposeNotificationRequest */
export type SendNotificationRequest = ComposeNotificationRequest;

export interface NotificationListParams {
  page?: number;
  page_size?: number;
  status?: NotificationStatus | "";
  channel?: NotificationChannel | "";
  q?: string;
  module_id?: string;
  entity_id?: string;
  entity_type?: string;
  date_from?: string;
  date_to?: string;
  template_id?: string;
}

export interface NotificationListResult {
  notifications: Notification[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

export interface NotificationMetrics {
  emails_sent_today: number;
  whatsapp_sent_today: number;
  failed_count: number;
  scheduled_count: number;
  draft_count: number;
  total_sent: number;
  total_delivered: number;
  delivery_rate: number;
  open_rate: number;
  read_rate: number;
}

export interface NotificationTemplate {
  id: string;
  channel: NotificationChannel;
  name: string;
  category: TemplateCategory;
  subject: string | null;
  body: string;
  body_html?: string | null;
  variables: string[];
  is_active: boolean;
  status?: "draft" | "published";
  version?: number;
  created_by?: string | null;
  created_at: string;
  updated_at: string;
}

export interface TemplateListResult {
  templates: NotificationTemplate[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

export interface CreateTemplatePayload {
  channel: NotificationChannel;
  name: string;
  category?: TemplateCategory;
  subject?: string;
  body: string;
  body_html?: string;
  variables?: string[];
  is_active?: boolean;
}

export interface UpdateTemplatePayload {
  name?: string;
  category?: TemplateCategory;
  subject?: string;
  body?: string;
  body_html?: string;
  variables?: string[];
  is_active?: boolean;
}

export interface PreviewTemplatePayload {
  data?: Record<string, unknown>;
  module_id?: string;
  entity_id?: string;
  entity_type?: string;
}

export interface PreviewTemplateResult {
  subject: string;
  body: string;
  body_html?: string;
}
