"use client";

import {
  Suspense,
  useEffect,
  useMemo,
  useState,
} from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { toast } from "sonner";
import { ArrowRight, Monitor, Smartphone } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import FormSelect from "@/components/common/form/FormSelect";

import DynamicForm from "@/features/metadata/components/DynamicForm";

import {
  getModuleFields,
  getModules,
  getValidationSchema,
} from "@/features/metadata/api";

import {
  FormValues,
  ModuleField,
  ModuleSummary,
  ValidationSchema,
  VisibilityRule,
} from "@/features/metadata/types";

import { getFormLayout } from "@/features/workspace/api";
import type { FormLayout } from "@/features/workspace/types";

type PreviewMode = "desktop" | "mobile";

export default function FormDesignerPage() {
  return (
    <Suspense
      fallback={
        <div className="rounded-3xl border border-slate-200 bg-white p-8 text-center text-slate-500">
          Loading form designer...
        </div>
      }
    >
      <FormDesignerInner />
    </Suspense>
  );
}

function FormDesignerInner() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const moduleFromQuery = searchParams.get("module") ?? "";

  const [modules, setModules] = useState<ModuleSummary[]>([]);
  const [moduleId, setModuleId] = useState(moduleFromQuery);

  const [fields, setFields] = useState<ModuleField[]>([]);
  const [schema, setSchema] = useState<ValidationSchema | null>(null);
  const [formLayout, setFormLayout] = useState<FormLayout | null>(null);
  const [loadedId, setLoadedId] = useState("");

  const [preview, setPreview] = useState<FormValues>({});
  const [conditionalDemo, setConditionalDemo] = useState(false);
  const [previewMode, setPreviewMode] = useState<PreviewMode>("desktop");

  const isLoading = !!moduleId && loadedId !== moduleId;

  useEffect(() => {
    (async () => {
      try {
        const all = await getModules();
        const dynamic = all.filter((m) => m.storage_strategy === "dynamic");
        setModules(dynamic);
      } catch {
        toast.error("Failed to load modules");
      }
    })();
  }, []);

  useEffect(() => {
    if (moduleFromQuery && moduleFromQuery !== moduleId) {
      setModuleId(moduleFromQuery);
    }
  }, [moduleFromQuery]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (!moduleId) return;

    (async () => {
      try {
        const [f, s, layout] = await Promise.all([
          getModuleFields(moduleId),
          getValidationSchema(moduleId),
          getFormLayout(moduleId, "create"),
        ]);
        setFields(f);
        setSchema(s);
        setFormLayout(layout);
        setPreview({});
        setLoadedId(moduleId);
      } catch {
        toast.error("Failed to load module metadata");
      }
    })();
  }, [moduleId]);

  const visibilityRules = useMemo<VisibilityRule[]>(() => {
    if (!conditionalDemo) return [];

    const visible = [...fields]
      .filter((f) => f.is_visible)
      .sort((a, b) => a.sort_order - b.sort_order);

    if (visible.length < 2) return [];

    const [first, ...rest] = visible;
    return [
      {
        when: { field: first.api_name, operator: "not_empty" },
        effect: "show",
        targets: rest.map((f) => f.api_name),
      },
    ];
  }, [conditionalDemo, fields]);

  function selectModule(id: string) {
    setModuleId(id);
    const params = new URLSearchParams();
    if (id) params.set("module", id);
    const qs = params.toString();
    router.replace(qs ? `/settings/forms?${qs}` : "/settings/forms");
  }

  const selectedModule = modules.find((m) => m.id === moduleId);
  const moduleHref = selectedModule
    ? `/m/${selectedModule.api_name}?create=1`
    : "/dashboard";

  const requiredCount = fields.filter((f) => f.is_required && f.is_visible).length;
  const sectionCount = formLayout?.sections?.length ?? 0;

  return (
    <div className="space-y-8" data-tutorial-surface="form-designer">
      <PageHeader
        badge="Metadata Engine"
        title="Form Designer"
        description="Preview how create forms are generated from field metadata, sections, and layout — without saving records."
      />

      <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div className="grid gap-5 md:grid-cols-2">
          <FormSelect
            label="Module"
            helperText="Only dynamic modules are listed."
            value={moduleId}
            onChange={(e) => selectModule(e.target.value)}
          >
            <option value="">Select a module...</option>
            {modules.map((m) => (
              <option key={m.id} value={m.id}>
                {m.singular_label} ({m.api_name})
              </option>
            ))}
          </FormSelect>

          <div className="flex flex-col justify-end gap-3 pb-1">
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                Preview
              </span>
              <button
                type="button"
                onClick={() => setPreviewMode("desktop")}
                className={`inline-flex items-center gap-1.5 rounded-full px-3 py-1.5 text-xs font-semibold transition ${
                  previewMode === "desktop"
                    ? "bg-emerald-500 text-white"
                    : "bg-slate-100 text-slate-600 hover:bg-slate-200"
                }`}
              >
                <Monitor size={14} />
                Desktop
              </button>
              <button
                type="button"
                onClick={() => setPreviewMode("mobile")}
                className={`inline-flex items-center gap-1.5 rounded-full px-3 py-1.5 text-xs font-semibold transition ${
                  previewMode === "mobile"
                    ? "bg-emerald-500 text-white"
                    : "bg-slate-100 text-slate-600 hover:bg-slate-200"
                }`}
              >
                <Smartphone size={14} />
                Mobile
              </button>
            </div>
            <label className="flex items-center gap-3 text-sm font-medium text-slate-700">
              <input
                type="checkbox"
                checked={conditionalDemo}
                onChange={(e) => setConditionalDemo(e.target.checked)}
                className="h-4 w-4 accent-emerald-500"
              />
              Test conditional visibility (reveal fields after the first is filled)
            </label>
          </div>
        </div>

        {selectedModule && (
          <div className="mt-4 flex flex-wrap gap-3 text-xs text-slate-500">
            <span className="rounded-full bg-slate-100 px-3 py-1">
              {fields.length} fields
            </span>
            <span className="rounded-full bg-slate-100 px-3 py-1">
              {requiredCount} required
            </span>
            <span className="rounded-full bg-slate-100 px-3 py-1">
              {sectionCount} sections
            </span>
          </div>
        )}
      </div>

      {isLoading && (
        <div className="rounded-3xl border border-slate-200 bg-white p-8 text-center text-slate-500 shadow-sm">
          Loading metadata...
        </div>
      )}

      {!isLoading && moduleId && schema && fields.length > 0 && (
        <div className="grid gap-8 lg:grid-cols-[1.4fr_1fr]">
          <div
            className={
              previewMode === "mobile"
                ? "mx-auto w-full max-w-sm"
                : "w-full"
            }
          >
            <DynamicForm
              key={`${moduleId}-${conditionalDemo}-${previewMode}`}
              fields={fields}
              schema={schema}
              formLayout={formLayout}
              previewOnly
              visibilityRules={visibilityRules}
              sectionTitle={selectedModule?.singular_label ?? "Record"}
              sectionDescription="Rendered dynamically from field metadata."
              onChange={setPreview}
              footerSlot={
                <Link
                  href={moduleHref}
                  className="inline-flex items-center gap-2 text-sm font-semibold text-emerald-700 hover:text-emerald-800"
                >
                  Open {selectedModule?.singular_label ?? "module"} module
                  <ArrowRight size={16} />
                </Link>
              }
            />
          </div>

          <aside className="space-y-4">
            <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
              <h3 className="mb-3 text-sm font-semibold text-slate-900">
                Live payload preview
              </h3>
              <pre className="max-h-96 overflow-auto rounded-2xl bg-slate-900 p-4 text-xs text-emerald-200">
                {JSON.stringify(preview, null, 2)}
              </pre>
            </div>

            <div className="rounded-3xl border border-slate-200 bg-white p-6 text-sm text-slate-500 shadow-sm">
              <h3 className="mb-2 text-sm font-semibold text-slate-900">
                How this works
              </h3>
              <ul className="list-disc space-y-1 pl-5">
                <li>Fields come from <code>GET /modules/:id/fields</code>.</li>
                <li>
                  Form sections come from{" "}
                  <code>GET /modules/:id/layouts/form</code>.
                </li>
                <li>
                  Edit sections and fields in Settings → Fields — this page only
                  previews the result.
                </li>
                <li>
                  Create real records from the module listing page (Add …).
                </li>
              </ul>
            </div>
          </aside>
        </div>
      )}

      {!isLoading && moduleId && fields.length === 0 && schema && (
        <div className="rounded-3xl border border-slate-200 bg-white p-8 text-center text-slate-500 shadow-sm">
          This module has no fields yet. Add fields in Settings → Fields.
        </div>
      )}
    </div>
  );
}
