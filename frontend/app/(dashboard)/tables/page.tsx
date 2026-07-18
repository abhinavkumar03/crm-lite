"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Plus, X } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import FormSelect from "@/components/common/form/FormSelect";

import DynamicForm from "@/features/metadata/components/DynamicForm";
import DynamicTable from "@/features/metadata/components/DynamicTable";
import ViewBar from "@/features/metadata/components/ViewBar";

import { useDynamicTable } from "@/features/metadata/hooks/useDynamicTable";

import {
  createView,
  deleteView,
  getModuleFields,
  getModules,
  getValidationSchema,
  getViews,
  setDefaultView,
  validateRecord,
} from "@/features/metadata/api";

import {
  FormValues,
  ModuleField,
  ModuleSummary,
  SavedView,
  TableRow,
  ValidationSchema,
} from "@/features/metadata/types";

import { errorListToMap } from "@/features/metadata/lib/validation";

// Builds a few illustrative rows so the table has data to sort/filter/paginate
// before the record runtime (Phase 10) is wired up.
function makeSampleRows(fields: ModuleField[], count: number): TableRow[] {
  return Array.from({ length: count }, (_, i) => {
    const row: TableRow = { id: `sample-${i + 1}` };
    for (const field of fields) {
      switch (field.field_type) {
        case "number":
        case "currency":
          row[field.api_name] = (i + 1) * 100;
          break;
        case "boolean":
        case "checkbox":
          row[field.api_name] = i % 2 === 0;
          break;
        case "email":
          row[field.api_name] = `user${i + 1}@example.com`;
          break;
        case "date":
        case "datetime":
          row[field.api_name] = new Date(
            Date.now() - i * 86400000
          ).toISOString();
          break;
        case "dropdown":
        case "radio":
        case "user":
        case "lookup":
          row[field.api_name] =
            field.options[i % Math.max(1, field.options.length)]?.value ?? "";
          break;
        case "multiselect":
          row[field.api_name] = field.options.slice(0, 1).map((o) => o.value);
          break;
        default:
          row[field.api_name] = `${field.label} ${i + 1}`;
      }
    }
    return row;
  });
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
  const [rows, setRows] = useState<TableRow[]>(() => makeSampleRows(fields, 8));
  const [views, setViews] = useState<SavedView[]>(initialViews);
  const [activeViewId, setActiveViewId] = useState<string | null>(null);
  const [adding, setAdding] = useState(false);
  const [addErrors, setAddErrors] = useState<Record<string, string>>({});

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
    const result = await validateRecord(moduleId, values);
    if (!result.valid) {
      setAddErrors(errorListToMap(result.errors));
      toast.error("Server validation failed.");
      return;
    }
    setAddErrors({});
    setRows((prev) => [
      { id: crypto.randomUUID(), ...(values as TableRow) },
      ...prev,
    ]);
    setAdding(false);
    toast.success("Record added");
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
            sectionDescription="Validated by the backend before it is added to the table."
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
      />
    </div>
  );
}

export default function DynamicTablesPage() {
  const [modules, setModules] = useState<ModuleSummary[]>([]);
  const [moduleId, setModuleId] = useState("");

  const [fields, setFields] = useState<ModuleField[]>([]);
  const [schema, setSchema] = useState<ValidationSchema | null>(null);
  const [views, setViews] = useState<SavedView[]>([]);
  const [loadedId, setLoadedId] = useState("");

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

  const selectedModule = modules.find((m) => m.id === moduleId);

  return (
    <div className="space-y-8">
      <PageHeader
        badge="Metadata Engine"
        title="Dynamic Tables"
        description="Metadata-driven tables with client-side sorting, filtering and pagination, plus saved views persisted per module. No table is hand-coded."
      />

      <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <FormSelect
          label="Module"
          helperText="Pick any module to render its records table from metadata."
          value={moduleId}
          onChange={(e) => setModuleId(e.target.value)}
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
