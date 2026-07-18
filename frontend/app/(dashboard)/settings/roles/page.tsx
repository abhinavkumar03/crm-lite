"use client";

import { useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { Plus, Save, Trash2 } from "lucide-react";

import Modal from "@/components/common/Modal";
import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";
import FormTextarea from "@/components/common/form/FormTextarea";
import Toggle from "@/components/common/form/Toggle";

import { listFields, listModules } from "@/features/settings/metadata";
import { ModuleDetail } from "@/features/settings/types";
import { apiErrorMessage } from "@/features/settings/errors";
import { ModuleField } from "@/features/metadata/types";

import {
  createRole,
  deleteRole,
  getRole,
  listPermissions,
  listRoles,
  setFieldAccess,
  setModuleAccess,
  setRolePermissions,
  updateRole,
} from "@/features/roles/api";
import {
  FieldAccess,
  FieldAccessLevel,
  ModuleAccess,
  Permission,
  RoleDetail,
  RoleSummary,
} from "@/features/roles/types";

type Tab = "permissions" | "modules" | "fields";

export default function RolesSettingsPage() {
  const [roles, setRoles] = useState<RoleSummary[]>([]);
  const [catalog, setCatalog] = useState<Permission[]>([]);
  const [modules, setModules] = useState<ModuleDetail[]>([]);
  const [selectedId, setSelectedId] = useState("");
  const [detail, setDetail] = useState<RoleDetail | null>(null);
  const [tab, setTab] = useState<Tab>("permissions");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  // Local editable state for the selected role
  const [permSet, setPermSet] = useState<Set<string>>(new Set());
  const [moduleACL, setModuleACL] = useState<Record<string, ModuleAccess>>({});
  const [fieldACL, setFieldACL] = useState<Record<string, FieldAccessLevel>>({});
  const [fieldModuleId, setFieldModuleId] = useState("");
  const [fields, setFields] = useState<ModuleField[]>([]);

  const [createOpen, setCreateOpen] = useState(false);
  const [newName, setNewName] = useState("");
  const [newSlug, setNewSlug] = useState("");
  const [newDesc, setNewDesc] = useState("");

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const [r, p, m] = await Promise.all([
          listRoles(),
          listPermissions(),
          listModules(),
        ]);
        if (!active) return;
        setRoles(r);
        setCatalog(p);
        setModules(m);
        if (r.length) setSelectedId((cur) => cur || r[0].id);
        if (m.length) setFieldModuleId((cur) => cur || m[0].id);
      } catch (err) {
        toast.error(apiErrorMessage(err, "Failed to load roles"));
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    if (!selectedId) return;
    let active = true;
    (async () => {
      try {
        const d = await getRole(selectedId);
        if (!active) return;
        setDetail(d);
        setPermSet(new Set(d.permissions));
        const macl: Record<string, ModuleAccess> = {};
        for (const a of d.module_access) macl[a.module_id] = a;
        setModuleACL(macl);
        const facl: Record<string, FieldAccessLevel> = {};
        for (const a of d.field_access) facl[a.field_id] = a.access;
        setFieldACL(facl);
      } catch (err) {
        toast.error(apiErrorMessage(err, "Failed to load role"));
      }
    })();
    return () => {
      active = false;
    };
  }, [selectedId]);

  useEffect(() => {
    if (!fieldModuleId) return;
    let active = true;
    (async () => {
      try {
        const data = await listFields(fieldModuleId);
        if (active) setFields(data);
      } catch {
        if (active) toast.error("Failed to load fields");
      }
    })();
    return () => {
      active = false;
    };
  }, [fieldModuleId]);

  const byCategory = useMemo(() => {
    const map = new Map<string, Permission[]>();
    for (const p of catalog) {
      const list = map.get(p.category) ?? [];
      list.push(p);
      map.set(p.category, list);
    }
    return Array.from(map.entries());
  }, [catalog]);

  function togglePerm(key: string) {
    setPermSet((prev) => {
      const next = new Set(prev);
      if (next.has(key)) next.delete(key);
      else next.add(key);
      return next;
    });
  }

  function ensureModuleAccess(moduleId: string): ModuleAccess {
    return (
      moduleACL[moduleId] ?? {
        module_id: moduleId,
        can_view: true,
        can_create: true,
        can_update: true,
        can_delete: false,
      }
    );
  }

  function patchModuleAccess(
    moduleId: string,
    patch: Partial<ModuleAccess>,
    enabled: boolean
  ) {
    setModuleACL((prev) => {
      const next = { ...prev };
      if (!enabled) {
        delete next[moduleId];
        return next;
      }
      next[moduleId] = { ...ensureModuleAccess(moduleId), ...patch, module_id: moduleId };
      return next;
    });
  }

  async function savePermissions() {
    if (!selectedId) return;
    try {
      setSaving(true);
      const d = await setRolePermissions(selectedId, Array.from(permSet));
      setDetail(d);
      toast.success("Permission matrix saved");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to save permissions"));
    } finally {
      setSaving(false);
    }
  }

  async function saveModuleAccess() {
    if (!selectedId) return;
    try {
      setSaving(true);
      const d = await setModuleAccess(selectedId, Object.values(moduleACL));
      setDetail(d);
      toast.success("Module access saved");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to save module access"));
    } finally {
      setSaving(false);
    }
  }

  async function saveFieldAccess() {
    if (!selectedId) return;
    try {
      setSaving(true);
      // Keep ACL rows for other modules; replace only rows for the current module's fields.
      const otherFieldIds = new Set(fields.map((f) => f.id));
      const kept: FieldAccess[] = (detail?.field_access ?? []).filter(
        (a) => !otherFieldIds.has(a.field_id)
      );
      const current: FieldAccess[] = Object.entries(fieldACL)
        .filter(([id, level]) => otherFieldIds.has(id) && level !== "write")
        .map(([field_id, access]) => ({ field_id, access }));
      // Also include explicit write? We only persist restrictions (hidden/read).
      // But if user set write after it was restricted, we drop it (absence = write).
      const d = await setFieldAccess(selectedId, [...kept, ...current]);
      setDetail(d);
      const facl: Record<string, FieldAccessLevel> = {};
      for (const a of d.field_access) facl[a.field_id] = a.access;
      setFieldACL(facl);
      toast.success("Field access saved");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to save field access"));
    } finally {
      setSaving(false);
    }
  }

  async function handleCreate() {
    if (!newName.trim() || !newSlug.trim()) {
      toast.error("Name and slug are required");
      return;
    }
    try {
      setSaving(true);
      const d = await createRole({
        name: newName.trim(),
        slug: newSlug.trim().toLowerCase(),
        description: newDesc.trim() || null,
      });
      setRoles((prev) => [...prev, d]);
      setSelectedId(d.id);
      setCreateOpen(false);
      setNewName("");
      setNewSlug("");
      setNewDesc("");
      toast.success("Role created");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to create role"));
    } finally {
      setSaving(false);
    }
  }

  async function handleRename(name: string) {
    if (!selectedId || !detail) return;
    try {
      const d = await updateRole(selectedId, { name });
      setDetail(d);
      setRoles((prev) => prev.map((r) => (r.id === d.id ? d : r)));
      toast.success("Role updated");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to update role"));
    }
  }

  async function handleDelete() {
    if (!detail || detail.is_system) return;
    if (!confirm(`Delete role "${detail.name}"?`)) return;
    try {
      await deleteRole(detail.id);
      const remaining = roles.filter((r) => r.id !== detail.id);
      setRoles(remaining);
      setSelectedId(remaining[0]?.id ?? "");
      setDetail(null);
      toast.success("Role deleted");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to delete role"));
    }
  }

  if (loading) {
    return (
      <div className="rounded-3xl border border-slate-200 bg-white p-8 text-slate-400 shadow-sm">
        Loading roles...
      </div>
    );
  }

  return (
    <div className="space-y-5">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">
            Roles & Permissions
          </h2>
          <p className="text-sm text-slate-500">
            Global permission matrix, plus per-module and per-field access
            control.
          </p>
        </div>
        <button
          type="button"
          onClick={() => setCreateOpen(true)}
          className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-4 py-2 text-sm font-semibold text-white transition hover:bg-emerald-600"
        >
          <Plus className="h-4 w-4" />
          New role
        </button>
      </div>

      <div className="grid gap-5 lg:grid-cols-[240px_1fr]">
        {/* Role list */}
        <div className="space-y-1 rounded-3xl border border-slate-200 bg-white p-3 shadow-sm">
          {roles.map((r) => (
            <button
              key={r.id}
              type="button"
              onClick={() => setSelectedId(r.id)}
              className={`flex w-full flex-col rounded-2xl px-3 py-2.5 text-left transition ${
                selectedId === r.id
                  ? "bg-emerald-50 text-emerald-800"
                  : "hover:bg-slate-50"
              }`}
            >
              <span className="text-sm font-semibold">{r.name}</span>
              <span className="text-xs text-slate-500">
                {r.slug}
                {` · L${r.hierarchy_level ?? 100}`}
                {r.is_system ? " · system" : ""}
                {` · ${r.member_count} member${r.member_count === 1 ? "" : "s"}`}
              </span>
            </button>
          ))}
        </div>

        {/* Editor */}
        {detail ? (
          <div className="space-y-5">
            <div className="flex flex-wrap items-center justify-between gap-3 rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
              <div className="min-w-0 flex-1">
                <FormInput
                  label="Role name"
                  value={detail.name}
                  onChange={(e) =>
                    setDetail({ ...detail, name: e.target.value })
                  }
                  onBlur={(e) => {
                    if (e.target.value.trim() !== roles.find((r) => r.id === detail.id)?.name) {
                      handleRename(e.target.value.trim());
                    }
                  }}
                />
              </div>
              {!detail.is_system && (
                <button
                  type="button"
                  onClick={handleDelete}
                  className="mt-6 inline-flex items-center gap-2 rounded-full border border-red-200 px-4 py-2 text-sm font-semibold text-red-600 transition hover:bg-red-50"
                >
                  <Trash2 className="h-4 w-4" />
                  Delete
                </button>
              )}
            </div>

            <div className="flex gap-2">
              {(
                [
                  ["permissions", "Permission matrix"],
                  ["modules", "Module access"],
                  ["fields", "Field access"],
                ] as const
              ).map(([key, label]) => (
                <button
                  key={key}
                  type="button"
                  onClick={() => setTab(key)}
                  className={`rounded-full px-4 py-1.5 text-sm font-semibold transition ${
                    tab === key
                      ? "bg-emerald-500 text-white"
                      : "bg-slate-100 text-slate-600 hover:bg-slate-200"
                  }`}
                >
                  {label}
                </button>
              ))}
            </div>

            {tab === "permissions" && (
              <section className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
                <p className="text-sm text-slate-500">
                  Tick the global capabilities this role may use. Import/export
                  and schema management are controlled here.
                </p>
                {byCategory.map(([category, perms]) => (
                  <div key={category}>
                    <h3 className="mb-2 text-xs font-semibold uppercase tracking-widest text-slate-400">
                      {category}
                    </h3>
                    <div className="grid gap-2 sm:grid-cols-2">
                      {perms.map((p) => (
                        <label
                          key={p.key}
                          className="flex cursor-pointer items-start gap-3 rounded-2xl border border-slate-200 px-3 py-2.5 hover:bg-slate-50"
                        >
                          <input
                            type="checkbox"
                            checked={permSet.has(p.key)}
                            onChange={() => togglePerm(p.key)}
                            className="mt-1 h-4 w-4 rounded border-slate-300 text-emerald-600 focus:ring-emerald-500"
                          />
                          <span>
                            <span className="block text-sm font-semibold text-slate-800">
                              {p.key}
                            </span>
                            {p.description && (
                              <span className="block text-xs text-slate-500">
                                {p.description}
                              </span>
                            )}
                          </span>
                        </label>
                      ))}
                    </div>
                  </div>
                ))}
                <div className="flex justify-end pt-2">
                  <button
                    type="button"
                    onClick={savePermissions}
                    disabled={saving}
                    className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:opacity-50"
                  >
                    <Save className="h-4 w-4" />
                    {saving ? "Saving..." : "Save matrix"}
                  </button>
                </div>
              </section>
            )}

            {tab === "modules" && (
              <section className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
                <p className="text-sm text-slate-500">
                  Optional per-module CRUD overrides. No row means unrestricted
                  (the global record.* permissions apply). Enable a restriction
                  to lock a module down for this role.
                </p>
                <div className="overflow-hidden rounded-2xl border border-slate-200">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50 text-left text-slate-600">
                        <th className="px-3 py-2 font-semibold">Module</th>
                        <th className="px-3 py-2 font-semibold">Restrict</th>
                        <th className="px-3 py-2 font-semibold">View</th>
                        <th className="px-3 py-2 font-semibold">Create</th>
                        <th className="px-3 py-2 font-semibold">Update</th>
                        <th className="px-3 py-2 font-semibold">Delete</th>
                      </tr>
                    </thead>
                    <tbody>
                      {modules.map((m) => {
                        const restricted = !!moduleACL[m.id];
                        const a = ensureModuleAccess(m.id);
                        return (
                          <tr
                            key={m.id}
                            className="border-b border-slate-100 last:border-0"
                          >
                            <td className="px-3 py-2 font-medium text-slate-800">
                              {m.plural_label}
                            </td>
                            <td className="px-3 py-2">
                              <Toggle
                                checked={restricted}
                                onChange={(v) =>
                                  patchModuleAccess(
                                    m.id,
                                    {
                                      can_view: true,
                                      can_create: true,
                                      can_update: true,
                                      can_delete: false,
                                    },
                                    v
                                  )
                                }
                              />
                            </td>
                            {(
                              [
                                "can_view",
                                "can_create",
                                "can_update",
                                "can_delete",
                              ] as const
                            ).map((key) => (
                              <td key={key} className="px-3 py-2">
                                <input
                                  type="checkbox"
                                  disabled={!restricted}
                                  checked={a[key]}
                                  onChange={(e) =>
                                    patchModuleAccess(
                                      m.id,
                                      { [key]: e.target.checked },
                                      true
                                    )
                                  }
                                  className="h-4 w-4 rounded border-slate-300 text-emerald-600 disabled:opacity-30"
                                />
                              </td>
                            ))}
                          </tr>
                        );
                      })}
                    </tbody>
                  </table>
                </div>
                <div className="flex justify-end">
                  <button
                    type="button"
                    onClick={saveModuleAccess}
                    disabled={saving}
                    className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:opacity-50"
                  >
                    <Save className="h-4 w-4" />
                    {saving ? "Saving..." : "Save module access"}
                  </button>
                </div>
              </section>
            )}

            {tab === "fields" && (
              <section className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
                <p className="text-sm text-slate-500">
                  Field-level ACL: hidden strips the field from forms and
                  responses, read shows it but blocks writes, write (default) is
                  full access.
                </p>
                <div className="w-56">
                  <FormSelect
                    label="Module"
                    value={fieldModuleId}
                    onChange={(e) => setFieldModuleId(e.target.value)}
                  >
                    {modules.map((m) => (
                      <option key={m.id} value={m.id}>
                        {m.plural_label}
                      </option>
                    ))}
                  </FormSelect>
                </div>
                <div className="overflow-hidden rounded-2xl border border-slate-200">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-slate-200 bg-slate-50 text-left text-slate-600">
                        <th className="px-3 py-2 font-semibold">Field</th>
                        <th className="px-3 py-2 font-semibold">API name</th>
                        <th className="px-3 py-2 font-semibold">Access</th>
                      </tr>
                    </thead>
                    <tbody>
                      {fields.length === 0 ? (
                        <tr>
                          <td
                            colSpan={3}
                            className="px-3 py-8 text-center text-slate-400"
                          >
                            No fields in this module.
                          </td>
                        </tr>
                      ) : (
                        fields.map((f) => (
                          <tr
                            key={f.id}
                            className="border-b border-slate-100 last:border-0"
                          >
                            <td className="px-3 py-2 font-medium text-slate-800">
                              {f.label}
                            </td>
                            <td className="px-3 py-2">
                              <code className="rounded bg-slate-100 px-1.5 py-0.5 text-xs">
                                {f.api_name}
                              </code>
                            </td>
                            <td className="px-3 py-2">
                              <select
                                value={fieldACL[f.id] ?? "write"}
                                onChange={(e) =>
                                  setFieldACL((prev) => ({
                                    ...prev,
                                    [f.id]: e.target.value as FieldAccessLevel,
                                  }))
                                }
                                className="rounded-lg border border-slate-200 bg-white px-2 py-1 text-sm focus:border-emerald-400 focus:outline-none"
                              >
                                <option value="write">Write</option>
                                <option value="read">Read</option>
                                <option value="hidden">Hidden</option>
                              </select>
                            </td>
                          </tr>
                        ))
                      )}
                    </tbody>
                  </table>
                </div>
                <div className="flex justify-end">
                  <button
                    type="button"
                    onClick={saveFieldAccess}
                    disabled={saving}
                    className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:opacity-50"
                  >
                    <Save className="h-4 w-4" />
                    {saving ? "Saving..." : "Save field access"}
                  </button>
                </div>
              </section>
            )}
          </div>
        ) : (
          <div className="rounded-3xl border border-slate-200 bg-white p-8 text-slate-400 shadow-sm">
            Select a role to edit.
          </div>
        )}
      </div>

      <Modal
        open={createOpen}
        title="New role"
        onClose={() => setCreateOpen(false)}
      >
        <div className="space-y-5">
          <FormInput
            label="Name"
            value={newName}
            requiredMark
            placeholder="Support Agent"
            onChange={(e) => {
              setNewName(e.target.value);
              if (!newSlug || newSlug === slugify(newName)) {
                setNewSlug(slugify(e.target.value));
              }
            }}
          />
          <FormInput
            label="Slug"
            value={newSlug}
            requiredMark
            helperText="Lowercase letters, numbers and underscores."
            onChange={(e) => setNewSlug(e.target.value)}
          />
          <FormTextarea
            label="Description"
            rows={2}
            value={newDesc}
            onChange={(e) => setNewDesc(e.target.value)}
          />
          <div className="flex justify-end gap-2">
            <button
              type="button"
              onClick={() => setCreateOpen(false)}
              className="rounded-full border border-slate-200 px-5 py-2.5 text-sm font-semibold text-slate-600 hover:bg-slate-50"
            >
              Cancel
            </button>
            <button
              type="button"
              onClick={handleCreate}
              disabled={saving}
              className="rounded-full bg-emerald-500 px-5 py-2.5 text-sm font-semibold text-white hover:bg-emerald-600 disabled:opacity-50"
            >
              {saving ? "Creating..." : "Create role"}
            </button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

function slugify(value: string): string {
  return value
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "_")
    .replace(/^_|_$/g, "")
    .slice(0, 100);
}
