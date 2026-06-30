import api from "@/services/api";

export async function getLeads(params: {
  page?: number;
  limit?: number;
  search?: string;
  status?: string;
  sort_by?: string;
  sort_order?: "asc" | "desc";
}) {
  const response = await api.get("/leads", {
    params,
  });

  return response.data;
}

export interface CreateLeadPayload {
    name: string;
    email: string;
    phone: string;
    company: string;
    status: string;
    notes: string;
}

export async function createLead(
    payload: CreateLeadPayload
) {
    const response = await api.post(
        "/leads",
        payload
    );

    return response.data;
}

export async function updateLead(
    id: string,
    payload: CreateLeadPayload
) {
    const response = await api.put(
        `/leads/${id}`,
        payload
    );

    return response.data;
}

export async function deleteLead(
    id: string
) {

    const response = await api.delete(
        `/leads/${id}`
    );

    return response.data;
}