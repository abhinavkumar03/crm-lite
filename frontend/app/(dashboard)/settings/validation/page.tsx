"use client";

import { useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { Plus, Pencil, Trash2 } from "lucide-react";

import Modal from "@/components/common/Modal";
import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";
import Toggle from "@/components/common/form/Toggle";

import {
  createRule,
  deleteRule,
  listFields,
  listModules,
  listRules,
  updateRule,
} from "@/features/settings/metadata";
import { ModuleDetail, RuleType, ValidationRule } from "@/features/settings/types";
import { apiErrorMessage } from "@/features/settings/errors";
import { ModuleField } from "@/features/metadata/types";

const RULE_TYPES: RuleType[] = [
  "required",
  "email",
  "url",
  "min_length",
  "max_length",
  "min",
  "max",
  "pattern",
  "in",
  "not_in",
  "required_if",
];

const VALUE_PARAM: RuleType[] = ["min_length", "max_length", "min", "max"];
const VALUES_PARAM: RuleType[] = ["in", "not_in"];

function ruleCategory(rt: RuleType): "none" | "value" | "pattern" | "values" | "required_if" {
  if (rt === "required_if") return "required_if";
  if (rt === "pattern") return "pattern";
  if (VALUE_PARAM.includes(rt)) return "value";
  if (VALUES_PARAM.includes(rt)) return "values";
  return "none";
}

function summarizeParams(rule: ValidationRule): string {
  const p = rule.params ?? {};
  switch (ruleCategory(rule.rule_type)) {
    case "value":
      return `value = ${p.value ?? "?"}`;
    case "pattern":
      return `/${p.pattern ?? ""}/`;
    case "values":
      return Array.isArray(p.values) ? (p.values as unknown[]).join(", ") : "";
    case "required_if":
      return `if ${p.field} = "${p.equals}" ⇒ ${p.target} required`;
    default:
      return "—";
  }
}

export default function ValidationSettingsPage() {
  const [modules, setModules] = useState<ModuleDetail[]>([]);
  const [moduleId, setModuleId] = useState("");
  const [fields, setFields] = useState<ModuleField[]>([]);
  const [rules, setRules] = useState<ValidationRule[]>([]);
  const [loading, setLoading] = useState(false);
  const [reloadKey, setReloadKey] = useState(0);

  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<ValidationRule | null>(null);
  const [saving, setSaving] = useState(false);

  // Form state
  const [ruleType, setRuleType] = useState<RuleType>("required");
  const [fieldId, setFieldId] = useState("");
  const [valueParam, setValueParam] = useState("");
  const [patternParam, setPatternParam] = useState("");
  const [valuesParam, setValuesParam] = useState("");
  const [ifField, setIfField] = useState("");
  const [ifEquals, setIfEquals] = useState("");
  const [ifTarget, setIfTarget] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const [isActive, setIsActive] = useState(true);

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
      setLoading(true);
      try {
        const [f, r] = await Promise.all([
          listFields(moduleId),
          listRules(moduleId),
        ]);
        if (!active) return;
        setFields(f);
        setRules(r);
      } catch {
        if (active) toast.error("Failed to load validation rules");
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [moduleId, reloadKey]);

  const fieldLabel = useMemo(() => {
    const map = new Map(fields.map((f) => [f.id, f.label]));
    return (id: string | null) => (id ? map.get(id) ?? id : "—");
  }, [fields]);

  function resetForm() {
    setRuleType("required");
    setFieldId(fields[0]?.id ?? "");
    setValueParam("");
    setPatternParam("");
    setValuesParam("");
    setIfField(fields[0]?.api_name ?? "");
    setIfEquals("");
    setIfTarget(fields[0]?.api_name ?? "");
    setErrorMessage("");
    setIsActive(true);
  }

  function openCreate() {
    setEditing(null);
    resetForm();
    setModalOpen(true);
  }

  function openEdit(rule: ValidationRule) {
    setEditing(rule);
    setRuleType(rule.rule_type);
    setFieldId(rule.field_id ?? "");
    const p = rule.params ?? {};
    setValueParam(p.value != null ? String(p.value) : "");
    setPatternParam(typeof p.pattern === "string" ? p.pattern : "");
    setValuesParam(Array.isArray(p.values) ? (p.values as unknown[]).join(", ") : "");
    setIfField(typeof p.field === "string" ? p.field : "");
    setIfEquals(p.equals != null ? String(p.equals) : "");
    setIfTarget(typeof p.target === "string" ? p.target : "");
    setErrorMessage(rule.error_message ?? "");
    setIsActive(rule.is_active);
    setModalOpen(true);
  }

  function buildParams(): Record<string, unknown> | null {
    switch (ruleCategory(ruleType)) {
      case "value": {
        const n = Number(valueParam);
        if (!valueParam.trim() || !Number.isFinite(n)) {
          toast.error("A numeric value is required");
          return null;
        }
        return { value: n };
      }
      case "pattern":
        if (!patternParam.trim()) {
          toast.error("A regex pattern is required");
          return null;
        }
        return { pattern: patternParam.trim() };
      case "values": {
        const values = valuesParam
          .split(",")
          .map((v) => v.trim())
          .filter(Boolean);
        if (values.length === 0) {
          toast.error("Provide at least one comma-separated value");
          return null;
        }
        return { values };
      }
      case "required_if":
        if (!ifField || !ifTarget || !ifEquals.trim()) {
          toast.error("Field, equals and target are required");
          return null;
        }
        return { field: ifField, equals: ifEquals.trim(), target: ifTarget };
      default:
        return {};
    }
  }

  async function handleSubmit() {
    if (!moduleId) return;
    const params = buildParams();
    if (params === null) return;

    const isFieldLevel = ruleType !== "required_if";
    if (isFieldLevel && !editing && !fieldId) {
      toast.error("Select the field this rule applies to");
      return;
    }

    try {
      setSaving(true);
      if (editing) {
        await updateRule(moduleId, editing.id, {
          params,
          error_message: errorMessage.trim() || null,
          is_active: isActive,
        });
        toast.success("Rule updated");
      } else {
        await createRule(moduleId, {
          rule_type: ruleType,
          field_id: isFieldLevel ? fieldId : null,
          params,
          error_message: errorMessage.trim() || null,
          is_active: isActive,
        });
        toast.success("Rule created");
      }
      setModalOpen(false);
      setReloadKey((k) => k + 1);
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to save rule"));
    } finally {
      setSaving(false);
    }
  }

  async function handleToggle(rule: ValidationRule) {
    try {
      const updated = await updateRule(moduleId, rule.id, {
        is_active: !rule.is_active,
      });
      setRules((prev) => prev.map((r) => (r.id === rule.id ? updated : r)));
    } catch {
      toast.error("Failed to update rule");
    }
  }

  async function handleDelete(rule: ValidationRule) {
    if (!confirm("Delete this validation rule?")) return;
    try {
      await deleteRule(moduleId, rule.id);
      toast.success("Rule deleted");
      setReloadKey((k) => k + 1);
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to delete rule"));
    }
  }

  const category = ruleCategory(ruleType);

  return (
    <div className="space-y-5" data-tutorial-surface="validation">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div className="min-w-[220px]">
          <h2 className="text-lg font-semibold text-slate-900">Validation</h2>
          <p className="text-sm text-slate-500">
            Database-driven rules enforced by the validation engine and mirrored
            to the client schema.
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
            New rule
          </button>
        </div>
      </div>

      <div className="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-sm">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-slate-200 bg-slate-50 text-left text-slate-600">
              <th className="px-4 py-3 font-semibold">Rule</th>
              <th className="px-4 py-3 font-semibold">Target</th>
              <th className="px-4 py-3 font-semibold">Params</th>
              <th className="px-4 py-3 font-semibold">Active</th>
              <th className="px-4 py-3 text-right font-semibold">Actions</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={5} className="px-4 py-10 text-center text-slate-400">
                  Loading rules...
                </td>
              </tr>
            ) : rules.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-4 py-10 text-center text-slate-400">
                  No validation rules yet.
                </td>
              </tr>
            ) : (
              rules.map((rule) => (
                <tr
                  key={rule.id}
                  className="border-b border-slate-100 last:border-0 hover:bg-slate-50/60"
                >
                  <td className="px-4 py-3">
                    <span className="inline-flex rounded-full bg-slate-100 px-2 py-0.5 text-xs font-semibold text-slate-700">
                      {rule.rule_type}
                    </span>
                    {rule.error_message && (
                      <div className="mt-1 max-w-xs truncate text-xs text-slate-400">
                        {rule.error_message}
                      </div>
                    )}
                  </td>
                  <td className="px-4 py-3 text-slate-700">
                    {rule.rule_type === "required_if" ? (
                      <span className="text-slate-500">module-level</span>
                    ) : (
                      fieldLabel(rule.field_id)
                    )}
                  </td>
                  <td className="px-4 py-3 text-slate-600">
                    {summarizeParams(rule)}
                  </td>
                  <td className="px-4 py-3">
                    <Toggle
                      checked={rule.is_active}
                      onChange={() => handleToggle(rule)}
                    />
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        type="button"
                        onClick={() => openEdit(rule)}
                        className="rounded-lg p-2 text-slate-500 transition hover:bg-slate-100 hover:text-slate-700"
                        aria-label="Edit"
                      >
                        <Pencil className="h-4 w-4" />
                      </button>
                      <button
                        type="button"
                        onClick={() => handleDelete(rule)}
                        className="rounded-lg p-2 text-red-500 transition hover:bg-red-50"
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
        title={editing ? "Edit rule" : "New validation rule"}
        onClose={() => setModalOpen(false)}
      >
        <div className="space-y-5">
          <div className="grid gap-4 sm:grid-cols-2">
            <FormSelect
              label="Rule type"
              value={ruleType}
              disabled={!!editing}
              helperText={editing ? "Immutable." : undefined}
              onChange={(e) => setRuleType(e.target.value as RuleType)}
            >
              {RULE_TYPES.map((rt) => (
                <option key={rt} value={rt}>
                  {rt}
                </option>
              ))}
            </FormSelect>

            {ruleType !== "required_if" && (
              <FormSelect
                label="Field"
                value={fieldId}
                requiredMark={!editing}
                disabled={!!editing}
                helperText={editing ? "Immutable." : undefined}
                onChange={(e) => setFieldId(e.target.value)}
              >
                <option value="">Select a field…</option>
                {fields.map((f) => (
                  <option key={f.id} value={f.id}>
                    {f.label}
                  </option>
                ))}
              </FormSelect>
            )}
          </div>

          {category === "value" && (
            <FormInput
              label="Value"
              type="number"
              value={valueParam}
              requiredMark
              onChange={(e) => setValueParam(e.target.value)}
            />
          )}

          {category === "pattern" && (
            <FormInput
              label="Regex pattern"
              placeholder="^[A-Z].*"
              value={patternParam}
              requiredMark
              onChange={(e) => setPatternParam(e.target.value)}
            />
          )}

          {category === "values" && (
            <FormInput
              label="Allowed values"
              placeholder="draft, active, closed"
              helperText="Comma-separated."
              value={valuesParam}
              requiredMark
              onChange={(e) => setValuesParam(e.target.value)}
            />
          )}

          {category === "required_if" && (
            <div className="grid gap-4 sm:grid-cols-3">
              <FormSelect
                label="When field"
                value={ifField}
                onChange={(e) => setIfField(e.target.value)}
              >
                {fields.map((f) => (
                  <option key={f.id} value={f.api_name}>
                    {f.label}
                  </option>
                ))}
              </FormSelect>
              <FormInput
                label="Equals"
                value={ifEquals}
                onChange={(e) => setIfEquals(e.target.value)}
              />
              <FormSelect
                label="Then require"
                value={ifTarget}
                onChange={(e) => setIfTarget(e.target.value)}
              >
                {fields.map((f) => (
                  <option key={f.id} value={f.api_name}>
                    {f.label}
                  </option>
                ))}
              </FormSelect>
            </div>
          )}

          <FormInput
            label="Custom error message"
            placeholder="Leave blank for the default message"
            value={errorMessage}
            onChange={(e) => setErrorMessage(e.target.value)}
          />

          <div className="rounded-2xl border border-slate-200 p-4">
            <Toggle
              label="Active"
              description="Inactive rules are stored but not enforced."
              checked={isActive}
              onChange={setIsActive}
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
              {saving ? "Saving..." : editing ? "Save changes" : "Create rule"}
            </button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
