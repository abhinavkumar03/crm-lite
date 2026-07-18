import api from "@/services/api";

import {
  FormValues,
  ModuleField,
  ModuleSummary,
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
