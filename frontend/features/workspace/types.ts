export interface LayoutSection {
  key: string;
  label: string;
  description?: string;
  order?: number;
  collapsed?: boolean;
  columns?: number;
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

/** Hydrated form layout from GET /layouts/form */
export interface FormLayoutField {
  id: string;
  key: string;
  label: string;
  type: string;
  required: boolean;
  editable: boolean;
  locked: boolean;
  display_order: number;
  placeholder?: string | null;
  description?: string | null;
  default_value?: string | null;
  lock_mode?: string;
  lookup_module_id?: string | null;
  options?: { label: string; value: string }[];
}

export interface FormLayoutSection {
  id: string;
  title: string;
  description?: string;
  order: number;
  collapsed: boolean;
  columns: number;
  fields: FormLayoutField[];
}

export interface FormLayout {
  id: string;
  name: string;
  layout_type: string;
  is_default: boolean;
  mode: string;
  sections: FormLayoutSection[];
}

export interface ListColumn {
  field_key: string;
  field_id?: string;
  label?: string;
  visible: boolean;
  order: number;
  sortable: boolean;
  searchable: boolean;
  system: boolean;
  locked?: boolean;
  width?: number | null;
}

export interface ListLayout {
  id: string;
  name: string;
  layout_type: string;
  is_default: boolean;
  columns: ListColumn[];
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
