import api from "@/services/api";

import { invalidateMetadataCache } from "@/features/metadata/cache";
import { ModuleField } from "@/features/metadata/types";

import {
  CreateFieldPayload,
  CreateModulePayload,
  CreateRulePayload,
  ModuleDetail,
  UpdateFieldPayload,
  UpdateModulePayload,
  UpdateRulePayload,
  ValidationRule,
} from "./types";

// --- Modules ---------------------------------------------------------------

export async function listModules(): Promise<ModuleDetail[]> {
  const res = await api.get("/modules");
  return res.data.data;
}

export async function createModule(
  payload: CreateModulePayload
): Promise<ModuleDetail> {
  const res = await api.post("/modules", payload);
  invalidateMetadataCache();
  return res.data.data;
}

export async function updateModule(
  moduleId: string,
  payload: UpdateModulePayload
): Promise<ModuleDetail> {
  const res = await api.put(`/modules/${moduleId}`, payload);
  invalidateMetadataCache();
  return res.data.data;
}

export async function setModuleStatus(
  moduleId: string,
  enabled: boolean
): Promise<ModuleDetail> {
  const res = await api.patch(`/modules/${moduleId}/status`, { enabled });
  invalidateMetadataCache();
  return res.data.data;
}

export async function deleteModule(moduleId: string): Promise<void> {
  await api.delete(`/modules/${moduleId}`);
  invalidateMetadataCache();
}

// --- Fields ----------------------------------------------------------------

export async function listFields(moduleId: string): Promise<ModuleField[]> {
  const res = await api.get(`/modules/${moduleId}/fields`);
  return res.data.data;
}

export async function createField(
  moduleId: string,
  payload: CreateFieldPayload
): Promise<ModuleField> {
  const res = await api.post(`/modules/${moduleId}/fields`, payload);
  invalidateMetadataCache(`fields:${moduleId}`);
  invalidateMetadataCache("modules");
  return res.data.data;
}

export async function updateField(
  moduleId: string,
  fieldId: string,
  payload: UpdateFieldPayload
): Promise<ModuleField> {
  const res = await api.put(`/modules/${moduleId}/fields/${fieldId}`, payload);
  invalidateMetadataCache(`fields:${moduleId}`);
  invalidateMetadataCache(`validation:${moduleId}`);
  return res.data.data;
}

export async function deleteField(
  moduleId: string,
  fieldId: string
): Promise<void> {
  await api.delete(`/modules/${moduleId}/fields/${fieldId}`);
  invalidateMetadataCache(`fields:${moduleId}`);
  invalidateMetadataCache(`validation:${moduleId}`);
}

// --- Validation rules ------------------------------------------------------

export async function listRules(moduleId: string): Promise<ValidationRule[]> {
  const res = await api.get(`/modules/${moduleId}/validation-rules`);
  return res.data.data;
}

export async function createRule(
  moduleId: string,
  payload: CreateRulePayload
): Promise<ValidationRule> {
  const res = await api.post(`/modules/${moduleId}/validation-rules`, payload);
  invalidateMetadataCache(`validation:${moduleId}`);
  return res.data.data;
}

export async function updateRule(
  moduleId: string,
  ruleId: string,
  payload: UpdateRulePayload
): Promise<ValidationRule> {
  const res = await api.put(
    `/modules/${moduleId}/validation-rules/${ruleId}`,
    payload
  );
  invalidateMetadataCache(`validation:${moduleId}`);
  return res.data.data;
}

export async function deleteRule(
  moduleId: string,
  ruleId: string
): Promise<void> {
  await api.delete(`/modules/${moduleId}/validation-rules/${ruleId}`);
  invalidateMetadataCache(`validation:${moduleId}`);
}
