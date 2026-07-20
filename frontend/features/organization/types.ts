export interface OrgMembership {
  id: string;
  name: string;
  slug: string;
  logo_url?: string | null;
  description?: string | null;
  role_slug: string;
  is_active: boolean;
}

export interface OrgMember {
  user_id: string;
  name: string;
  email: string;
  role_id?: string | null;
  role_slug: string;
  manager_user_id?: string | null;
  department_id?: string | null;
  team_id?: string | null;
  branch_id?: string | null;
  designation?: string | null;
  hierarchy_level: number;
  status: string;
}

export interface InviteResponse {
  id: string;
  email: string;
  token: string;
  status: string;
  expires_at: string;
  simulated_email_body: string;
}

export interface CreateInvitePayload {
  email: string;
  role_id: string;
  manager_user_id?: string | null;
  department_id?: string | null;
  team_id?: string | null;
}

export interface StructureItem {
  id: string;
  name: string;
  description?: string | null;
  location?: string | null;
  department_id?: string | null;
}

export interface CreateOrganizationPayload {
  name: string;
  slug?: string;
  description?: string;
  industry?: string;
  company_size?: string;
  country?: string;
  logo_url?: string;
  general?: {
    timezone?: string;
    currency?: string;
    locale?: string;
    date_format?: string;
  };
}

export interface CreateOrganizationResult {
  id: string;
}
