"use client";

import { Suspense, useEffect, useMemo, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { toast } from "sonner";
import { Plus, X } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import FormSelect from "@/components/common/form/FormSelect";

import DynamicForm from "@/features/metadata/components/DynamicForm";
import DynamicTable from "@/features/metadata/components/DynamicTable";
import ViewBar from "@/features/metadata/components/ViewBar";

import { useDynamicTable } from "@/features/metadata/hooks/useDynamicTable";

import {
  createRecord,
  createView,
  deleteRecord,
  deleteView,
  getModuleFields,
  getModules,
  getValidationSchema,
  getViews,
  listRecords,
  setDefaultView,
} from "@/features/metadata/api";

import {
  FieldError,
  FormValues,
  ModuleField,
  ModuleSummary,
  RecordResponse,
  SavedView,
  TableRow,
  ValidationSchema,
} from "@/features/metadata/types";

import { errorListToMap } from "@/features/metadata/lib/validation";

// Flattens a record into a table row: field values plus resolved relation labels
// (so lookup/user columns show a name instead of a raw id).
function recordToRow(rec: RecordResponse, moduleFields: ModuleField[]): TableRow {
  const row: TableRow = { id: rec.id };
  for (const field of moduleFields) {
    row[field.api_name] = rec.data[field.api_name] ?? null;
  }
  if (rec.relations) {
    for (const [key, ref] of Object.entries(rec.relations)) {
      row[key] = ref.label;
    }
  }
  return row;
}

interface WorkspaceProps {
  moduleId: string;
  moduleLabel: string;
  fields: ModuleField[];
  schema: ValidationSchema;
  initialViews: SavedView[];
}

// RecordsWorkspace is keyed by moduleId so its table state (columns, sort,
// filters, pagination) resets whenever a different module is selected.
function RecordsWorkspace({
  moduleId,
  moduleLabel,
  fields,
  schema,
  initialViews,
}: WorkspaceProps) {
  const router = useRouter();
  const [records, setRecords] = useState<RecordResponse[]>([]);
  const [views, setViews] = useState<SavedView[]>(initialViews);
  const [activeViewId, setActiveViewId] = useState<string | null>(null);
  const [adding, setAdding] = useState(false);
  const [addErrors, setAddErrors] = useState<Record<string, string>>({});
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        // Fetch a generous page and let the table handle in-view sort/filter.
        const result = await listRecords(moduleId, {
          page_size: 100,
          expand: true,
        });
        if (active) setRecords(result.records);
      } catch {
        toast.error("Failed to load records");
      }
    })();
    return () => {
      active = false;
    };
  }, [moduleId, reloadKey]);

  const rows = useMemo(
    () => records.map((rec) => recordToRow(rec, fields)),
    [records, fields]
  );

  const table = useDynamicTable({ rows, fields });

  async function refreshViews() {
    try {
      setViews(await getViews(moduleId));
    } catch {
      toast.error("Failed to refresh views");
    }
  }

  async function handleSaveView(name: string, isPublic: boolean) {
    try {
      const created = await createView(moduleId, {
        name,
        columns: table.currentConfig.columns,
        filters: table.currentConfig.filters,
        sort: table.currentConfig.sort,
        is_public: isPublic,
      });
      setActiveViewId(created.id);
      await refreshViews();
      toast.success(`View "${name}" saved`);
    } catch {
      toast.error("Failed to save view");
    }
  }

  async function handleSetDefault(view: SavedView) {
    try {
      await setDefaultView(moduleId, view.id);
      await refreshViews();
      toast.success(`"${view.name}" is now the default view`);
    } catch {
      toast.error("Failed to set default view");
    }
  }

  async function handleDeleteView(view: SavedView) {
    try {
      await deleteView(moduleId, view.id);
      if (activeViewId === view.id) setActiveViewId(null);
      await refreshViews();
      toast.success(`View "${view.name}" deleted`);
    } catch {
      toast.error("Failed to delete view");
    }
  }

  async function handleAddRecord(values: FormValues) {
    try {
      await createRecord(moduleId, values);
      setAddErrors({});
      setAdding(false);
      setReloadKey((k) => k + 1);
      toast.success("Record created");
    } catch (err: unknown) {
      const fieldErrors = extractFieldErrors(err);
      if (fieldErrors) {
        setAddErrors(errorListToMap(fieldErrors));
        toast.error("Server validation failed.");
      } else {
        toast.error("Failed to create record");
      }
    }
  }

  async function handleDeleteRecord(row: TableRow) {
    const id = row.id as string;
    try {
      await deleteRecord(moduleId, id);
      setReloadKey((k) => k + 1);
      toast.success("Record deleted");
    } catch {
      toast.error("Failed to delete record");
    }
  }

  return (
    <div className="space-y-5">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <ViewBar
          views={views}
          activeViewId={activeViewId}
          onApply={(view) => {
            table.applyView(view);
            setActiveViewId(view.id);
          }}
          onSaveCurrent={handleSaveView}
          onSetDefault={handleSetDefault}
          onDelete={handleDeleteView}
        />

        <button
          type="button"
          onClick={() => setAdding((v) => !v)}
          className="inline-flex items-center gap-1.5 rounded-full bg-emerald-500 px-4 py-2 text-sm font-semibold text-white hover:bg-emerald-600"
        >
          {adding ? <X className="h-4 w-4" /> : <Plus className="h-4 w-4" />}
          {adding ? "Close" : "Add record"}
        </button>
      </div>

      {adding && (
        <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <DynamicForm
            fields={fields}
            schema={schema}
            submitText="Add record"
            sectionTitle={`New ${moduleLabel}`}
            sectionDescription="Validated by the backend, then persisted via the record runtime."
            externalErrors={addErrors}
            onSubmit={handleAddRecord}
          />
        </div>
      )}

      <DynamicTable
        fields={fields}
        columns={table.columns}
        rows={table.result.rows}
        sort={table.sort}
        onToggleSort={table.toggleSort}
        filters={table.filters}
        onFilter={table.setFilter}
        page={table.result.page}
        totalPages={table.result.totalPages}
        total={table.result.total}
        pageSize={table.pageSize}
        onPage={table.setPage}
        onPageSize={table.setPageSize}
        onDeleteRow={handleDeleteRecord}
        onRowClick={(row) => {
          const id = row.id as string;
          if (id) router.push(`/tables/${moduleId}/${id}`);
        }}
      />
    </div>
  );
}

function extractFieldErrors(err: unknown): FieldError[] | null {
  if (
    typeof err === "object" &&
    err !== null &&
    "response" in err &&
    typeof (err as { response?: unknown }).response === "object"
  ) {
    const response = (err as { response?: { data?: { errors?: unknown } } })
      .response;
    const errors = response?.data?.errors;
    if (Array.isArray(errors)) return errors as FieldError[];
  }
  return null;
}

export default function DynamicTablesPage() {
  return (
    <Suspense
      fallback={
        <div className="rounded-3xl border border-slate-200 bg-white p-8 text-center text-slate-500">
          Loading tables...
        </div>
      }
    >
      <DynamicTablesInner />
    </Suspense>
  );
}

function DynamicTablesInner() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const moduleFromQuery = searchParams.get("module") ?? "";

  const [modules, setModules] = useState<ModuleSummary[]>([]);
  const [moduleId, setModuleId] = useState(moduleFromQuery);

  const [fields, setFields] = useState<ModuleField[]>([]);
  const [schema, setSchema] = useState<ValidationSchema | null>(null);
  const [views, setViews] = useState<SavedView[]>([]);
  const [loadedId, setLoadedId] = useState("");

  const isLoading = !!moduleId && loadedId !== moduleId;

  useEffect(() => {
    (async () => {
      try {
        const all = await getModules();
        setModules(all.filter((m) => m.storage_strategy === "dynamic"));
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
        const [f, s, v] = await Promise.all([
          getModuleFields(moduleId),
          getValidationSchema(moduleId),
          getViews(moduleId),
        ]);
        setFields(f);
        setSchema(s);
        setViews(v);
        setLoadedId(moduleId);
      } catch {
        toast.error("Failed to load module metadata");
      }
    })();
  }, [moduleId]);

  function selectModule(id: string) {
    setModuleId(id);
    const params = new URLSearchParams();
    if (id) params.set("module", id);
    const qs = params.toString();
    router.replace(qs ? `/tables?${qs}` : "/tables");
  }

  const selectedModule = modules.find((m) => m.id === moduleId);

  return (
    <div className="space-y-8">
      <PageHeader
        badge="Metadata Engine"
        title="Dynamic Tables"
        description="Metadata-driven tables backed by the generic record runtime: create, query and delete records that live entirely in JSONB, with saved views per module."
      />

      <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <FormSelect
          label="Module"
          helperText="Dynamic (JSONB-backed) modules are served by the record runtime."
          value={moduleId}
          onChange={(e) => selectModule(e.target.value)}
        >
          <option value="">Select a module...</option>
          {modules.map((m) => (
            <option key={m.id} value={m.id}>
              {m.plural_label} ({m.api_name})
            </option>
          ))}
        </FormSelect>
      </div>

      {isLoading && (
        <div className="rounded-3xl border border-slate-200 bg-white p-8 text-center text-slate-500 shadow-sm">
          Loading metadata...
        </div>
      )}

      {!isLoading && moduleId && schema && fields.length > 0 && (
        <RecordsWorkspace
          key={moduleId}
          moduleId={moduleId}
          moduleLabel={selectedModule?.singular_label ?? "record"}
          fields={fields}
          schema={schema}
          initialViews={views}
        />
      )}

      {!isLoading && moduleId && fields.length === 0 && schema && (
        <div className="rounded-3xl border border-slate-200 bg-white p-8 text-center text-slate-500 shadow-sm">
          This module has no fields yet.
        </div>
      )}
    </div>
  );
}
