"use client";

import {
  Suspense,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { toast } from "sonner";

import PageHeader from "@/components/common/PageHeader";
import FormSelect from "@/components/common/form/FormSelect";

import DynamicForm from "@/features/metadata/components/DynamicForm";

import { useDemo } from "@/features/demo/DemoProvider";
import {
  createRecord,
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

import { errorListToMap } from "@/features/metadata/lib/validation";

export default function DynamicFormsPage() {
  return (
    <Suspense
      fallback={
        <div className="rounded-3xl border border-slate-200 bg-white p-8 text-center text-slate-500">
          Loading forms...
        </div>
      }
    >
      <DynamicFormsInner />
    </Suspense>
  );
}

function DynamicFormsInner() {
  const demo = useDemo();
  const tutorialCreate =
    demo?.mode === "running" &&
    demo.currentStep?.step_key === "create_record";

  const router = useRouter();
  const searchParams = useSearchParams();
  const moduleFromQuery = searchParams.get("module") ?? "";

  const [modules, setModules] = useState<ModuleSummary[]>([]);
  const [moduleId, setModuleId] = useState(moduleFromQuery);

  const [fields, setFields] = useState<ModuleField[]>([]);
  const [schema, setSchema] = useState<ValidationSchema | null>(null);
  const [loadedId, setLoadedId] = useState("");

  const [preview, setPreview] = useState<FormValues>({});
  const [serverErrors, setServerErrors] = useState<Record<string, string>>({});
  const [conditionalDemo, setConditionalDemo] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const isLoading = !!moduleId && loadedId !== moduleId;

  useEffect(() => {
    (async () => {
      try {
        const all = await getModules();
        const dynamic = all.filter((m) => m.storage_strategy === "dynamic");
        setModules(dynamic);
        if (tutorialCreate) {
          const tutorial = dynamic.find((m) => m.api_name === "tutorial_lead");
          if (tutorial) {
            setModuleId(tutorial.id);
            const params = new URLSearchParams({ module: tutorial.id });
            router.replace(`/forms?${params.toString()}`);
          }
        }
      } catch {
        toast.error("Failed to load modules");
      }
    })();
  }, [tutorialCreate]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (moduleFromQuery && moduleFromQuery !== moduleId) {
      setModuleId(moduleFromQuery);
    }
  }, [moduleFromQuery]); // eslint-disable-line react-hooks/exhaustive-deps

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
    router.replace(qs ? `/forms?${qs}` : "/forms");
  }

  async function handleSubmit(values: FormValues) {
    if (!moduleId || submitting) return;
    setSubmitting(true);
    try {
      await createRecord(moduleId, values);
      setServerErrors({});
      setPreview({});
      toast.success("Record created");
      if (tutorialCreate) {
        // Stay briefly so demo validate can see the new row, then continue.
        await demo?.validate({ silent: true });
      }
      router.push(`/tables?module=${moduleId}`);
    } catch (err: unknown) {
      const axiosErr = err as {
        response?: { data?: { data?: { errors?: { field: string; message: string }[] }; message?: string } };
      };
      const fieldErrors = axiosErr.response?.data?.data?.errors;
      if (fieldErrors?.length) {
        setServerErrors(errorListToMap(fieldErrors));
        toast.error("Validation failed");
      } else {
        toast.error(
          axiosErr.response?.data?.message || "Failed to create record"
        );
      }
    } finally {
      setSubmitting(false);
    }
  }

  const selectedModule = modules.find((m) => m.id === moduleId);

  const tutorialInitialValues = useMemo<FormValues>(() => {
    if (!tutorialCreate) return {};
    const values: FormValues = {};
    for (const f of fields) {
      if (f.api_name === "company_name") {
        values.company_name = "Acme Tutorial Co";
      } else if (f.is_required && f.field_type === "text") {
        values[f.api_name] = "Tutorial";
      } else if (f.is_required && f.field_type === "email") {
        values[f.api_name] = "demo@example.com";
      } else if (f.is_required && (f.field_type === "number" || f.field_type === "currency")) {
        values[f.api_name] = 1;
      }
    }
    return values;
  }, [tutorialCreate, fields]);

  return (
    <div className="space-y-8" data-tutorial-surface="create-record">
      <PageHeader
        badge="Metadata Engine"
        title="Dynamic Forms"
        description="Create records with forms generated from module field metadata and backend validation."
      />

      {tutorialCreate && (
        <p className="rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-800">
          Tutorial Lead is selected and <strong>company_name</strong> is
          pre-filled. Click <strong>Create record</strong> to continue.
        </p>
      )}

      <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div className="grid gap-5 md:grid-cols-2">
          <FormSelect
            label="Module"
            helperText={
              tutorialCreate
                ? "Locked to Tutorial Lead for the walkthrough."
                : "Only dynamic modules are listed."
            }
            value={moduleId}
            disabled={tutorialCreate}
            onChange={(e) => selectModule(e.target.value)}
          >
            <option value="">Select a module...</option>
            {modules.map((m) => (
              <option key={m.id} value={m.id}>
                {m.singular_label} ({m.api_name})
              </option>
            ))}
          </FormSelect>

          {!tutorialCreate && (
            <label className="flex items-end gap-3 pb-3 text-sm font-medium text-slate-700">
              <input
                type="checkbox"
                checked={conditionalDemo}
                onChange={(e) => setConditionalDemo(e.target.checked)}
                className="h-4 w-4 accent-emerald-500"
              />
              Conditional demo: reveal fields after the first is filled
            </label>
          )}
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
            key={`${moduleId}-${conditionalDemo}-${tutorialCreate ? "tut" : "std"}`}
            fields={fields}
            schema={schema}
            submitText={submitting ? "Creating..." : "Create record"}
            initialValues={tutorialInitialValues}
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
                  Submit creates a row via{" "}
                  <code>POST /modules/:id/records</code>.
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
