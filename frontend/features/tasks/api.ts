import api from "@/services/api";
import { CreateTaskPayload } from "./types";

export async function getTasks(
    page = 1,
    search = "",
    status = ""
) {
    const response = await api.get("/tasks", {
        params: {
            page,
            search,
            status,
        },
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