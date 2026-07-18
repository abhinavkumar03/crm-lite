"use client";

import { useEffect, useMemo, useState } from "react";
import axios from "axios";
import { toast } from "sonner";
import { RefreshCw, Upload, FileSpreadsheet, ArrowRight, X } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import FormSelect from "@/components/common/form/FormSelect";

import { getModules, getModuleFields } from "@/features/metadata/api";
import { ModuleField, ModuleSummary } from "@/features/metadata/types";
import {
  analyzeImport,
  createImport,
  listImports,
} from "@/features/import/api";
import { AnalyzeResult, ImportJob, ImportStatus } from "@/features/import/types";

const STATUS_STYLES: Record<ImportStatus, string> = {
  pending: "bg-amber-100 text-amber-700",
  processing: "bg-sky-100 text-sky-700",
  completed: "bg-emerald-100 text-emerald-700",
  failed: "bg-red-100 text-red-700",
};

const SKIP = "";

function StatusBadge({ status }: { status: ImportStatus }) {
  return (
    <span
      className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold capitalize ${STATUS_STYLES[status]}`}
    >
      {status}
    </span>
  );
}

// A field is importable if a user can write to it. Read-only, system and derived
// (formula) fields are never valid import targets.
function isWritable(f: ModuleField): boolean {
  return !f.is_read_only && !f.is_system && f.field_type !== "formula";
}

function apiError(err: unknown, fallback: string): string {
  if (axios.isAxiosError(err)) {
    return err.response?.data?.message ?? fallback;
  }
  return fallback;
}

export default function ImportsPage() {
  const [modules, setModules] = useState<ModuleSummary[]>([]);
  const [moduleId, setModuleId] = useState("");
  const [fields, setFields] = useState<ModuleField[]>([]);

  const [file, setFile] = useState<File | null>(null);
  const [analysis, setAnalysis] = useState<AnalyzeResult | null>(null);
  const [analyzing, setAnalyzing] = useState(false);
  const [importing, setImporting] = useState(false);

  // fieldSource maps a field api_name -> chosen source column header ("" = skip).
  const [fieldSource, setFieldSource] = useState<Record<string, string>>({});

  const [imports, setImports] = useState<ImportJob[]>([]);
  const [reloadKey, setReloadKey] = useState(0);

  const writableFields = useMemo(() => fields.filter(isWritable), [fields]);

  // Load dynamic modules once (only dynamic modules use the record runtime).
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

  // Load the selected module's fields.
  useEffect(() => {
    if (!moduleId) return;
    let active = true;
    (async () => {
      try {
        const f = await getModuleFields(moduleId);
        if (active) setFields(f);
      } catch {
        toast.error("Failed to load fields");
      }
    })();
    return () => {
      active = false;
    };
  }, [moduleId]);

  // Load import history for the selected module.
  useEffect(() => {
    if (!moduleId) return;
    let active = true;
    (async () => {
      try {
        const result = await listImports(moduleId, { page_size: 20 });
        if (active) setImports(result.imports);
      } catch {
        toast.error("Failed to load import history");
      }
    })();
    return () => {
      active = false;
    };
  }, [moduleId, reloadKey]);

  // Auto-refresh while any job is still running so progress is live.
  const hasActive = imports.some(
    (j) => j.status === "pending" || j.status === "processing"
  );
  useEffect(() => {
    if (!hasActive) return;
    const timer = setInterval(() => setReloadKey((k) => k + 1), 3000);
    return () => clearInterval(timer);
  }, [hasActive]);

  function resetWizard() {
    setFile(null);
    setAnalysis(null);
    setFieldSource({});
  }

  function handleModuleChange(id: string) {
    setModuleId(id);
    setFields([]);
    setImports([]);
    resetWizard();
  }

  async function handleFile(selected: File | null) {
    if (!selected || !moduleId) return;
    setFile(selected);
    setAnalysis(null);
    try {
      setAnalyzing(true);
      const result = await analyzeImport(moduleId, selected);
      setAnalysis(result);
      // Invert the suggested header->field mapping into field->header for the UI.
      const initial: Record<string, string> = {};
      Object.entries(result.suggested_mapping).forEach(([header, api]) => {
        initial[api] = header;
      });
      setFieldSource(initial);
    } catch (err) {
      toast.error(apiError(err, "Failed to analyze file"));
      setFile(null);
    } finally {
      setAnalyzing(false);
    }
  }

  const mappedCount = useMemo(
    () => writableFields.filter((f) => fieldSource[f.api_name]).length,
    [writableFields, fieldSource]
  );

  async function handleStart() {
    if (!file || !moduleId) return;

    const mapping: Record<string, string> = {};
    writableFields.forEach((f) => {
      const header = fieldSource[f.api_name];
      if (header) mapping[header] = f.api_name;
    });

    if (Object.keys(mapping).length === 0) {
      toast.error("Map at least one column to a field");
      return;
    }

    try {
      setImporting(true);
      await createImport(moduleId, file, mapping);
      toast.success("Import started");
      resetWizard();
      setReloadKey((k) => k + 1);
    } catch (err) {
      toast.error(apiError(err, "Failed to start import"));
    } finally {
      setImporting(false);
    }
  }

  return (
    <div className="space-y-8">
      <PageHeader
        badge="Data"
        title="Import Engine"
        description="Bring records into any dynamic module from CSV or Excel. Columns are auto-matched to fields, every row runs through the same validation engine as the API, and processing happens asynchronously in the worker with a full error report."
      />

      <FormSelect
        label="Target module"
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
          {/* Upload + mapping */}
          <div className="space-y-5 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <h3 className="text-sm font-semibold text-slate-900">
              1. Upload a file
            </h3>

            {!analysis ? (
              <label className="flex cursor-pointer flex-col items-center justify-center gap-3 rounded-2xl border-2 border-dashed border-slate-300 bg-slate-50 px-6 py-10 text-center transition hover:border-emerald-400 hover:bg-emerald-50/40">
                <Upload className="h-8 w-8 text-slate-400" />
                <div>
                  <p className="text-sm font-medium text-slate-700">
                    {analyzing ? "Analyzing…" : "Click to choose a .csv or .xlsx file"}
                  </p>
                  <p className="text-xs text-slate-400">Max 10 MiB, 5000 rows</p>
                </div>
                <input
                  type="file"
                  accept=".csv,.xlsx"
                  className="hidden"
                  disabled={analyzing}
                  onChange={(e) => handleFile(e.target.files?.[0] ?? null)}
                />
              </label>
            ) : (
              <div className="flex items-center justify-between rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                <div className="flex items-center gap-3">
                  <FileSpreadsheet className="h-5 w-5 text-emerald-600" />
                  <div>
                    <p className="text-sm font-medium text-slate-800">
                      {file?.name}
                    </p>
                    <p className="text-xs text-slate-500">
                      {analysis.row_count} rows · {analysis.headers.length} columns
                    </p>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={resetWizard}
                  className="rounded-lg p-1.5 text-slate-500 hover:bg-slate-200"
                  aria-label="Remove file"
                >
                  <X className="h-4 w-4" />
                </button>
              </div>
            )}

            {analysis && (
              <>
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-semibold text-slate-900">
                    2. Map columns to fields
                  </h3>
                  <span className="text-xs text-slate-500">
                    {mappedCount} of {writableFields.length} mapped
                  </span>
                </div>

                <div className="space-y-3">
                  {writableFields.map((f) => (
                    <div
                      key={f.id}
                      className="grid grid-cols-[1fr_auto_1fr] items-center gap-2"
                    >
                      <div className="min-w-0">
                        <p className="truncate text-sm font-medium text-slate-700">
                          {f.label}
                          {f.is_required && (
                            <span className="ml-1 text-red-500">*</span>
                          )}
                        </p>
                        <p className="truncate text-xs text-slate-400">
                          {f.api_name} · {f.field_type}
                        </p>
                      </div>
                      <ArrowRight className="h-4 w-4 shrink-0 text-slate-300" />
                      <select
                        value={fieldSource[f.api_name] ?? SKIP}
                        onChange={(e) =>
                          setFieldSource((prev) => ({
                            ...prev,
                            [f.api_name]: e.target.value,
                          }))
                        }
                        className="w-full rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-800 focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-100"
                      >
                        <option value={SKIP}>— Skip —</option>
                        {analysis.headers.map((h) => (
                          <option key={h} value={h}>
                            {h}
                          </option>
                        ))}
                      </select>
                    </div>
                  ))}
                </div>

                <button
                  type="button"
                  onClick={handleStart}
                  disabled={importing || mappedCount === 0}
                  className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:opacity-50"
                >
                  <Upload className="h-4 w-4" />
                  {importing ? "Starting…" : `Import ${analysis.row_count} rows`}
                </button>
              </>
            )}

            {analysis && analysis.sample_rows.length > 0 && (
              <div className="space-y-2">
                <h4 className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                  Preview
                </h4>
                <div className="overflow-x-auto rounded-xl border border-slate-200">
                  <table className="w-full text-xs">
                    <thead>
                      <tr className="bg-slate-50 text-left text-slate-600">
                        {analysis.headers.map((h) => (
                          <th key={h} className="whitespace-nowrap px-3 py-2 font-semibold">
                            {h}
                          </th>
                        ))}
                      </tr>
                    </thead>
                    <tbody>
                      {analysis.sample_rows.map((row, i) => (
                        <tr key={i} className="border-t border-slate-100">
                          {analysis.headers.map((h) => (
                            <td key={h} className="whitespace-nowrap px-3 py-1.5 text-slate-600">
                              {row[h]}
                            </td>
                          ))}
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </div>

          {/* History */}
          <div className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex items-center justify-between gap-3">
              <h3 className="text-sm font-semibold text-slate-900">
                Import history
              </h3>
              <button
                type="button"
                onClick={() => setReloadKey((k) => k + 1)}
                className="rounded-lg border border-slate-200 p-1.5 text-slate-600 hover:bg-slate-100"
                aria-label="Refresh"
              >
                <RefreshCw className="h-4 w-4" />
              </button>
            </div>

            {imports.length === 0 ? (
              <p className="rounded-2xl border border-dashed border-slate-200 px-4 py-10 text-center text-sm text-slate-400">
                No imports yet for this module.
              </p>
            ) : (
              <ul className="space-y-3">
                {imports.map((job) => (
                  <li
                    key={job.id}
                    className="rounded-2xl border border-slate-200 p-4"
                  >
                    <div className="flex items-center justify-between gap-3">
                      <div className="min-w-0">
                        <p className="truncate text-sm font-medium text-slate-800">
                          {job.filename}
                        </p>
                        <p className="text-xs text-slate-400">
                          {new Date(job.created_at).toLocaleString()}
                        </p>
                      </div>
                      <StatusBadge status={job.status} />
                    </div>

                    <div className="mt-3 flex flex-wrap gap-x-5 gap-y-1 text-xs text-slate-600">
                      <span>
                        Total: <strong>{job.total_rows}</strong>
                      </span>
                      <span className="text-emerald-600">
                        Success: <strong>{job.success_rows}</strong>
                      </span>
                      <span className={job.error_rows > 0 ? "text-red-600" : ""}>
                        Errors: <strong>{job.error_rows}</strong>
                      </span>
                    </div>

                    {job.errors.length > 0 && (
                      <details className="mt-3">
                        <summary className="cursor-pointer text-xs font-medium text-red-600">
                          View error report ({job.errors.length})
                        </summary>
                        <div className="mt-2 max-h-48 overflow-y-auto rounded-lg bg-red-50/60 p-2">
                          <table className="w-full text-xs">
                            <thead>
                              <tr className="text-left text-red-500">
                                <th className="px-2 py-1 font-semibold">Row</th>
                                <th className="px-2 py-1 font-semibold">Field</th>
                                <th className="px-2 py-1 font-semibold">Message</th>
                              </tr>
                            </thead>
                            <tbody>
                              {job.errors.map((e, i) => (
                                <tr key={i} className="border-t border-red-100">
                                  <td className="px-2 py-1 text-slate-600">{e.row}</td>
                                  <td className="px-2 py-1 text-slate-600">
                                    {e.field || "—"}
                                  </td>
                                  <td className="px-2 py-1 text-slate-700">
                                    {e.message}
                                  </td>
                                </tr>
                              ))}
                            </tbody>
                          </table>
                        </div>
                      </details>
                    )}
                  </li>
                ))}
              </ul>
            )}

            <p className="text-xs text-slate-500">
              Rows are processed asynchronously. Run the worker with
              <code className="mx-1 rounded bg-slate-100 px-1">make run-worker</code>
              to move imports from queued to completed.
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
