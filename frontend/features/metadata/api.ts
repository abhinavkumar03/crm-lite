import api from "@/services/api";

import {
  FormValues,
  ModuleField,
  ModuleSummary,
  SavedView,
  SaveViewPayload,
  ValidateResult,
  ValidationSchema,
} from "./types";

export async function getModules(): Promise<ModuleSummary[]> {
  const res = await api.get("/modules");
  return res.data.data;
}

export async function getModuleFields(
  moduleId: string
): Promise<ModuleField[]> {
  const res = await api.get(`/modules/${moduleId}/fields`);
  return res.data.data;
}

export async function getValidationSchema(
  moduleId: string
): Promise<ValidationSchema> {
  const res = await api.get(`/modules/${moduleId}/validation-schema`);
  return res.data.data;
}

export async function validateRecord(
  moduleId: string,
  data: FormValues
): Promise<ValidateResult> {
  const res = await api.post(`/modules/${moduleId}/validate`, { data });
  return res.data.data;
}

// --- Saved views -----------------------------------------------------------

export async function getViews(moduleId: string): Promise<SavedView[]> {
  const res = await api.get(`/modules/${moduleId}/views`);
  return res.data.data;
}

export async function createView(
  moduleId: string,
  payload: SaveViewPayload
): Promise<SavedView> {
  const res = await api.post(`/modules/${moduleId}/views`, payload);
  return res.data.data;
}

export async function updateView(
  moduleId: string,
  viewId: string,
  payload: Partial<SaveViewPayload>
): Promise<SavedView> {
  const res = await api.put(`/modules/${moduleId}/views/${viewId}`, payload);
  return res.data.data;
}

export async function deleteView(
  moduleId: string,
  viewId: string
): Promise<void> {
  await api.delete(`/modules/${moduleId}/views/${viewId}`);
}

export async function setDefaultView(
  moduleId: string,
  viewId: string
): Promise<void> {
  await api.post(`/modules/${moduleId}/views/${viewId}/default`);
}
