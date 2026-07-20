"use client";

import { Suspense, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { toast } from "sonner";
import { Loader2, RotateCcw } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import FormSelect from "@/components/common/form/FormSelect";
import { getModules } from "@/features/metadata/api";
import type { ModuleSummary } from "@/features/metadata/types";
import ListingColumnsEditor from "@/features/settings/components/ListingColumnsEditor";
import {
  getListLayout,
  reorderListColumns,
  resetListLayout,
  toggleListColumn,
} from "@/features/workspace/api";
import type { ListColumn } from "@/features/workspace/types";

function TablesSettingsInner() {
  const searchParams = useSearchParams();
  const moduleFromQuery = searchParams.get("module") ?? "";

  const [modules, setModules] = useState<ModuleSummary[]>([]);
  const [moduleId, setModuleId] = useState(moduleFromQuery);
  const [columns, setColumns] = useState<ListColumn[]>([]);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    (async () => {
      try {
        const all = await getModules();
        const enabled = all.filter((m) => m.is_enabled);
        setModules(enabled);
        if (!moduleId && enabled[0]) {
          setModuleId(enabled[0].id);
        }
      } catch {
        toast.error("Failed to load modules");
      }
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (!moduleId) return;
    let active = true;
    (async () => {
      try {
        setLoading(true);
        const layout = await getListLayout(moduleId, { includeHidden: true });
        if (active) setColumns(layout.columns);
      } catch {
        toast.error("Failed to load list layout");
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [moduleId]);

  async function handleReorder(next: ListColumn[]) {
    const payload = next
      .filter((c) => !c.system && c.field_key !== "_actions")
      .map((c, i) => ({ field_key: c.field_key, order: i + 1 }));
    const previous = columns;
    setColumns(next);
    try {
      setSaving(true);
      const updated = await reorderListColumns(moduleId, payload);
      setColumns(updated.columns);
      toast.success("Column order updated");
    } catch {
      setColumns(previous);
      toast.error("Failed to reorder columns");
    } finally {
      setSaving(false);
    }
  }

  async function handleToggle(col: ListColumn) {
    if (col.locked || col.system) {
      toast.error("Locked columns cannot be hidden");
      return;
    }
    try {
      setSaving(true);
      const updated = await toggleListColumn(moduleId, col.field_key, !col.visible);
      setColumns(updated.columns);
      toast.success(col.visible ? "Column hidden" : "Column shown");
    } catch {
      toast.error("Failed to toggle column");
    } finally {
      setSaving(false);
    }
  }

  async function handleReset() {
    if (!moduleId) return;
    try {
      setSaving(true);
      const updated = await resetListLayout(moduleId);
      setColumns(updated.columns);
      toast.success("Listing columns reset to defaults");
    } catch {
      toast.error("Failed to reset listing columns");
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Listing columns"
        description="Configure which fields appear in module tables, and in what order. Independent from form field sections."
        action={
          moduleId ? (
            <button
              type="button"
              disabled={saving || loading}
              onClick={handleReset}
              className="inline-flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm font-semibold text-slate-700 hover:bg-slate-50 disabled:opacity-50"
            >
              <RotateCcw className="h-4 w-4" />
              Reset to default
            </button>
          ) : undefined
        }
      />

      <div className="max-w-sm">
        <FormSelect
          label="Module"
          value={moduleId}
          onChange={(e) => setModuleId(e.target.value)}
        >
          <option value="">Select a module</option>
          {modules.map((m) => (
            <option key={m.id} value={m.id}>
              {m.plural_label}
            </option>
          ))}
        </FormSelect>
      </div>

      {!moduleId ? (
        <p className="text-sm text-slate-500">
          Select a module to edit its listing columns.
        </p>
      ) : loading ? (
        <div className="flex items-center gap-2 py-12 text-sm text-slate-400">
          <Loader2 className="h-4 w-4 animate-spin" />
          Loading listing columns…
        </div>
      ) : (
        <ListingColumnsEditor
          columns={columns}
          saving={saving}
          onReorder={handleReorder}
          onToggle={handleToggle}
        />
      )}
    </div>
  );
}

/** Settings → Listing Columns: org default list layout editor. */
export default function TablesSettingsPage() {
  return (
    <Suspense
      fallback={
        <div className="py-16 text-center text-sm text-slate-400">
          Loading…
        </div>
      }
    >
      <TablesSettingsInner />
    </Suspense>
  );
}
