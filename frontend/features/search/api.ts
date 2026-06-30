import api from "@/services/api";

import { SearchResponse } from "./types";

export async function searchCRM(
  query: string
): Promise<SearchResponse> {
  const res = await api.get(
    "/search",
    {
      params: {
        q: query,
      },
    }
  );

  return res.data.data;
}