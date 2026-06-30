import api from "@/services/api";

import { DashboardResponse } from "./types";

export async function getDashboard(
  refresh = false
) {
  const response = await api.get<{
    data: DashboardResponse;
  }>("/dashboard", {
    params: {
      refresh,
    },
  });

  return response.data.data;
}