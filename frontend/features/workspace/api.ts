import api from "@/services/api";
import { uploadFile } from "@/features/uploads/api";

import {
  DetailLayout,
  FormLayout,
  LayoutSection,
  ListColumn,
  ListLayout,
  RelatedDescriptor,
  WorkspaceActivity,
  WorkspaceAttachment,
  WorkspaceNote,
} from "./types";

export async function getDetailLayout(moduleId: string): Promise<DetailLayout> {
  const res = await api.get(`/modules/${moduleId}/layouts/detail`);
  const data = res.data.data;
  if (typeof data.config === "string") {
    try {
      data.config = JSON.parse(data.config);
    } catch {
      data.config = {};
    }
  }
  return data;
}

export async function updateDetailLayout(
  moduleId: string,
  payload: {
    sections: LayoutSection[];
    tabs?: string[];
  }
): Promise<DetailLayout> {
  const res = await api.put(`/modules/${moduleId}/layouts/detail`, payload);
  const data = res.data.data;
  if (typeof data.config === "string") {
    try {
      data.config = JSON.parse(data.config);
    } catch {
      data.config = {};
    }
  }
  return data;
}

export async function getFormLayout(
  moduleId: string,
  mode: "create" | "edit" = "create"
): Promise<FormLayout> {
  const res = await api.get(`/modules/${moduleId}/layouts/form`, {
    params: { mode },
  });
  return res.data.data;
}

export async function updateFormLayout(
  moduleId: string,
  sections: LayoutSection[]
): Promise<FormLayout> {
  const res = await api.put(`/modules/${moduleId}/layouts/form`, { sections });
  return res.data.data;
}

export async function getListLayout(
  moduleId: string,
  opts?: { includeHidden?: boolean }
): Promise<ListLayout> {
  const res = await api.get(`/modules/${moduleId}/layouts/list`, {
    params: opts?.includeHidden ? { include_hidden: "true" } : undefined,
  });
  return res.data.data;
}

export async function updateListLayout(
  moduleId: string,
  columns: ListColumn[]
): Promise<ListLayout> {
  const res = await api.put(`/modules/${moduleId}/layouts/list`, { columns });
  return res.data.data;
}

export async function reorderListColumns(
  moduleId: string,
  columns: { field_key: string; order: number }[]
): Promise<ListLayout> {
  const res = await api.put(`/modules/${moduleId}/layouts/list/reorder`, {
    columns,
  });
  return res.data.data;
}

export async function toggleListColumn(
  moduleId: string,
  fieldKey: string,
  visible: boolean
): Promise<ListLayout> {
  const res = await api.put(`/modules/${moduleId}/layouts/list/toggle`, {
    field_key: fieldKey,
    visible,
  });
  return res.data.data;
}

export async function resetListLayout(moduleId: string): Promise<ListLayout> {
  const res = await api.post(`/modules/${moduleId}/layouts/list/reset`);
  return res.data.data;
}

export async function listRecordNotes(
  moduleId: string,
  recordId: string
): Promise<WorkspaceNote[]> {
  const res = await api.get(`/modules/${moduleId}/records/${recordId}/notes`);
  return res.data.data ?? [];
}

export async function createRecordNote(
  moduleId: string,
  recordId: string,
  body: string,
  title?: string
): Promise<WorkspaceNote> {
  const res = await api.post(`/modules/${moduleId}/records/${recordId}/notes`, {
    body,
    title: title || null,
  });
  return res.data.data;
}

export async function deleteRecordNote(
  moduleId: string,
  recordId: string,
  noteId: string
): Promise<void> {
  await api.delete(`/modules/${moduleId}/records/${recordId}/notes/${noteId}`);
}

export async function listRecordAttachments(
  moduleId: string,
  recordId: string
): Promise<WorkspaceAttachment[]> {
  const res = await api.get(
    `/modules/${moduleId}/records/${recordId}/attachments`
  );
  return res.data.data ?? [];
}

export async function uploadRecordAttachment(
  moduleId: string,
  recordId: string,
  file: File
): Promise<WorkspaceAttachment> {
  const uploaded = await uploadFile(file);
  const res = await api.post(
    `/modules/${moduleId}/records/${recordId}/attachments`,
    {
      file_name: file.name,
      file_url: uploaded.url,
      public_id: uploaded.public_id,
      resource_type: uploaded.resource_type,
      file_size: uploaded.bytes,
    }
  );
  return res.data.data;
}

export async function deleteRecordAttachment(
  moduleId: string,
  recordId: string,
  attachmentId: string
): Promise<void> {
  await api.delete(
    `/modules/${moduleId}/records/${recordId}/attachments/${attachmentId}`
  );
}

export async function listRecordActivities(
  moduleId: string,
  recordId: string
): Promise<WorkspaceActivity[]> {
  const res = await api.get(
    `/modules/${moduleId}/records/${recordId}/activities`
  );
  return res.data.data ?? [];
}

export async function listRelatedDescriptors(
  moduleId: string
): Promise<RelatedDescriptor[]> {
  const res = await api.get(`/modules/${moduleId}/related`);
  return res.data.data ?? [];
}
