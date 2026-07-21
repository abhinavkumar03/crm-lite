"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Copy, LayoutTemplate } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import {
  cloneWorkflowTemplate,
  listWorkflowTemplates,
} from "@/features/workflows/api";
import type { WorkflowTemplate } from "@/features/workflows/types";

const CATEGORY_LABELS: Record<string, string> = {
  sales: "Sales",
  nurture: "Nurture",
  tasks: "Tasks",
  lifecycle: "Lifecycle",
};

export default function WorkflowTemplatesPage() {
  const router = useRouter();
  const [items, setItems] = useState<WorkflowTemplate[]>([]);
  const [loading, setLoading] = useState(true);
  const [category, setCategory] = useState("");
  const [cloningId, setCloningId] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const list = await listWorkflowTemplates();
        if (active) setItems(list ?? []);
      } catch {
        toast.error("Failed to load templates");
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, []);

  const categories = useMemo(() => {
    const set = new Set<string>();
    for (const t of items) {
      if (t.category) set.add(t.category);
    }
    return Array.from(set).sort();
  }, [items]);

  const filtered = useMemo(() => {
    if (!category) return items;
    return items.filter((t) => t.category === category);
  }, [items, category]);

  async function onClone(id: string) {
    try {
      setCloningId(id);
      const wf = await cloneWorkflowTemplate(id);
      toast.success("Cloned as draft — review and publish");
      router.push(`/settings/automation/workflows/${wf.id}`);
    } catch {
      toast.error("Clone failed");
    } finally {
      setCloningId(null);
    }
  }

  function previewTriggers(t: WorkflowTemplate): string {
    const triggers = (t.definition?.triggers as Array<{ type?: string }>) ?? [];
    return triggers.map((x) => x.type).filter(Boolean).join(", ") || "—";
  }

  function previewActions(t: WorkflowTemplate): string {
    const actions = (t.definition?.actions as Array<{ type?: string }>) ?? [];
    return actions.map((x) => x.type).filter(Boolean).join(", ") || "—";
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Workflow templates"
        description="Clone battle-tested starters for leads, contacts, and tasks — then customize and publish."
        action={
          <Link
            href="/settings/automation"
            className="rounded-xl border border-slate-200 px-3 py-2 text-sm"
          >
            ← Automation
          </Link>
        }
      />

      <div className="flex flex-wrap items-center gap-2">
        <button
          type="button"
          onClick={() => setCategory("")}
          className={`rounded-full px-3 py-1.5 text-xs font-semibold ${
            !category ? "bg-slate-900 text-white" : "bg-slate-100 text-slate-600"
          }`}
        >
          All ({items.length})
        </button>
        {categories.map((c) => (
          <button
            key={c}
            type="button"
            onClick={() => setCategory(c)}
            className={`rounded-full px-3 py-1.5 text-xs font-semibold ${
              category === c
                ? "bg-slate-900 text-white"
                : "bg-slate-100 text-slate-600"
            }`}
          >
            {CATEGORY_LABELS[c] ?? c}
          </button>
        ))}
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        {loading ? (
          <p className="text-sm text-slate-400">Loading gallery…</p>
        ) : filtered.length === 0 ? (
          <p className="text-sm text-slate-500">
            No templates in this category yet.
          </p>
        ) : (
          filtered.map((t) => (
            <div
              key={t.id}
              className="flex flex-col gap-3 rounded-3xl border border-slate-200 bg-white p-5 shadow-sm"
            >
              <div className="flex items-start gap-3">
                <div className="rounded-2xl bg-slate-100 p-2">
                  <LayoutTemplate className="h-4 w-4 text-slate-700" />
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <h3 className="font-semibold text-slate-900">{t.name}</h3>
                    {t.is_builtin && (
                      <span className="rounded-full bg-violet-50 px-2 py-0.5 text-[10px] font-semibold uppercase text-violet-800">
                        Built-in
                      </span>
                    )}
                    {t.category && (
                      <span className="rounded-full bg-slate-100 px-2 py-0.5 text-[10px] font-semibold uppercase text-slate-600">
                        {CATEGORY_LABELS[t.category] ?? t.category}
                      </span>
                    )}
                  </div>
                  <p className="mt-1 text-sm text-slate-500">{t.description}</p>
                  <dl className="mt-3 space-y-1 text-xs text-slate-500">
                    <div>
                      <span className="font-medium text-slate-700">Module:</span>{" "}
                      {t.module_api_name ?? "any"}
                    </div>
                    <div>
                      <span className="font-medium text-slate-700">
                        Triggers:
                      </span>{" "}
                      {previewTriggers(t)}
                    </div>
                    <div>
                      <span className="font-medium text-slate-700">
                        Actions:
                      </span>{" "}
                      {previewActions(t)}
                    </div>
                  </dl>
                </div>
              </div>
              <button
                type="button"
                disabled={cloningId === t.id}
                onClick={() => onClone(t.id)}
                className="inline-flex w-fit items-center gap-2 rounded-xl bg-slate-900 px-3 py-2 text-sm text-white disabled:opacity-60"
              >
                <Copy className="h-4 w-4" />
                {cloningId === t.id ? "Cloning…" : "Clone to workspace"}
              </button>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
