import api from "@/services/api";
import { uploadFile } from "@/features/uploads/api";

import {
  DetailLayout,
  RelatedDescriptor,
  WorkspaceActivity,
  WorkspaceAttachment,
  WorkspaceNote,
} from "./types";

export async function getDetailLayout(moduleId: string): Promise<DetailLayout> {
  const res = await api.get(`/modules/${moduleId}/layouts/detail`);
  const data = res.data.data;
  // config may arrive as object already
  if (typeof data.config === "string") {
    try {
      data.config = JSON.parse(data.config);
    } catch {
      data.config = {};
    }
  }
  return data;
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
