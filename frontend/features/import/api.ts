import api from "@/services/api";

import {
  AnalyzeResult,
  ImportJob,
  ImportListParams,
  ImportListResult,
} from "./types";

// analyzeImport uploads the file for parsing only — nothing is persisted. The
// browser sets the multipart Content-Type (with boundary) automatically for a
// FormData body.
export async function analyzeImport(
  moduleId: string,
  file: File
): Promise<AnalyzeResult> {
  const form = new FormData();
  form.append("file", file);
  const res = await api.post(`/modules/${moduleId}/imports/analyze`, form);
  return res.data.data;
}

export async function createImport(
  moduleId: string,
  file: File,
  mapping: Record<string, string>,
  options?: Record<string, unknown>
): Promise<ImportJob> {
  const form = new FormData();
  form.append("file", file);
  form.append("mapping", JSON.stringify(mapping));
  if (options) form.append("options", JSON.stringify(options));
  const res = await api.post(`/modules/${moduleId}/imports`, form);
  return res.data.data;
}

export async function listImports(
  moduleId: string,
  params: ImportListParams = {}
): Promise<ImportListResult> {
  const res = await api.get(`/modules/${moduleId}/imports`, { params });
  return res.data.data;
}

export async function getImport(
  moduleId: string,
  importId: string
): Promise<ImportJob> {
  const res = await api.get(`/modules/${moduleId}/imports/${importId}`);
  return res.data.data;
}
