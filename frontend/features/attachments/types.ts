export interface Attachment {
  id: string;

  entity_type: "LEAD" | "CONTACT" | "TASK";

  entity_id: string;

  file_name: string;

  file_url: string;

  public_id: string;

  resource_type: string;

  file_size: number;

  uploaded_by: string;

  created_at: string;
}

export interface UploadAttachmentPayload {
  file: File;
}