import api from "@/services/api";

import {
  CreateTemplatePayload,
  ExportJob,
  ExportListParams,
  ExportListResult,
  ExportSpec,
  ExportTemplate,
} from "./types";

// buildParams turns a spec into query params for the synchronous endpoint.
function specToParams(spec: ExportSpec): Record<string, string> {
  const params: Record<string, string> = { format: spec.format };
  if (spec.columns && spec.columns.length) params.columns = spec.columns.join(",");
  if (spec.search) params.search = spec.search;
  if (spec.sort) params.sort = spec.sort;
  if (spec.order) params.order = spec.order;
  if (spec.expand) params.expand = "true";
  if (spec.filters && spec.filters.length) params.filters = JSON.stringify(spec.filters);
  return params;
}

function filenameFromDisposition(header: unknown, fallback: string): string {
  if (typeof header === "string") {
    const match = /filename="?([^"]+)"?/.exec(header);
    if (match) return match[1];
  }
  return fallback;
}

// triggerBrowserDownload saves a blob to disk via a transient anchor element.
function triggerBrowserDownload(blob: Blob, filename: string) {
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  a.remove();
  window.URL.revokeObjectURL(url);
}

// exportNow streams a file synchronously and saves it in the browser.
export async function exportNow(moduleId: string, spec: ExportSpec): Promise<void> {
  const res = await api.get(`/modules/${moduleId}/export`, {
    params: specToParams(spec),
    responseType: "blob",
  });
  const filename = filenameFromDisposition(
    res.headers["content-disposition"],
    `export.${spec.format}`
  );
  triggerBrowserDownload(res.data as Blob, filename);
}

export async function createExport(moduleId: string, spec: ExportSpec): Promise<ExportJob> {
  const res = await api.post(`/modules/${moduleId}/exports`, spec);
  return res.data.data;
}

export async function listExports(
  moduleId: string,
  params: ExportListParams = {}
): Promise<ExportListResult> {
  const res = await api.get(`/modules/${moduleId}/exports`, { params });
  return res.data.data;
}

// downloadExport fetches a completed job's stored file and saves it.
export async function downloadExport(moduleId: string, job: ExportJob): Promise<void> {
  const res = await api.get(`/modules/${moduleId}/exports/${job.id}/download`, {
    responseType: "blob",
  });
  const filename = filenameFromDisposition(
    res.headers["content-disposition"],
    job.filename
  );
  triggerBrowserDownload(res.data as Blob, filename);
}

// --- Templates -------------------------------------------------------------

export async function listTemplates(moduleId: string): Promise<ExportTemplate[]> {
  const res = await api.get(`/modules/${moduleId}/export-templates`);
  return res.data.data;
}

export async function createTemplate(
  moduleId: string,
  payload: CreateTemplatePayload
): Promise<ExportTemplate> {
  const res = await api.post(`/modules/${moduleId}/export-templates`, payload);
  return res.data.data;
}

export async function deleteTemplate(moduleId: string, templateId: string): Promise<void> {
  await api.delete(`/modules/${moduleId}/export-templates/${templateId}`);
}
