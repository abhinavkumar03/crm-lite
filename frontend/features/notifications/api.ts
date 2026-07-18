import api from "@/services/api";

import {
  Notification,
  NotificationListParams,
  NotificationListResult,
  SendNotificationRequest,
} from "./types";

export async function sendNotification(
  payload: SendNotificationRequest
): Promise<Notification> {
  const res = await api.post("/notifications", payload);
  return res.data.data;
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
