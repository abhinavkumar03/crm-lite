"use client";

import {
  useEffect,
  useMemo,
  useState,
} from "react";

import { toast } from "sonner";

import PageHeader from "@/components/common/PageHeader";
import FormSelect from "@/components/common/form/FormSelect";

import DynamicForm from "@/features/metadata/components/DynamicForm";

import {
  getModuleFields,
  getModules,
  getValidationSchema,
  validateRecord,
} from "@/features/metadata/api";

import {
  FormValues,
  ModuleField,
  ModuleSummary,
  ValidationSchema,
  VisibilityRule,
} from "@/features/metadata/types";

import { errorListToMap } from "@/features/metadata/lib/validation";

export default function DynamicFormsPage() {
  const [modules, setModules] = useState<ModuleSummary[]>([]);
  const [moduleId, setModuleId] = useState("");

  const [fields, setFields] = useState<ModuleField[]>([]);
  const [schema, setSchema] = useState<ValidationSchema | null>(null);
  const [loadedId, setLoadedId] = useState("");

  const [preview, setPreview] = useState<FormValues>({});
  const [serverErrors, setServerErrors] = useState<Record<string, string>>({});
  const [conditionalDemo, setConditionalDemo] = useState(false);

  // Loading is derived rather than stored so no state is set synchronously
  // inside an effect (avoids cascading renders).
  const isLoading = !!moduleId && loadedId !== moduleId;

  useEffect(() => {
    (async () => {
      try {
        setModules(await getModules());
      } catch {
        toast.error("Failed to load modules");
      }
    })();
  }, []);

  useEffect(() => {
    if (!moduleId) return;

    (async () => {
      try {
        const [f, s] = await Promise.all([
          getModuleFields(moduleId),
          getValidationSchema(moduleId),
        ]);
        setFields(f);
        setSchema(s);
        setServerErrors({});
        setPreview({});
        setLoadedId(moduleId);
      } catch {
        toast.error("Failed to load module metadata");
      }
    })();
  }, [moduleId]);

  // Demo: when enabled, every field after the first is only shown once the first
  // field has a value — a live demonstration of the conditional-rendering engine.
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

  async function handleSubmit(values: FormValues) {
    const result = await validateRecord(moduleId, values);

    if (result.valid) {
      setServerErrors({});
      toast.success("Payload is valid (server-verified).");
      return;
    }

    setServerErrors(errorListToMap(result.errors));
    toast.error("Server validation failed.");
  }

  const selectedModule = modules.find((m) => m.id === moduleId);

  return (
    <div className="space-y-8">
      <PageHeader
        badge="Metadata Engine"
        title="Dynamic Forms"
        description="Forms generated entirely from module field metadata, with conditional rendering and validation driven by the backend schema. No form is hand-coded."
      />

      <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div className="grid gap-5 md:grid-cols-2">
          <FormSelect
            label="Module"
            helperText="Pick any module to generate its form from metadata."
            value={moduleId}
            onChange={(e) => setModuleId(e.target.value)}
          >
            <option value="">Select a module...</option>
            {modules.map((m) => (
              <option key={m.id} value={m.id}>
                {m.singular_label} ({m.api_name})
              </option>
            ))}
          </FormSelect>

          <label className="flex items-end gap-3 pb-3 text-sm font-medium text-slate-700">
            <input
              type="checkbox"
              checked={conditionalDemo}
              onChange={(e) => setConditionalDemo(e.target.checked)}
              className="h-4 w-4 accent-emerald-500"
            />
            Conditional demo: reveal fields after the first is filled
          </label>
        </div>
      </div>

      {isLoading && (
        <div className="rounded-3xl border border-slate-200 bg-white p-8 text-center text-slate-500 shadow-sm">
          Loading metadata...
        </div>
      )}

      {!isLoading && moduleId && schema && fields.length > 0 && (
        <div className="grid gap-8 lg:grid-cols-[1.4fr_1fr]">
          <DynamicForm
            key={`${moduleId}-${conditionalDemo}`}
            fields={fields}
            schema={schema}
            submitText="Validate"
            visibilityRules={visibilityRules}
            sectionTitle={selectedModule?.plural_label ?? "Record"}
            sectionDescription="Rendered dynamically from field metadata."
            externalErrors={serverErrors}
            onSubmit={handleSubmit}
            onChange={setPreview}
          />

          <aside className="space-y-4">
            <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
              <h3 className="mb-3 text-sm font-semibold text-slate-900">
                Live payload
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
                  Validation uses the compiled schema and the server
                  <code> /validate</code> endpoint.
                </li>
                <li>Conditional visibility is metadata-driven.</li>
              </ul>
            </div>
          </aside>
        </div>
      )}

      {!isLoading && moduleId && fields.length === 0 && schema && (
        <div className="rounded-3xl border border-slate-200 bg-white p-8 text-center text-slate-500 shadow-sm">
          This module has no fields yet.
        </div>
      )}
    </div>
  );
}
