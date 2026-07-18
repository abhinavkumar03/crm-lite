import api from "@/services/api";

import {
  CreateInvitePayload,
  InviteResponse,
  OrgMember,
  OrgMembership,
  StructureItem,
} from "./types";

export async function listMyOrganizations(): Promise<OrgMembership[]> {
  const res = await api.get("/me/organizations");
  return res.data.data;
}

export async function switchOrganization(organizationId: string): Promise<void> {
  await api.post("/me/organizations/switch", {
    organization_id: organizationId,
  });
}

export async function listMembers(): Promise<OrgMember[]> {
  const res = await api.get("/organizations/members");
  return res.data.data;
}

export async function inviteMember(
  payload: CreateInvitePayload
): Promise<InviteResponse> {
  const res = await api.post("/organizations/invitations", payload);
  return res.data.data;
}

export async function listDepartments(): Promise<StructureItem[]> {
  const res = await api.get("/departments");
  return res.data.data;
}

export async function createDepartment(payload: {
  name: string;
  description?: string | null;
}): Promise<StructureItem> {
  const res = await api.post("/departments", payload);
  return res.data.data;
}

export async function listTeams(): Promise<StructureItem[]> {
  const res = await api.get("/teams");
  return res.data.data;
}

export async function createTeam(payload: {
  name: string;
  description?: string | null;
  department_id?: string | null;
}): Promise<StructureItem> {
  const res = await api.post("/teams", payload);
  return res.data.data;
}
