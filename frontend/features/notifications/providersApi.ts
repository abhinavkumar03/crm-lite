import api from "@/services/api";

export type ProviderChannel = "email" | "whatsapp";

export interface CommunicationProvider {
  id: string;
  channel: ProviderChannel;
  provider_type: string;
  name: string;
  config: Record<string, unknown>;
  secrets_configured: boolean;
  is_default: boolean;
  is_active: boolean;
  last_health_at?: string | null;
  last_error?: string | null;
  created_at: string;
  updated_at: string;
}

export interface SenderIdentity {
  id: string;
  provider_id?: string | null;
  channel: ProviderChannel;
  display_name?: string | null;
  from_address: string;
  reply_to?: string | null;
  is_default: boolean;
  created_at: string;
  updated_at: string;
}

export async function listProviders(channel?: string): Promise<CommunicationProvider[]> {
  const res = await api.get("/communication-providers", {
    params: channel ? { channel } : undefined,
  });
  return res.data.data;
}

export async function createProvider(payload: {
  channel: ProviderChannel;
  provider_type: string;
  name: string;
  config?: Record<string, unknown>;
  secrets?: Record<string, unknown>;
  is_default?: boolean;
  is_active?: boolean;
}): Promise<CommunicationProvider> {
  const res = await api.post("/communication-providers", payload);
  return res.data.data;
}

export async function updateProvider(
  id: string,
  payload: Record<string, unknown>
): Promise<CommunicationProvider> {
  const res = await api.put(`/communication-providers/${id}`, payload);
  return res.data.data;
}

export async function deleteProvider(id: string): Promise<void> {
  await api.delete(`/communication-providers/${id}`);
}

export async function testProvider(id: string, to: string): Promise<void> {
  await api.post(`/communication-providers/${id}/test`, { to });
}

export async function listSenders(channel?: string): Promise<SenderIdentity[]> {
  const res = await api.get("/communication-senders", {
    params: channel ? { channel } : undefined,
  });
  return res.data.data;
}

export async function createSender(payload: {
  channel: ProviderChannel;
  from_address: string;
  display_name?: string;
  reply_to?: string;
  provider_id?: string;
  is_default?: boolean;
}): Promise<SenderIdentity> {
  const res = await api.post("/communication-senders", payload);
  return res.data.data;
}

export async function deleteSender(id: string): Promise<void> {
  await api.delete(`/communication-senders/${id}`);
}
