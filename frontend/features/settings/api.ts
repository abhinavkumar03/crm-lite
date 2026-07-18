import api from "@/services/api";

import { OrgSettings, UpdateSettingsPayload } from "./types";

// getSettings returns the current organization's settings, with backend defaults
// filled in for anything never saved.
export async function getSettings(): Promise<OrgSettings> {
  const res = await api.get("/settings");
  return res.data.data;
}

// updateSettings applies a partial change (any subset of name/general/automation).
export async function updateSettings(
  payload: UpdateSettingsPayload
): Promise<OrgSettings> {
  const res = await api.put("/settings", payload);
  return res.data.data;
}
