"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Play, Workflow } from "lucide-react";

import { listWorkflows, runWorkflow } from "@/features/workflows/api";
import type { WorkflowSummary } from "@/features/workflows/types";

type Props = {
  moduleId: string;
  recordId: string;
};

export default function RecordWorkflowsPanel({ moduleId, recordId }: Props) {
  const [items, setItems] = useState<WorkflowSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [runningId, setRunningId] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const result = await listWorkflows({
          status: "active",
          module_id: moduleId,
          page_size: 50,
        });
        if (!active) return;
        const list = result.items ?? [];
        // Prefer workflows that declare a manual trigger; fall back to all active.
        const manual = list.filter((w) =>
          (w.trigger_types ?? []).includes("manual")
        );
        setItems(manual.length > 0 ? manual : list);
      } catch {
        if (active) setItems([]);
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [moduleId]);

  async function onRun(workflowId: string, name: string) {
    try {
      setRunningId(workflowId);
      await runWorkflow(workflowId, recordId, moduleId);
      toast.success(`Queued “${name}”`);
    } catch {
      toast.error("Failed to queue workflow");
    } finally {
      setRunningId(null);
    }
  }

  if (loading) {
    return (
      <section className="rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
        <p className="text-sm text-slate-400">Loading workflows…</p>
      </section>
    );
  }

  if (items.length === 0) {
    return null;
  }

  return (
    <section className="rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
      <div className="mb-3 flex items-center gap-2">
        <Workflow className="h-4 w-4 text-slate-700" />
        <h2 className="text-sm font-semibold text-slate-900">Run workflow</h2>
      </div>
      <p className="mb-3 text-xs text-slate-500">
        Manually execute an active automation against this record. Runs are
        asynchronous — check Automation → Logs for results.
      </p>
      <ul className="space-y-2">
        {items.map((w) => (
          <li
            key={w.id}
            className="flex items-center justify-between gap-3 rounded-2xl border border-slate-100 px-3 py-2"
          >
            <div className="min-w-0">
              <p className="truncate text-sm font-medium text-slate-900">
                {w.name}
              </p>
              <p className="truncate text-xs text-slate-500">
                {(w.trigger_types ?? []).join(", ") || "active"}
              </p>
            </div>
            <button
              type="button"
              disabled={runningId === w.id}
              onClick={() => onRun(w.id, w.name)}
              className="inline-flex shrink-0 items-center gap-1.5 rounded-full bg-slate-900 px-3 py-1.5 text-xs font-semibold text-white hover:bg-slate-800 disabled:opacity-60"
            >
              <Play className="h-3 w-3" />
              {runningId === w.id ? "Queuing…" : "Run"}
            </button>
          </li>
        ))}
      </ul>
    </section>
  );
}
