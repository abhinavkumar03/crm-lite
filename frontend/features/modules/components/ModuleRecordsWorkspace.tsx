"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Plus, X } from "lucide-react";

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
  getValidationSchema,
  listRecords,
  setDefaultView,
  getViews,
} from "@/features/metadata/api";
import {
  FieldError,
  FormValues,
  ModuleField,
  RecordResponse,
  SavedView,
  TableRow,
  ValidationSchema,
} from "@/features/metadata/types";
import { errorListToMap } from "@/features/metadata/lib/validation";
import { moduleRecordPath } from "@/features/modules/paths";
import { getFormLayout, getListLayout } from "@/features/workspace/api";
import type { FormLayout } from "@/features/workspace/types";
import { useDemo } from "@/features/demo/DemoProvider";

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

export type ModuleRecordsWorkspaceProps = {
  moduleId: string;
  apiName: string;
  moduleLabel: string;
  pluralLabel: string;
  fields: ModuleField[];
  schema: ValidationSchema;
  initialViews: SavedView[];
};

/**
 * Per-module records list: saved views, create, view/edit/delete.
 * Key this component by moduleId so table state resets across modules.
 */
export default function ModuleRecordsWorkspace({
  moduleId,
  apiName,
  moduleLabel,
  pluralLabel,
  fields,
  schema,
  initialViews,
}: ModuleRecordsWorkspaceProps) {
  const router = useRouter();
  const demo = useDemo();
  const emphasizeCreate =
    demo?.mode === "running" && demo.currentStep?.step_key === "create_record";
  const [records, setRecords] = useState<RecordResponse[]>([]);
  const [views, setViews] = useState<SavedView[]>(initialViews);
  const [activeViewId, setActiveViewId] = useState<string | null>(null);
  const [adding, setAdding] = useState(false);
  const [addErrors, setAddErrors] = useState<Record<string, string>>({});
  const [reloadKey, setReloadKey] = useState(0);
  const [listColumns, setListColumns] = useState<string[]>([]);
  const [formLayout, setFormLayout] = useState<FormLayout | null>(null);
  const [liveFields, setLiveFields] = useState<ModuleField[]>(fields);
  const [liveSchema, setLiveSchema] = useState<ValidationSchema>(schema);

  useEffect(() => {
    setLiveFields(fields);
    setLiveSchema(schema);
  }, [fields, schema]);

  async function refreshMetadata() {
    try {
      const [f, s, list, form] = await Promise.all([
        getModuleFields(moduleId),
        getValidationSchema(moduleId),
        getListLayout(moduleId),
        getFormLayout(moduleId, "create"),
      ]);
      setLiveFields(f);
      setLiveSchema(s);
      setListColumns(
        list.columns
          .filter((c) => !c.system && c.field_key !== "_actions")
          .map((c) => c.field_key)
      );
      setFormLayout(form);
    } catch {
      // Keep existing metadata if refresh fails.
    }
  }

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const [list, form] = await Promise.all([
          getListLayout(moduleId),
          getFormLayout(moduleId, "create"),
        ]);
        if (!active) return;
        setListColumns(
          list.columns
            .filter((c) => !c.system && c.field_key !== "_actions")
            .map((c) => c.field_key)
        );
        setFormLayout(form);
      } catch {
        // Fall back to all fields / flat form.
      }
    })();
    return () => {
      active = false;
    };
  }, [moduleId]);

  async function openCreateForm() {
    setAddErrors({});
    setAdding(true);
    await refreshMetadata();
  }

  // Deep link / demo: /m/{apiName}?create=1 opens the add form once.
  useEffect(() => {
    if (typeof window === "undefined") return;
    const params = new URLSearchParams(window.location.search);
    if (params.get("create") !== "1") return;
    void openCreateForm();
    params.delete("create");
    const qs = params.toString();
    const path = `/m/${apiName}${qs ? `?${qs}` : ""}`;
    window.history.replaceState({}, "", path);
    // eslint-disable-next-line react-hooks/exhaustive-deps -- open once on mount when create=1
  }, [apiName]);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const result = await listRecords(moduleId, {
          page_size: 100,
          expand: true,
        });
        if (!active) return;
        setRecords(result.records);
        if (result.columns && result.columns.length > 0) {
          setListColumns(result.columns.map((c) => c.field));
        }
      } catch {
        toast.error("Failed to load records");
      }
    })();
    return () => {
      active = false;
    };
  }, [moduleId, reloadKey]);

  const rows = useMemo(
    () => records.map((rec) => recordToRow(rec, liveFields)),
    [records, liveFields]
  );

  const table = useDynamicTable({
    rows,
    fields: liveFields,
    initialColumns: listColumns.length ? listColumns : undefined,
  });

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
    if (!id || !confirm(`Delete this ${moduleLabel}?`)) return;
    try {
      await deleteRecord(moduleId, id);
      setReloadKey((k) => k + 1);
      toast.success("Record deleted");
    } catch {
      toast.error("Failed to delete record");
    }
  }

  function openRecord(row: TableRow, edit = false) {
    const id = row.id as string;
    if (!id) return;
    const base = moduleRecordPath(apiName, id);
    router.push(edit ? `${base}?edit=1` : base);
  }

  return (
    <div className="space-y-5" data-tutorial-surface="tables-records">
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
          data-tutorial-action="add-record"
          onClick={() => {
            if (adding) {
              setAdding(false);
              return;
            }
            void openCreateForm();
          }}
          className="inline-flex items-center gap-1.5 rounded-full bg-emerald-500 px-4 py-2 text-sm font-semibold text-white hover:bg-emerald-600"
        >
          {adding ? <X className="h-4 w-4" /> : <Plus className="h-4 w-4" />}
          {adding ? "Close" : `Add ${moduleLabel}`}
        </button>
      </div>

      {adding && (
        <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <DynamicForm
            key={`${moduleId}-${liveFields.map((f) => f.id).join("-")}-${formLayout?.id ?? "flat"}`}
            fields={liveFields}
            schema={liveSchema}
            formLayout={formLayout}
            submitText={`Add ${moduleLabel}`}
            sectionTitle={`New ${moduleLabel}`}
            sectionDescription={`Create a ${moduleLabel} for ${pluralLabel}. Validated by the backend.`}
            externalErrors={addErrors}
            onSubmit={handleAddRecord}
            emphasizeSubmit={emphasizeCreate}
          />
        </div>
      )}

      <DynamicTable
        fields={liveFields}
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
        onRowClick={(row) => openRecord(row)}
        onViewRow={(row) => openRecord(row)}
        onEditRow={(row) => openRecord(row, true)}
        onDeleteRow={handleDeleteRecord}
      />
    </div>
  );
}
