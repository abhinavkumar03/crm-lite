import api from "@/services/api";
import { CreateTaskPayload } from "./types";

export async function getTasks(params: {
  page?: number;
  limit?: number;
  search?: string;
  status?: string;
  sort_by?: string;
  sort_order?: "asc" | "desc";
}) {
  const response = await api.get("/tasks", {
    params,
  });

  return response.data;
}

export async function createTask(
    payload: CreateTaskPayload
) {
    return (await api.post("/tasks", payload)).data;
}

export async function updateTask(
    id: string,
    payload: CreateTaskPayload
) {
    return (
        await api.put(`/tasks/${id}`, payload)
    ).data;
}

export async function deleteTask(
    id: string
) {
    return (
        await api.delete(`/tasks/${id}`)
    ).data;
}