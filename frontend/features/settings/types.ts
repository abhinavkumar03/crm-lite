// Types mirroring the backend organization-settings engine (Phase 15).

export interface GeneralSettings {
  timezone: string;
  date_format: string;
  time_format: "12h" | "24h";
  currency: string;
  locale: string;
  week_start: "sunday" | "monday";
}

export interface AutomationSettings {
  notifications_enabled: boolean;
  default_channel: "whatsapp" | "email";
  daily_digest: boolean;
}

export interface CommunicationSettings {
  email_provider: string;
  whatsapp_provider: string;
  enabled_channels: string[];
}

export interface OrgSettings {
  id: string;
  name: string;
  slug: string;
  plan: string;
  subscription_plan?: string;
  logo_url?: string | null;
  industry?: string | null;
  company_size?: string | null;
  country?: string | null;
  status?: string;
  general: GeneralSettings;
  automation: AutomationSettings;
  communication?: CommunicationSettings;
  updated_at: string;
}

// Partial update: send only the sections that changed.
export interface UpdateSettingsPayload {
  name?: string;
  logo_url?: string | null;
  industry?: string | null;
  company_size?: string | null;
  country?: string | null;
  general?: GeneralSettings;
  automation?: AutomationSettings;
  communication?: CommunicationSettings;
}

// ---------------------------------------------------------------------------
// Metadata administration (modules / fields / validation rules)
// These mirror the existing metadata engines' admin DTOs and are used by the
// Module, Field and Validation tabs of the Settings Center.
// ---------------------------------------------------------------------------

import { FieldOption, FieldType } from "@/features/metadata/types";

// Full module record (superset of metadata's ModuleSummary) used by the module
// management table + edit form.
export interface ModuleDetail {
  id: string;
  api_name: string;
  singular_label: string;
  plural_label: string;
  description: string | null;
  icon: string | null;
  color: string | null;
  storage_strategy: "native" | "dynamic";
  native_table: string | null;
  is_system: boolean;
  is_enabled: boolean;
  is_visible_sidebar: boolean;
  sort_order: number;
  default_sort_field: string;
  default_sort_order: "asc" | "desc";
  created_at: string;
  updated_at: string;
}

export interface CreateModulePayload {
  api_name: string;
  singular_label: string;
  plural_label: string;
  description?: string | null;
  icon?: string | null;
  color?: string | null;
  is_visible_sidebar?: boolean;
  default_sort_field?: string;
  default_sort_order?: "asc" | "desc";
}

export interface UpdateModulePayload {
  singular_label?: string;
  plural_label?: string;
  description?: string | null;
  icon?: string | null;
  color?: string | null;
  is_visible_sidebar?: boolean;
  default_sort_field?: string;
  default_sort_order?: "asc" | "desc";
}

export interface CreateFieldPayload {
  api_name: string;
  label: string;
  field_type: FieldType;
  is_required?: boolean;
  is_unique?: boolean;
  is_read_only?: boolean;
  default_value?: string | null;
  placeholder?: string | null;
  description?: string | null;
  help_text?: string | null;
  min_length?: number | null;
  max_length?: number | null;
  regex?: string | null;
  validation_message?: string | null;
  options?: FieldOption[];
  lookup_module_id?: string | null;
  is_visible?: boolean;
  is_searchable?: boolean;
  is_filterable?: boolean;
  /** Detail layout section key (defaults to "general" on the server). */
  section_key?: string;
}

export interface UpdateFieldPayload {
  label?: string;
  is_required?: boolean;
  is_unique?: boolean;
  is_read_only?: boolean;
  default_value?: string | null;
  placeholder?: string | null;
  description?: string | null;
  help_text?: string | null;
  min_length?: number | null;
  max_length?: number | null;
  regex?: string | null;
  validation_message?: string | null;
  options?: FieldOption[];
  is_visible?: boolean;
  is_searchable?: boolean;
  is_filterable?: boolean;
}

export type RuleType =
  | "required"
  | "min_length"
  | "max_length"
  | "min"
  | "max"
  | "pattern"
  | "email"
  | "url"
  | "in"
  | "not_in"
  | "required_if";

export interface ValidationRule {
  id: string;
  module_id: string;
  field_id: string | null;
  rule_type: RuleType;
  params: Record<string, unknown>;
  error_message: string | null;
  is_active: boolean;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface CreateRulePayload {
  field_id?: string | null;
  rule_type: RuleType;
  params?: Record<string, unknown>;
  error_message?: string | null;
  is_active?: boolean;
  sort_order?: number;
}

export interface UpdateRulePayload {
  params?: Record<string, unknown>;
  error_message?: string | null;
  is_active?: boolean;
  sort_order?: number;
}
