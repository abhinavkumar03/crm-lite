import api from "@/services/api";

import {
  AccessContext,
  CreateRolePayload,
  FieldAccess,
  ModuleAccess,
  Permission,
  RoleDetail,
  RoleSummary,
  UpdateRolePayload,
} from "./types";

export async function getMyAccess(): Promise<AccessContext> {
  const res = await api.get("/me/access");
  return res.data.data;
}

export async function listPermissions(): Promise<Permission[]> {
  const res = await api.get("/permissions");
  return res.data.data;
}

export async function listRoles(): Promise<RoleSummary[]> {
  const res = await api.get("/roles");
  return res.data.data;
}

export async function getRole(roleId: string): Promise<RoleDetail> {
  const res = await api.get(`/roles/${roleId}`);
  return res.data.data;
}

export async function createRole(payload: CreateRolePayload): Promise<RoleDetail> {
  const res = await api.post("/roles", payload);
  return res.data.data;
}

export async function updateRole(
  roleId: string,
  payload: UpdateRolePayload
): Promise<RoleDetail> {
  const res = await api.put(`/roles/${roleId}`, payload);
  return res.data.data;
}

export async function deleteRole(roleId: string): Promise<void> {
  await api.delete(`/roles/${roleId}`);
}

export async function setRolePermissions(
  roleId: string,
  permissions: string[]
): Promise<RoleDetail> {
  const res = await api.put(`/roles/${roleId}/permissions`, { permissions });
  return res.data.data;
}

export async function setModuleAccess(
  roleId: string,
  access: ModuleAccess[]
): Promise<RoleDetail> {
  const res = await api.put(`/roles/${roleId}/module-access`, { access });
  return res.data.data;
}

export async function setFieldAccess(
  roleId: string,
  access: FieldAccess[]
): Promise<RoleDetail> {
  const res = await api.put(`/roles/${roleId}/field-access`, { access });
  return res.data.data;
}
