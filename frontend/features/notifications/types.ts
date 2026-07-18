// Types mirroring the backend notification pipeline (Phase 11).

export type NotificationChannel = "email" | "whatsapp";

export type NotificationStatus = "queued" | "sent" | "failed";

export interface Notification {
  id: string;
  channel: NotificationChannel;
  recipient: string;
  subject: string | null;
  body: string;
  template: string | null;
  data: Record<string, unknown>;
  status: NotificationStatus;
  provider: string | null;
  error: string | null;
  entity_type: string | null;
  entity_id: string | null;
  created_by: string | null;
  sent_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface SendNotificationRequest {
  channel: NotificationChannel;
  to: string;
  subject?: string;
  body: string;
  template?: string;
  data?: Record<string, unknown>;
  entity_type?: string;
  entity_id?: string;
}

export interface NotificationListParams {
  page?: number;
  page_size?: number;
  status?: NotificationStatus;
  channel?: NotificationChannel;
}

export interface NotificationListResult {
  notifications: Notification[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}
