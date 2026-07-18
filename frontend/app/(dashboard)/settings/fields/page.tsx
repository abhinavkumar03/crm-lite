"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Plus, Pencil, Trash2, X } from "lucide-react";

import Modal from "@/components/common/Modal";
import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";
import FormTextarea from "@/components/common/form/FormTextarea";
import Toggle from "@/components/common/form/Toggle";

import {
  createField,
  deleteField,
  listFields,
  listModules,
  updateField,
} from "@/features/settings/metadata";
import { ModuleDetail } from "@/features/settings/types";
import { apiErrorMessage } from "@/features/settings/errors";
import { FieldOption, FieldType, ModuleField } from "@/features/metadata/types";

const FIELD_TYPES: FieldType[] = [
  "text",
  "textarea",
  "email",
  "phone",
  "number",
  "currency",
  "date",
  "datetime",
  "boolean",
  "dropdown",
  "multiselect",
  "radio",
  "checkbox",
  "url",
  "user",
  "lookup",
  "json",
  "richtext",
];

const OPTION_TYPES: FieldType[] = ["dropdown", "multiselect", "radio"];

type FormState = {
  api_name: string;
  label: string;
  field_type: FieldType;
  is_required: boolean;
  is_unique: boolean;
  is_read_only: boolean;
  is_visible: boolean;
  is_searchable: boolean;
  is_filterable: boolean;
  placeholder: string;
  description: string;
  help_text: string;
  default_value: string;
  validation_message: string;
  regex: string;
  min_length: string;
  max_length: string;
  options: FieldOption[];
  lookup_module_id: string;
};

const EMPTY: FormState = {
  api_name: "",
  label: "",
  field_type: "text",
  is_required: false,
  is_unique: false,
  is_read_only: false,
  is_visible: true,
  is_searchable: false,
  is_filterable: false,
  placeholder: "",
  description: "",
  help_text: "",
  default_value: "",
  validation_message: "",
  regex: "",
  min_length: "",
  max_length: "",
  options: [],
  lookup_module_id: "",
};

function intOrNull(v: string): number | null {
  const t = v.trim();
  if (!t) return null;
  const n = Number(t);
  return Number.isFinite(n) ? n : null;
}

export default function FieldsSettingsPage() {
  const [modules, setModules] = useState<ModuleDetail[]>([]);
  const [moduleId, setModuleId] = useState("");
  const [fields, setFields] = useState<ModuleField[]>([]);
  const [loadingFields, setLoadingFields] = useState(false);

  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<ModuleField | null>(null);
  const [form, setForm] = useState<FormState>(EMPTY);
  const [saving, setSaving] = useState(false);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const data = await listModules();
        if (!active) return;
        setModules(data);
        if (data.length) setModuleId((cur) => cur || data[0].id);
      } catch {
        toast.error("Failed to load modules");
      }
    })();
    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    if (!moduleId) return;
    let active = true;
    (async () => {
      setLoadingFields(true);
      try {
        const data = await listFields(moduleId);
        if (active) setFields(data);
      } catch {
        if (active) toast.error("Failed to load fields");
      } finally {
        if (active) setLoadingFields(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [moduleId, reloadKey]);

  const selectedModule = modules.find((m) => m.id === moduleId);

  function openCreate() {
    setEditing(null);
    setForm(EMPTY);
    setModalOpen(true);
  }

  function openEdit(f: ModuleField) {
    setEditing(f);
    setForm({
      api_name: f.api_name,
      label: f.label,
      field_type: f.field_type,
      is_required: f.is_required,
      is_unique: f.is_unique,
      is_read_only: f.is_read_only,
      is_visible: f.is_visible,
      is_searchable: f.is_searchable,
      is_filterable: f.is_filterable,
      placeholder: f.placeholder ?? "",
      description: f.description ?? "",
      help_text: f.help_text ?? "",
      default_value: f.default_value ?? "",
      validation_message: f.validation_message ?? "",
      regex: f.regex ?? "",
      min_length: f.min_length?.toString() ?? "",
      max_length: f.max_length?.toString() ?? "",
      options: f.options ?? [],
      lookup_module_id: f.lookup_module_id ?? "",
    });
    setModalOpen(true);
  }

  function patch(p: Partial<FormState>) {
    setForm((prev) => ({ ...prev, ...p }));
  }

  function addOption() {
    patch({ options: [...form.options, { label: "", value: "" }] });
  }

  function updateOption(idx: number, p: Partial<FieldOption>) {
    patch({
      options: form.options.map((o, i) => (i === idx ? { ...o, ...p } : o)),
    });
  }

  function removeOption(idx: number) {
    patch({ options: form.options.filter((_, i) => i !== idx) });
  }

  async function handleSubmit() {
    if (!moduleId) return;
    if (!form.label.trim()) {
      toast.error("Label is required");
      return;
    }
    if (!editing && !form.api_name.trim()) {
      toast.error("API name is required");
      return;
    }

    const needsOptions = OPTION_TYPES.includes(form.field_type);
    const cleanOptions = form.options
      .map((o) => ({ label: o.label.trim(), value: o.value.trim() }))
      .filter((o) => o.label && o.value);
    if (needsOptions && cleanOptions.length === 0) {
      toast.error("Add at least one option for this field type");
      return;
    }

    try {
      setSaving(true);
      if (editing) {
        await updateField(moduleId, editing.id, {
          label: form.label.trim(),
          is_required: form.is_required,
          is_unique: form.is_unique,
          is_read_only: form.is_read_only,
          is_visible: form.is_visible,
          is_searchable: form.is_searchable,
          is_filterable: form.is_filterable,
          placeholder: form.placeholder.trim() || null,
          description: form.description.trim() || null,
          help_text: form.help_text.trim() || null,
          default_value: form.default_value.trim() || null,
          validation_message: form.validation_message.trim() || null,
          regex: form.regex.trim() || null,
          min_length: intOrNull(form.min_length),
          max_length: intOrNull(form.max_length),
          options: needsOptions ? cleanOptions : [],
        });
        toast.success("Field updated");
      } else {
        await createField(moduleId, {
          api_name: form.api_name.trim(),
          label: form.label.trim(),
          field_type: form.field_type,
          is_required: form.is_required,
          is_unique: form.is_unique,
          is_read_only: form.is_read_only,
          is_visible: form.is_visible,
          is_searchable: form.is_searchable,
          is_filterable: form.is_filterable,
          placeholder: form.placeholder.trim() || null,
          description: form.description.trim() || null,
          help_text: form.help_text.trim() || null,
          default_value: form.default_value.trim() || null,
          validation_message: form.validation_message.trim() || null,
          regex: form.regex.trim() || null,
          min_length: intOrNull(form.min_length),
          max_length: intOrNull(form.max_length),
          options: needsOptions ? cleanOptions : [],
          lookup_module_id:
            form.field_type === "lookup" && form.lookup_module_id
              ? form.lookup_module_id
              : null,
        });
        toast.success("Field created");
      }
      setModalOpen(false);
      setReloadKey((k) => k + 1);
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to save field"));
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(f: ModuleField) {
    if (!confirm(`Delete field "${f.label}"? This cannot be undone.`)) return;
    try {
      await deleteField(moduleId, f.id);
      toast.success("Field deleted");
      setReloadKey((k) => k + 1);
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to delete field"));
    }
  }

  const showOptions = OPTION_TYPES.includes(form.field_type);
  const showLookup = form.field_type === "lookup";

  return (
    <div className="space-y-5">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div className="min-w-[220px]">
          <h2 className="text-lg font-semibold text-slate-900">Fields</h2>
          <p className="text-sm text-slate-500">
            Define the fields that make up each module.
          </p>
        </div>
        <div className="flex items-end gap-2">
          <div className="w-56">
            <FormSelect
              label="Module"
              value={moduleId}
              onChange={(e) => setModuleId(e.target.value)}
            >
              {modules.map((m) => (
                <option key={m.id} value={m.id}>
                  {m.plural_label}
                </option>
              ))}
            </FormSelect>
          </div>
          <button
            type="button"
            onClick={openCreate}
            disabled={!moduleId}
            className="inline-flex h-[46px] items-center gap-2 rounded-full bg-emerald-500 px-4 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:opacity-50"
          >
            <Plus className="h-4 w-4" />
            New field
          </button>
        </div>
      </div>

      {selectedModule?.storage_strategy === "native" && (
        <div className="rounded-2xl border border-indigo-200 bg-indigo-50 px-4 py-3 text-sm text-indigo-700">
          This is a native module. Its fields map to real table columns — edit
          with care.
        </div>
      )}

      <div className="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-sm">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-slate-200 bg-slate-50 text-left text-slate-600">
              <th className="px-4 py-3 font-semibold">Field</th>
              <th className="px-4 py-3 font-semibold">API name</th>
              <th className="px-4 py-3 font-semibold">Type</th>
              <th className="px-4 py-3 font-semibold">Flags</th>
              <th className="px-4 py-3 text-right font-semibold">Actions</th>
            </tr>
          </thead>
          <tbody>
            {loadingFields ? (
              <tr>
                <td colSpan={5} className="px-4 py-10 text-center text-slate-400">
                  Loading fields...
                </td>
              </tr>
            ) : fields.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-4 py-10 text-center text-slate-400">
                  No fields yet.
                </td>
              </tr>
            ) : (
              fields.map((f) => (
                <tr
                  key={f.id}
                  className="border-b border-slate-100 last:border-0 hover:bg-slate-50/60"
                >
                  <td className="px-4 py-3 font-semibold text-slate-800">
                    {f.label}
                  </td>
                  <td className="px-4 py-3">
                    <code className="rounded bg-slate-100 px-1.5 py-0.5 text-xs text-slate-700">
                      {f.api_name}
                    </code>
                  </td>
                  <td className="px-4 py-3">
                    <span className="inline-flex rounded-full bg-slate-100 px-2 py-0.5 text-xs font-semibold text-slate-600">
                      {f.field_type}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex flex-wrap gap-1">
                      {f.is_required && (
                        <span className="rounded-full bg-amber-100 px-2 py-0.5 text-xs font-semibold text-amber-700">
                          required
                        </span>
                      )}
                      {f.is_unique && (
                        <span className="rounded-full bg-sky-100 px-2 py-0.5 text-xs font-semibold text-sky-700">
                          unique
                        </span>
                      )}
                      {f.is_system && (
                        <span className="rounded-full bg-slate-100 px-2 py-0.5 text-xs font-semibold text-slate-600">
                          system
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        type="button"
                        onClick={() => openEdit(f)}
                        className="rounded-lg p-2 text-slate-500 transition hover:bg-slate-100 hover:text-slate-700"
                        aria-label="Edit"
                      >
                        <Pencil className="h-4 w-4" />
                      </button>
                      <button
                        type="button"
                        onClick={() => handleDelete(f)}
                        disabled={f.is_system}
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
        title={editing ? `Edit ${editing.label}` : "New field"}
        onClose={() => setModalOpen(false)}
      >
        <div className="space-y-5">
          <div className="grid gap-4 sm:grid-cols-2">
            <FormInput
              label="API name"
              placeholder="amount_due"
              value={form.api_name}
              requiredMark={!editing}
              disabled={!!editing}
              helperText={editing ? "Immutable." : "Unique within the module."}
              onChange={(e) => patch({ api_name: e.target.value })}
            />
            <FormInput
              label="Label"
              placeholder="Amount due"
              value={form.label}
              requiredMark
              onChange={(e) => patch({ label: e.target.value })}
            />
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <FormSelect
              label="Field type"
              value={form.field_type}
              disabled={!!editing}
              helperText={editing ? "Immutable." : undefined}
              onChange={(e) =>
                patch({ field_type: e.target.value as FieldType })
              }
            >
              {FIELD_TYPES.map((t) => (
                <option key={t} value={t}>
                  {t}
                </option>
              ))}
            </FormSelect>

            {showLookup && (
              <FormSelect
                label="Lookup module"
                value={form.lookup_module_id}
                onChange={(e) => patch({ lookup_module_id: e.target.value })}
              >
                <option value="">Select a module…</option>
                {modules.map((m) => (
                  <option key={m.id} value={m.id}>
                    {m.plural_label}
                  </option>
                ))}
              </FormSelect>
            )}
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <FormInput
              label="Placeholder"
              value={form.placeholder}
              onChange={(e) => patch({ placeholder: e.target.value })}
            />
            <FormInput
              label="Default value"
              value={form.default_value}
              onChange={(e) => patch({ default_value: e.target.value })}
            />
          </div>

          <FormTextarea
            label="Help text"
            rows={2}
            value={form.help_text}
            onChange={(e) => patch({ help_text: e.target.value })}
          />

          {showOptions && (
            <div className="space-y-2 rounded-2xl border border-slate-200 p-4">
              <div className="flex items-center justify-between">
                <p className="text-sm font-semibold text-slate-700">Options</p>
                <button
                  type="button"
                  onClick={addOption}
                  className="inline-flex items-center gap-1 rounded-lg border border-slate-200 px-2 py-1 text-xs font-semibold text-slate-600 hover:bg-slate-50"
                >
                  <Plus className="h-3.5 w-3.5" />
                  Add
                </button>
              </div>
              {form.options.length === 0 && (
                <p className="text-xs text-slate-400">No options yet.</p>
              )}
              {form.options.map((o, i) => (
                <div key={i} className="flex items-center gap-2">
                  <input
                    placeholder="Label"
                    value={o.label}
                    onChange={(e) => updateOption(i, { label: e.target.value })}
                    className="w-full rounded-xl border border-slate-300 px-3 py-2 text-sm focus:border-emerald-500 focus:outline-none"
                  />
                  <input
                    placeholder="Value"
                    value={o.value}
                    onChange={(e) => updateOption(i, { value: e.target.value })}
                    className="w-full rounded-xl border border-slate-300 px-3 py-2 text-sm focus:border-emerald-500 focus:outline-none"
                  />
                  <button
                    type="button"
                    onClick={() => removeOption(i)}
                    className="rounded-lg p-2 text-red-500 hover:bg-red-50"
                    aria-label="Remove option"
                  >
                    <X className="h-4 w-4" />
                  </button>
                </div>
              ))}
            </div>
          )}

          <div className="grid gap-4 sm:grid-cols-2">
            <FormInput
              label="Min length"
              type="number"
              value={form.min_length}
              onChange={(e) => patch({ min_length: e.target.value })}
            />
            <FormInput
              label="Max length"
              type="number"
              value={form.max_length}
              onChange={(e) => patch({ max_length: e.target.value })}
            />
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <FormInput
              label="Regex pattern"
              value={form.regex}
              onChange={(e) => patch({ regex: e.target.value })}
            />
            <FormInput
              label="Validation message"
              value={form.validation_message}
              onChange={(e) => patch({ validation_message: e.target.value })}
            />
          </div>

          <div className="grid gap-3 rounded-2xl border border-slate-200 p-4 sm:grid-cols-2">
            <Toggle
              label="Required"
              checked={form.is_required}
              onChange={(v) => patch({ is_required: v })}
            />
            <Toggle
              label="Unique"
              checked={form.is_unique}
              onChange={(v) => patch({ is_unique: v })}
            />
            <Toggle
              label="Read only"
              checked={form.is_read_only}
              onChange={(v) => patch({ is_read_only: v })}
            />
            <Toggle
              label="Visible"
              checked={form.is_visible}
              onChange={(v) => patch({ is_visible: v })}
            />
            <Toggle
              label="Searchable"
              checked={form.is_searchable}
              onChange={(v) => patch({ is_searchable: v })}
            />
            <Toggle
              label="Filterable"
              checked={form.is_filterable}
              onChange={(v) => patch({ is_filterable: v })}
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
              {saving ? "Saving..." : editing ? "Save changes" : "Create field"}
            </button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
