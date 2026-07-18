"use client";

import { useEffect, useMemo, useState } from "react";
import axios from "axios";
import { toast } from "sonner";
import { Download, RefreshCw, Save, Trash2, FileDown } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import FormSelect from "@/components/common/form/FormSelect";
import FormInput from "@/components/common/form/FormInput";

import { getModules, getModuleFields } from "@/features/metadata/api";
import { ModuleField, ModuleSummary } from "@/features/metadata/types";
import {
  createExport,
  createTemplate,
  deleteTemplate,
  downloadExport,
  exportNow,
  listExports,
  listTemplates,
} from "@/features/export/api";
import {
  ExportFormat,
  ExportJob,
  ExportSpec,
  ExportStatus,
  ExportTemplate,
} from "@/features/export/types";

const STATUS_STYLES: Record<ExportStatus, string> = {
  pending: "bg-amber-100 text-amber-700",
  processing: "bg-sky-100 text-sky-700",
  completed: "bg-emerald-100 text-emerald-700",
  failed: "bg-red-100 text-red-700",
};

// Record-level columns that are always exportable alongside user fields.
const META_COLUMNS = [
  { key: "id", label: "Record ID" },
  { key: "created_at", label: "Created At" },
  { key: "updated_at", label: "Updated At" },
];

function StatusBadge({ status }: { status: ExportStatus }) {
  return (
    <span
      className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold capitalize ${STATUS_STYLES[status]}`}
    >
      {status}
    </span>
  );
}

function apiError(err: unknown, fallback: string): string {
  if (axios.isAxiosError(err)) return err.response?.data?.message ?? fallback;
  return fallback;
}

function formatBytes(n: number): string {
  if (n < 1024) return `${n} B`;
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`;
  return `${(n / (1024 * 1024)).toFixed(1)} MB`;
}

export default function ExportsPage() {
  const [modules, setModules] = useState<ModuleSummary[]>([]);
  const [moduleId, setModuleId] = useState("");
  const [fields, setFields] = useState<ModuleField[]>([]);

  const [format, setFormat] = useState<ExportFormat>("csv");
  const [selected, setSelected] = useState<string[]>([]);
  const [search, setSearch] = useState("");

  const [templates, setTemplates] = useState<ExportTemplate[]>([]);
  const [templateName, setTemplateName] = useState("");
  const [exports, setExports] = useState<ExportJob[]>([]);
  const [reloadKey, setReloadKey] = useState(0);
  const [busy, setBusy] = useState(false);

  // The ordered set of selectable columns: user fields, then meta columns.
  const columnOptions = useMemo(
    () => [
      ...fields.map((f) => ({ key: f.api_name, label: f.label, type: f.field_type })),
      ...META_COLUMNS.map((m) => ({ key: m.key, label: m.label, type: "meta" })),
    ],
    [fields]
  );

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const mods = await getModules();
        if (active) {
          setModules(
            mods.filter((m) => m.storage_strategy === "dynamic" && m.is_enabled)
          );
        }
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
      try {
        const f = await getModuleFields(moduleId);
        if (active) {
          setFields(f);
          setSelected(f.filter((x) => x.is_visible).map((x) => x.api_name));
        }
      } catch {
        toast.error("Failed to load fields");
      }
    })();
    return () => {
      active = false;
    };
  }, [moduleId]);

  useEffect(() => {
    if (!moduleId) return;
    let active = true;
    (async () => {
      try {
        const [tpls, result] = await Promise.all([
          listTemplates(moduleId),
          listExports(moduleId, { page_size: 20 }),
        ]);
        if (active) {
          setTemplates(tpls);
          setExports(result.exports);
        }
      } catch {
        toast.error("Failed to load export history");
      }
    })();
    return () => {
      active = false;
    };
  }, [moduleId, reloadKey]);

  const hasActive = exports.some(
    (j) => j.status === "pending" || j.status === "processing"
  );
  useEffect(() => {
    if (!hasActive) return;
    const timer = setInterval(() => setReloadKey((k) => k + 1), 3000);
    return () => clearInterval(timer);
  }, [hasActive]);

  function handleModuleChange(id: string) {
    setModuleId(id);
    setFields([]);
    setSelected([]);
    setTemplates([]);
    setExports([]);
    setTemplateName("");
  }

  function toggleColumn(key: string) {
    setSelected((prev) =>
      prev.includes(key) ? prev.filter((k) => k !== key) : [...prev, key]
    );
  }

  // Preserve the option order when serializing the selection.
  function orderedColumns(): string[] {
    return columnOptions.map((c) => c.key).filter((k) => selected.includes(k));
  }

  function buildSpec(): ExportSpec {
    return {
      format,
      columns: orderedColumns(),
      search: search.trim() || undefined,
    };
  }

  async function handleDownloadNow() {
    if (selected.length === 0) {
      toast.error("Select at least one column");
      return;
    }
    try {
      setBusy(true);
      await exportNow(moduleId, buildSpec());
      toast.success("Download started");
    } catch (err) {
      toast.error(apiError(err, "Failed to export"));
    } finally {
      setBusy(false);
    }
  }

  async function handleQueue() {
    if (selected.length === 0) {
      toast.error("Select at least one column");
      return;
    }
    try {
      setBusy(true);
      await createExport(moduleId, buildSpec());
      toast.success("Export queued");
      setReloadKey((k) => k + 1);
    } catch (err) {
      toast.error(apiError(err, "Failed to queue export"));
    } finally {
      setBusy(false);
    }
  }

  async function handleSaveTemplate() {
    if (!templateName.trim()) {
      toast.error("Give the template a name");
      return;
    }
    if (selected.length === 0) {
      toast.error("Select at least one column");
      return;
    }
    try {
      await createTemplate(moduleId, {
        name: templateName.trim(),
        format,
        columns: orderedColumns(),
      });
      toast.success("Template saved");
      setTemplateName("");
      setReloadKey((k) => k + 1);
    } catch (err) {
      toast.error(apiError(err, "Failed to save template"));
    }
  }

  function applyTemplate(template: ExportTemplate) {
    setFormat(template.format);
    setSearch("");
    if (template.columns.length > 0) {
      setSelected(template.columns);
    } else {
      setSelected(fields.filter((f) => f.is_visible).map((f) => f.api_name));
    }
    toast.success(`Applied "${template.name}"`);
  }

  async function handleDeleteTemplate(template: ExportTemplate) {
    try {
      await deleteTemplate(moduleId, template.id);
      toast.success("Template deleted");
      setReloadKey((k) => k + 1);
    } catch (err) {
      toast.error(apiError(err, "Failed to delete template"));
    }
  }

  async function handleDownloadJob(job: ExportJob) {
    try {
      await downloadExport(moduleId, job);
    } catch (err) {
      toast.error(apiError(err, "Failed to download"));
    }
  }

  return (
    <div className="space-y-8">
      <PageHeader
        badge="Data"
        title="Export Engine"
        description="Export any dynamic module to CSV or Excel. Download instantly for small sets, or queue an asynchronous export the worker builds in the background. Save reusable templates and re-download from history."
      />

      <FormSelect
        label="Source module"
        helperText="Only dynamic modules support the record runtime."
        value={moduleId}
        onChange={(e) => handleModuleChange(e.target.value)}
      >
        <option value="">Select a module…</option>
        {modules.map((m) => (
          <option key={m.id} value={m.id}>
            {m.plural_label}
          </option>
        ))}
      </FormSelect>

      {moduleId && (
        <div className="grid gap-8 lg:grid-cols-[1fr_1.1fr]">
          {/* Configure */}
          <div className="space-y-5 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <h3 className="text-sm font-semibold text-slate-900">Configure export</h3>

            <FormSelect
              label="Format"
              value={format}
              onChange={(e) => setFormat(e.target.value as ExportFormat)}
            >
              <option value="csv">CSV (.csv)</option>
              <option value="xlsx">Excel (.xlsx)</option>
            </FormSelect>

            <FormInput
              label="Search filter (optional)"
              placeholder="Only rows matching this text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />

            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <p className="text-sm font-semibold text-slate-700">Columns</p>
                <span className="text-xs text-slate-500">
                  {selected.length} selected
                </span>
              </div>
              <div className="max-h-64 space-y-1 overflow-y-auto rounded-2xl border border-slate-200 p-3">
                {columnOptions.map((c) => (
                  <label
                    key={c.key}
                    className="flex cursor-pointer items-center gap-3 rounded-lg px-2 py-1.5 hover:bg-slate-50"
                  >
                    <input
                      type="checkbox"
                      checked={selected.includes(c.key)}
                      onChange={() => toggleColumn(c.key)}
                      className="h-4 w-4 rounded border-slate-300 text-emerald-600 focus:ring-emerald-500"
                    />
                    <span className="flex-1 text-sm text-slate-700">{c.label}</span>
                    <span className="text-xs text-slate-400">{c.type}</span>
                  </label>
                ))}
              </div>
            </div>

            <div className="flex flex-wrap gap-3">
              <button
                type="button"
                onClick={handleDownloadNow}
                disabled={busy || selected.length === 0}
                className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:opacity-50"
              >
                <Download className="h-4 w-4" />
                Download now
              </button>
              <button
                type="button"
                onClick={handleQueue}
                disabled={busy || selected.length === 0}
                className="inline-flex items-center gap-2 rounded-full border border-slate-300 px-5 py-2.5 text-sm font-semibold text-slate-700 transition hover:bg-slate-100 disabled:opacity-50"
              >
                <FileDown className="h-4 w-4" />
                Queue export
              </button>
            </div>

            {/* Templates */}
            <div className="space-y-3 border-t border-slate-100 pt-5">
              <p className="text-sm font-semibold text-slate-700">Templates</p>

              <div className="flex gap-2">
                <input
                  type="text"
                  value={templateName}
                  onChange={(e) => setTemplateName(e.target.value)}
                  placeholder="Save current config as…"
                  className="flex-1 rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-800 focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-100"
                />
                <button
                  type="button"
                  onClick={handleSaveTemplate}
                  className="inline-flex items-center gap-2 rounded-xl bg-slate-800 px-3 py-2 text-sm font-medium text-white hover:bg-slate-900"
                >
                  <Save className="h-4 w-4" />
                  Save
                </button>
              </div>

              {templates.length > 0 && (
                <ul className="space-y-1">
                  {templates.map((t) => (
                    <li
                      key={t.id}
                      className="flex items-center justify-between rounded-lg border border-slate-200 px-3 py-2"
                    >
                      <button
                        type="button"
                        onClick={() => applyTemplate(t)}
                        className="min-w-0 flex-1 text-left"
                      >
                        <span className="truncate text-sm font-medium text-slate-700">
                          {t.name}
                        </span>
                        <span className="ml-2 text-xs uppercase text-slate-400">
                          {t.format}
                        </span>
                      </button>
                      <button
                        type="button"
                        onClick={() => handleDeleteTemplate(t)}
                        className="rounded-lg p-1.5 text-slate-400 hover:bg-red-50 hover:text-red-500"
                        aria-label="Delete template"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>

          {/* History */}
          <div className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex items-center justify-between gap-3">
              <h3 className="text-sm font-semibold text-slate-900">Export history</h3>
              <button
                type="button"
                onClick={() => setReloadKey((k) => k + 1)}
                className="rounded-lg border border-slate-200 p-1.5 text-slate-600 hover:bg-slate-100"
                aria-label="Refresh"
              >
                <RefreshCw className="h-4 w-4" />
              </button>
            </div>

            {exports.length === 0 ? (
              <p className="rounded-2xl border border-dashed border-slate-200 px-4 py-10 text-center text-sm text-slate-400">
                No exports yet for this module.
              </p>
            ) : (
              <ul className="space-y-3">
                {exports.map((job) => (
                  <li key={job.id} className="rounded-2xl border border-slate-200 p-4">
                    <div className="flex items-center justify-between gap-3">
                      <div className="min-w-0">
                        <p className="truncate text-sm font-medium text-slate-800">
                          {job.filename}
                        </p>
                        <p className="text-xs text-slate-400">
                          {new Date(job.created_at).toLocaleString()} ·{" "}
                          {job.format.toUpperCase()}
                        </p>
                      </div>
                      <StatusBadge status={job.status} />
                    </div>

                    <div className="mt-3 flex items-center justify-between">
                      <div className="flex flex-wrap gap-x-5 text-xs text-slate-600">
                        <span>
                          Rows: <strong>{job.row_count}</strong>
                        </span>
                        {job.byte_size > 0 && (
                          <span>
                            Size: <strong>{formatBytes(job.byte_size)}</strong>
                          </span>
                        )}
                      </div>
                      {job.status === "completed" && (
                        <button
                          type="button"
                          onClick={() => handleDownloadJob(job)}
                          className="inline-flex items-center gap-1.5 rounded-full bg-emerald-50 px-3 py-1 text-xs font-semibold text-emerald-700 hover:bg-emerald-100"
                        >
                          <Download className="h-3.5 w-3.5" />
                          Download
                        </button>
                      )}
                    </div>

                    {job.status === "failed" && job.error && (
                      <p className="mt-2 truncate text-xs text-red-500">{job.error}</p>
                    )}
                  </li>
                ))}
              </ul>
            )}

            <p className="text-xs text-slate-500">
              Queued exports are built asynchronously. Run the worker with
              <code className="mx-1 rounded bg-slate-100 px-1">make run-worker</code>
              so jobs move from queued to completed and become downloadable.
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
