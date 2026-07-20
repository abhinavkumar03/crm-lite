import api from "@/services/api";

import {
  ComposeNotificationRequest,
  CreateTemplatePayload,
  Notification,
  NotificationListParams,
  NotificationListResult,
  NotificationMetrics,
  NotificationTemplate,
  PreviewTemplatePayload,
  PreviewTemplateResult,
  TemplateListResult,
  UpdateTemplatePayload,
} from "./types";

export async function composeNotification(
  payload: ComposeNotificationRequest
): Promise<Notification> {
  const res = await api.post("/notifications", payload);
  return res.data.data;
}

/** @deprecated Prefer composeNotification */
export async function sendNotification(
  payload: ComposeNotificationRequest
): Promise<Notification> {
  return composeNotification({ ...payload, mode: payload.mode ?? "send" });
}

export async function listNotifications(
  params: NotificationListParams = {}
): Promise<NotificationListResult> {
  const res = await api.get("/notifications", { params });
  return res.data.data;
}

export async function getNotification(id: string): Promise<Notification> {
  const res = await api.get(`/notifications/${id}`);
  return res.data.data;
}

export async function retryNotification(id: string): Promise<Notification> {
  const res = await api.post(`/notifications/${id}/retry`);
  return res.data.data;
}

export async function cancelNotification(id: string): Promise<Notification> {
  const res = await api.post(`/notifications/${id}/cancel`);
  return res.data.data;
}

export async function getNotificationMetrics(): Promise<NotificationMetrics> {
  const res = await api.get("/notifications/metrics");
  return res.data.data;
}

export async function listTemplates(params: {
  page?: number;
  page_size?: number;
  channel?: string;
  category?: string;
} = {}): Promise<TemplateListResult> {
  const res = await api.get("/notification-templates", { params });
  return res.data.data;
}

export async function getTemplate(id: string): Promise<NotificationTemplate> {
  const res = await api.get(`/notification-templates/${id}`);
  return res.data.data;
}

export async function createTemplate(
  payload: CreateTemplatePayload
): Promise<NotificationTemplate> {
  const res = await api.post("/notification-templates", payload);
  return res.data.data;
}

export async function updateTemplate(
  id: string,
  payload: UpdateTemplatePayload
): Promise<NotificationTemplate> {
  const res = await api.put(`/notification-templates/${id}`, payload);
  return res.data.data;
}

export async function deleteTemplate(id: string): Promise<void> {
  await api.delete(`/notification-templates/${id}`);
}

export async function previewTemplate(
  id: string,
  payload: PreviewTemplatePayload = {}
): Promise<PreviewTemplateResult> {
  const res = await api.post(`/notification-templates/${id}/preview`, payload);
  return res.data.data;
}
