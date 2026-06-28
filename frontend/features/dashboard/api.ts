import api from "@/services/api";

import { DashboardResponse } from "./types";

export async function getDashboard() {

    const response =
        await api.get<{
            data: DashboardResponse;
        }>("/dashboard");

    return response.data.data;
}