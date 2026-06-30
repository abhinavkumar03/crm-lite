import api from "@/services/api";
import { CreateContactPayload } from "./types";

export async function getContacts(params: {
  page?: number;
  limit?: number;
  search?: string;
  sort_by?: string;
  sort_order?: "asc" | "desc";
}) {
  const response = await api.get("/contacts", {
    params,
  });

  return response.data;
}

export async function createContact(
    payload: CreateContactPayload
) {
    return (await api.post("/contacts", payload)).data;
}

export async function updateContact(
    id: string,
    payload: CreateContactPayload
) {
    return (
        await api.put(`/contacts/${id}`, payload)
    ).data;
}

export async function deleteContact(
    id: string
) {
    return (
        await api.delete(`/contacts/${id}`)
    ).data;
}