// Types mirroring the backend roles & permissions engine (Phase 16).

export interface Permission {
  id: string;
  key: string;
  category: string;
  description: string | null;
}

export interface ModuleAccess {
  module_id: string;
  can_view: boolean;
  can_create: boolean;
  can_update: boolean;
  can_delete: boolean;
}

export type FieldAccessLevel = "hidden" | "read" | "write";

export interface FieldAccess {
  field_id: string;
  access: FieldAccessLevel;
}

export interface RoleSummary {
  id: string;
  name: string;
  slug: string;
  description: string | null;
  is_system: boolean;
  member_count: number;
  created_at: string;
  updated_at: string;
}

export interface RoleDetail extends RoleSummary {
  permissions: string[];
  module_access: ModuleAccess[];
  field_access: FieldAccess[];
}

export interface CreateRolePayload {
  name: string;
  slug: string;
  description?: string | null;
}

export interface UpdateRolePayload {
  name?: string;
  description?: string | null;
}

export interface AccessContext {
  role_id: string;
  role_slug: string;
  permissions: string[];
  module_access: ModuleAccess[];
  field_access: FieldAccess[];
}
