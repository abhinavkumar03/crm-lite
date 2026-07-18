export interface LayoutSection {
  key: string;
  label: string;
  fields: string[];
}

export interface DetailLayoutConfig {
  sections?: LayoutSection[];
  tabs?: string[];
}

export interface DetailLayout {
  id: string;
  name: string;
  layout_type: string;
  is_default: boolean;
  config: DetailLayoutConfig;
}

export interface WorkspaceNote {
  id: string;
  title?: string | null;
  body: string;
  created_by: string;
  author_name: string;
  created_at: string;
  updated_at: string;
}

export interface WorkspaceAttachment {
  id: string;
  file_name: string;
  file_url: string;
  public_id: string;
  resource_type?: string | null;
  file_size?: number | null;
  uploaded_by: string;
  uploader_name: string;
  created_at: string;
}

export interface WorkspaceActivity {
  id: string;
  action: string;
  description: string;
  performed_by: string;
  actor_name: string;
  metadata?: Record<string, unknown>;
  created_at: string;
}

export interface RelatedDescriptor {
  child_module_id: string;
  child_module_name: string;
  child_api_name: string;
  lookup_field_api_name: string;
  lookup_field_label: string;
}
