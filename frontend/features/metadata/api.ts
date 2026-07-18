import api from "@/services/api";

import { cached, invalidateMetadataCache } from "./cache";
import {
  FormValues,
  ModuleField,
  ModuleSummary,
  RecordListParams,
  RecordListResult,
  RecordResponse,
  SavedView,
  SaveViewPayload,
  ValidateResult,
  ValidationSchema,
} from "./types";

export { invalidateMetadataCache };

export async function getModules(): Promise<ModuleSummary[]> {
  return cached("modules", async () => {
    const res = await api.get("/modules");
    return res.data.data;
  });
}

export async function getModuleFields(
  moduleId: string
): Promise<ModuleField[]> {
  return cached(`fields:${moduleId}`, async () => {
    const res = await api.get(`/modules/${moduleId}/fields`);
    return res.data.data;
  });
}

export async function getValidationSchema(
  moduleId: string
): Promise<ValidationSchema> {
  return cached(`validation:${moduleId}`, async () => {
    const res = await api.get(`/modules/${moduleId}/validation-schema`);
    return res.data.data;
  });
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

// --- Records (Phase 10 module runtime) -------------------------------------

export async function listRecords(
  moduleId: string,
  params: RecordListParams = {}
): Promise<RecordListResult> {
  const res = await api.get(`/modules/${moduleId}/records`, { params });
  return res.data.data;
}

export async function getRecord(
  moduleId: string,
  recordId: string,
  expand = true
): Promise<RecordResponse> {
  const res = await api.get(`/modules/${moduleId}/records/${recordId}`, {
    params: expand ? { expand: true } : undefined,
  });
  return res.data.data;
}

export async function createRecord(
  moduleId: string,
  data: FormValues,
  ownerId?: string
): Promise<RecordResponse> {
  const res = await api.post(`/modules/${moduleId}/records`, {
    data,
    owner_id: ownerId,
  });
  return res.data.data;
}

export async function updateRecord(
  moduleId: string,
  recordId: string,
  data: FormValues,
  ownerId?: string
): Promise<RecordResponse> {
  const res = await api.put(`/modules/${moduleId}/records/${recordId}`, {
    data,
    owner_id: ownerId,
  });
  return res.data.data;
}

export async function deleteRecord(
  moduleId: string,
  recordId: string
): Promise<void> {
  await api.delete(`/modules/${moduleId}/records/${recordId}`);
}
