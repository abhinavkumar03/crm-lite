import api from "@/services/api";

import type {
  BuilderMetadata,
  CreateWorkflowPayload,
  ExecutionDetail,
  ListExecutionsResult,
  ListWorkflowsResult,
  UpdateWorkflowPayload,
  VersionSummary,
  WorkflowDetail,
  WorkflowMetrics,
  WorkflowTemplate,
} from "./types";

export async function listWorkflows(params: {
  page?: number;
  page_size?: number;
  status?: string;
  module_id?: string;
} = {}): Promise<ListWorkflowsResult> {
  const res = await api.get("/workflows", { params });
  return res.data.data;
}

export async function getWorkflow(id: string): Promise<WorkflowDetail> {
  const res = await api.get(`/workflows/${id}`);
  return res.data.data;
}

export async function createWorkflow(
  payload: CreateWorkflowPayload
): Promise<WorkflowDetail> {
  const res = await api.post("/workflows", payload);
  return res.data.data;
}

export async function updateWorkflow(
  id: string,
  payload: UpdateWorkflowPayload
): Promise<WorkflowDetail> {
  const res = await api.patch(`/workflows/${id}`, payload);
  return res.data.data;
}

export async function archiveWorkflow(id: string): Promise<void> {
  await api.delete(`/workflows/${id}`);
}

export async function publishWorkflow(
  id: string,
  changelog = ""
): Promise<WorkflowDetail> {
  const res = await api.post(`/workflows/${id}/publish`, { changelog });
  return res.data.data;
}

export async function disableWorkflow(id: string): Promise<WorkflowDetail> {
  const res = await api.post(`/workflows/${id}/disable`);
  return res.data.data;
}

export async function listWorkflowVersions(
  id: string
): Promise<VersionSummary[]> {
  const res = await api.get(`/workflows/${id}/versions`);
  return res.data.data;
}

export async function rollbackWorkflow(
  id: string,
  versionId: string
): Promise<WorkflowDetail> {
  const res = await api.post(`/workflows/${id}/versions/${versionId}/rollback`);
  return res.data.data;
}

export async function runWorkflow(
  id: string,
  recordId: string,
  moduleId?: string
): Promise<void> {
  await api.post(`/workflows/${id}/run`, {
    record_id: recordId,
    module_id: moduleId || undefined,
  });
}

export async function getBuilderMetadata(): Promise<BuilderMetadata> {
  const res = await api.get("/workflows/builder-metadata");
  return res.data.data;
}

export async function listExecutions(params: {
  page?: number;
  page_size?: number;
  workflow_id?: string;
  module_id?: string;
  record_id?: string;
  status?: string;
} = {}): Promise<ListExecutionsResult> {
  const res = await api.get("/workflows/executions", { params });
  return res.data.data;
}

export async function getExecution(id: string): Promise<ExecutionDetail> {
  const res = await api.get(`/workflows/executions/${id}`);
  return res.data.data;
}

export async function retryExecution(id: string): Promise<void> {
  await api.post(`/workflows/executions/${id}/retry`);
}

export async function listWorkflowTemplates(): Promise<WorkflowTemplate[]> {
  const res = await api.get("/workflows/templates");
  return res.data.data;
}

export async function cloneWorkflowTemplate(
  id: string
): Promise<WorkflowDetail> {
  const res = await api.post(`/workflows/templates/${id}/clone`);
  return res.data.data;
}

export async function getWorkflowMetrics(): Promise<WorkflowMetrics> {
  const res = await api.get("/workflows/metrics");
  return res.data.data;
}
