// Types mirroring the backend metadata engine (Phases 5-7). These drive the
// dynamic form renderer.

export type FieldType =
  | "text"
  | "textarea"
  | "email"
  | "phone"
  | "number"
  | "currency"
  | "date"
  | "datetime"
  | "boolean"
  | "dropdown"
  | "multiselect"
  | "radio"
  | "checkbox"
  | "url"
  | "file"
  | "image"
  | "user"
  | "lookup"
  | "formula"
  | "json"
  | "richtext";

export type FieldOption = {
  label: string;
  value: string;
};

export type StorageDescriptor = {
  kind: "column" | "jsonb";
  path: string;
};

export interface ModuleField {
  id: string;
  module_id: string;
  api_name: string;
  label: string;
  field_type: FieldType;
  is_required: boolean;
  is_unique: boolean;
  is_read_only: boolean;
  default_value: string | null;
  placeholder: string | null;
  description: string | null;
  help_text: string | null;
  min_length: number | null;
  max_length: number | null;
  regex: string | null;
  validation_message: string | null;
  options: FieldOption[];
  lookup_module_id: string | null;
  sort_order: number;
  is_visible: boolean;
  is_searchable: boolean;
  is_filterable: boolean;
  is_nullable: boolean;
  is_indexed: boolean;
  is_system: boolean;
  storage: StorageDescriptor;
  created_at: string;
  updated_at: string;
}

export interface ModuleSummary {
  id: string;
  api_name: string;
  singular_label: string;
  plural_label: string;
  icon: string | null;
  color: string | null;
  storage_strategy: "native" | "dynamic";
  is_enabled: boolean;
  is_system: boolean;
  sort_order: number;
}

// Compiled per-field validation schema (GET /modules/:id/validation-schema).
export interface FieldSchema {
  api_name: string;
  label: string;
  type: FieldType;
  required: boolean;
  min_length?: number;
  max_length?: number;
  min?: number;
  max?: number;
  pattern?: string;
  format?: "email" | "url";
  options?: string[];
  multiple?: boolean;
  messages?: Record<string, string>;
}

export interface ValidationSchema {
  module_id: string;
  fields: FieldSchema[];
}

export interface FieldError {
  field: string;
  message: string;
}

export interface ValidateResult {
  valid: boolean;
  errors: FieldError[];
}

export type FieldValue = string | number | boolean | string[] | null;

export type FormValues = Record<string, FieldValue>;

// --- Conditional rendering -------------------------------------------------

export type ConditionOperator =
  | "equals"
  | "not_equals"
  | "in"
  | "not_in"
  | "empty"
  | "not_empty"
  | "truthy"
  | "falsy";

export interface Condition {
  field: string;
  operator: ConditionOperator;
  value?: unknown;
}

// A VisibilityRule shows or hides target fields when its condition is met.
export interface VisibilityRule {
  when: Condition;
  effect: "show" | "hide";
  targets: string[];
}

// --- Saved views / dynamic tables ------------------------------------------

export type SortDirection = "asc" | "desc";

export type FilterOperator =
  | "contains"
  | "equals"
  | "not_equals"
  | "gt"
  | "lt"
  | "in";

export interface ViewFilter {
  field: string;
  operator: FilterOperator;
  value: unknown;
}

export interface ViewSort {
  field: string;
  order: SortDirection | "";
}

export interface SavedView {
  id: string;
  module_id: string;
  name: string;
  columns: string[];
  filters: ViewFilter[];
  sort: ViewSort;
  is_default: boolean;
  is_public: boolean;
  owner_id: string | null;
  is_owner: boolean;
  created_at: string;
  updated_at: string;
}

export interface SaveViewPayload {
  name: string;
  columns: string[];
  filters: ViewFilter[];
  sort?: ViewSort;
  is_public?: boolean;
}

// A table row is just a bag of field values keyed by api_name.
export type TableRow = Record<string, FieldValue>;

// --- Record runtime (Phase 10) --------------------------------------------

export interface RelationRef {
  id: string;
  label: string;
}

export interface RecordResponse {
  id: string;
  module_id: string;
  data: Record<string, FieldValue>;
  owner_id: string | null;
  created_by: string | null;
  updated_by: string | null;
  relations?: Record<string, RelationRef>;
  created_at: string;
  updated_at: string;
}

export interface RecordListResult {
  records: RecordResponse[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

export interface RecordListParams {
  page?: number;
  page_size?: number;
  search?: string;
  sort?: string;
  order?: SortDirection;
  expand?: boolean;
}
