"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Plus, Pencil, Trash2, RefreshCw } from "lucide-react";

import Modal from "@/components/common/Modal";
import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";
import FormTextarea from "@/components/common/form/FormTextarea";
import Toggle from "@/components/common/form/Toggle";

import {
  createModule,
  deleteModule,
  listModules,
  setModuleStatus,
  updateModule,
} from "@/features/settings/metadata";
import { ModuleDetail } from "@/features/settings/types";
import { apiErrorMessage } from "@/features/settings/errors";

type FormState = {
  api_name: string;
  singular_label: string;
  plural_label: string;
  description: string;
  icon: string;
  color: string;
  is_visible_sidebar: boolean;
  default_sort_field: string;
  default_sort_order: "asc" | "desc";
};

const EMPTY: FormState = {
  api_name: "",
  singular_label: "",
  plural_label: "",
  description: "",
  icon: "",
  color: "",
  is_visible_sidebar: true,
  default_sort_field: "created_at",
  default_sort_order: "desc",
};

export default function ModulesSettingsPage() {
  const [modules, setModules] = useState<ModuleDetail[]>([]);
  const [loading, setLoading] = useState(true);
  const [reloadKey, setReloadKey] = useState(0);

  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<ModuleDetail | null>(null);
  const [form, setForm] = useState<FormState>(EMPTY);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const data = await listModules();
        if (active) setModules(data);
      } catch {
        toast.error("Failed to load modules");
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [reloadKey]);

  function openCreate() {
    setEditing(null);
    setForm(EMPTY);
    setModalOpen(true);
  }

  function openEdit(m: ModuleDetail) {
    setEditing(m);
    setForm({
      api_name: m.api_name,
      singular_label: m.singular_label,
      plural_label: m.plural_label,
      description: m.description ?? "",
      icon: m.icon ?? "",
      color: m.color ?? "",
      is_visible_sidebar: m.is_visible_sidebar,
      default_sort_field: m.default_sort_field,
      default_sort_order: m.default_sort_order,
    });
    setModalOpen(true);
  }

  function patch(p: Partial<FormState>) {
    setForm((prev) => ({ ...prev, ...p }));
  }

  async function handleSubmit() {
    if (!form.singular_label.trim() || !form.plural_label.trim()) {
      toast.error("Singular and plural labels are required");
      return;
    }
    if (!editing && !form.api_name.trim()) {
      toast.error("API name is required");
      return;
    }

    try {
      setSaving(true);
      if (editing) {
        await updateModule(editing.id, {
          singular_label: form.singular_label.trim(),
          plural_label: form.plural_label.trim(),
          description: form.description.trim() || null,
          icon: form.icon.trim() || null,
          color: form.color.trim() || null,
          is_visible_sidebar: form.is_visible_sidebar,
          default_sort_field: form.default_sort_field.trim() || "created_at",
          default_sort_order: form.default_sort_order,
        });
        toast.success("Module updated");
      } else {
        await createModule({
          api_name: form.api_name.trim(),
          singular_label: form.singular_label.trim(),
          plural_label: form.plural_label.trim(),
          description: form.description.trim() || null,
          icon: form.icon.trim() || null,
          color: form.color.trim() || null,
          is_visible_sidebar: form.is_visible_sidebar,
          default_sort_field: form.default_sort_field.trim() || "created_at",
          default_sort_order: form.default_sort_order,
        });
        toast.success("Module created");
      }
      setModalOpen(false);
      setReloadKey((k) => k + 1);
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to save module"));
    } finally {
      setSaving(false);
    }
  }

  async function handleToggle(m: ModuleDetail) {
    try {
      await setModuleStatus(m.id, !m.is_enabled);
      setModules((prev) =>
        prev.map((x) =>
          x.id === m.id ? { ...x, is_enabled: !m.is_enabled } : x
        )
      );
    } catch {
      toast.error("Failed to update status");
    }
  }

  async function handleDelete(m: ModuleDetail) {
    if (!confirm(`Delete "${m.plural_label}"? This cannot be undone.`)) return;
    try {
      await deleteModule(m.id);
      toast.success("Module deleted");
      setReloadKey((k) => k + 1);
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to delete module"));
    }
  }

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between gap-3">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">Modules</h2>
          <p className="text-sm text-slate-500">
            Object types in your CRM. Dynamic modules store data with no schema
            change; native modules are backed by first-class tables.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => setReloadKey((k) => k + 1)}
            className="rounded-xl border border-slate-200 p-2 text-slate-600 hover:bg-slate-100"
            aria-label="Refresh"
          >
            <RefreshCw className="h-4 w-4" />
          </button>
          <button
            type="button"
            onClick={openCreate}
            className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-4 py-2 text-sm font-semibold text-white transition hover:bg-emerald-600"
          >
            <Plus className="h-4 w-4" />
            New module
          </button>
        </div>
      </div>

      <div className="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-sm">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-slate-200 bg-slate-50 text-left text-slate-600">
              <th className="px-4 py-3 font-semibold">Module</th>
              <th className="px-4 py-3 font-semibold">API name</th>
              <th className="px-4 py-3 font-semibold">Type</th>
              <th className="px-4 py-3 font-semibold">Sidebar</th>
              <th className="px-4 py-3 font-semibold">Enabled</th>
              <th className="px-4 py-3 text-right font-semibold">Actions</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={6} className="px-4 py-10 text-center text-slate-400">
                  Loading modules...
                </td>
              </tr>
            ) : modules.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-4 py-10 text-center text-slate-400">
                  No modules yet.
                </td>
              </tr>
            ) : (
              modules.map((m) => (
                <tr
                  key={m.id}
                  className="border-b border-slate-100 last:border-0 hover:bg-slate-50/60"
                >
                  <td className="px-4 py-3">
                    <div className="font-semibold text-slate-800">
                      {m.plural_label}
                    </div>
                    {m.description && (
                      <div className="max-w-xs truncate text-xs text-slate-500">
                        {m.description}
                      </div>
                    )}
                  </td>
                  <td className="px-4 py-3">
                    <code className="rounded bg-slate-100 px-1.5 py-0.5 text-xs text-slate-700">
                      {m.api_name}
                    </code>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex flex-wrap gap-1">
                      <span
                        className={`inline-flex rounded-full px-2 py-0.5 text-xs font-semibold capitalize ${
                          m.storage_strategy === "native"
                            ? "bg-indigo-100 text-indigo-700"
                            : "bg-emerald-100 text-emerald-700"
                        }`}
                      >
                        {m.storage_strategy}
                      </span>
                      {m.is_system && (
                        <span className="inline-flex rounded-full bg-slate-100 px-2 py-0.5 text-xs font-semibold text-slate-600">
                          system
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <Toggle
                      checked={m.is_visible_sidebar}
                      onChange={async () => {
                        try {
                          const updated = await updateModule(m.id, {
                            is_visible_sidebar: !m.is_visible_sidebar,
                          });
                          setModules((prev) =>
                            prev.map((x) => (x.id === m.id ? updated : x))
                          );
                        } catch {
                          toast.error("Failed to update visibility");
                        }
                      }}
                    />
                  </td>
                  <td className="px-4 py-3">
                    <Toggle
                      checked={m.is_enabled}
                      onChange={() => handleToggle(m)}
                    />
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        type="button"
                        onClick={() => openEdit(m)}
                        className="rounded-lg p-2 text-slate-500 transition hover:bg-slate-100 hover:text-slate-700"
                        aria-label="Edit"
                      >
                        <Pencil className="h-4 w-4" />
                      </button>
                      <button
                        type="button"
                        onClick={() => handleDelete(m)}
                        disabled={m.is_system}
                        className="rounded-lg p-2 text-red-500 transition hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-30"
                        aria-label="Delete"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <Modal
        open={modalOpen}
        title={editing ? `Edit ${editing.singular_label}` : "New module"}
        onClose={() => setModalOpen(false)}
      >
        <div className="space-y-5">
          <div className="grid gap-4 sm:grid-cols-2">
            <FormInput
              label="API name"
              placeholder="invoices"
              value={form.api_name}
              requiredMark={!editing}
              disabled={!!editing}
              helperText={
                editing
                  ? "Immutable — changing it would orphan stored data."
                  : "Lowercase, unique per organization."
              }
              onChange={(e) => patch({ api_name: e.target.value })}
            />
            <div className="grid grid-cols-2 gap-4">
              <FormInput
                label="Icon"
                placeholder="FileText"
                value={form.icon}
                onChange={(e) => patch({ icon: e.target.value })}
              />
              <FormInput
                label="Color"
                placeholder="#10b981"
                value={form.color}
                onChange={(e) => patch({ color: e.target.value })}
              />
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <FormInput
              label="Singular label"
              placeholder="Invoice"
              value={form.singular_label}
              requiredMark
              onChange={(e) => patch({ singular_label: e.target.value })}
            />
            <FormInput
              label="Plural label"
              placeholder="Invoices"
              value={form.plural_label}
              requiredMark
              onChange={(e) => patch({ plural_label: e.target.value })}
            />
          </div>

          <FormTextarea
            label="Description"
            rows={2}
            value={form.description}
            onChange={(e) => patch({ description: e.target.value })}
          />

          <div className="grid gap-4 sm:grid-cols-2">
            <FormInput
              label="Default sort field"
              value={form.default_sort_field}
              onChange={(e) => patch({ default_sort_field: e.target.value })}
            />
            <FormSelect
              label="Default sort order"
              value={form.default_sort_order}
              onChange={(e) =>
                patch({ default_sort_order: e.target.value as "asc" | "desc" })
              }
            >
              <option value="desc">Descending</option>
              <option value="asc">Ascending</option>
            </FormSelect>
          </div>

          <div className="rounded-2xl border border-slate-200 p-4">
            <Toggle
              label="Show in sidebar"
              description="Display this module in the main navigation."
              checked={form.is_visible_sidebar}
              onChange={(v) => patch({ is_visible_sidebar: v })}
            />
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <button
              type="button"
              onClick={() => setModalOpen(false)}
              className="rounded-full border border-slate-200 px-5 py-2.5 text-sm font-semibold text-slate-600 transition hover:bg-slate-50"
            >
              Cancel
            </button>
            <button
              type="button"
              onClick={handleSubmit}
              disabled={saving}
              className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:opacity-50"
            >
              {saving ? "Saving..." : editing ? "Save changes" : "Create module"}
            </button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
