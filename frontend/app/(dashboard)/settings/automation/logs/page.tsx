"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { toast } from "sonner";
import { RefreshCw, RotateCcw } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import {
  getExecution,
  listExecutions,
  retryExecution,
} from "@/features/workflows/api";
import type { ExecutionDetail, ExecutionSummary } from "@/features/workflows/types";

const statusClass: Record<string, string> = {
  succeeded: "text-emerald-700",
  failed: "text-rose-600",
  partial: "text-amber-700",
  running: "text-sky-700",
  queued: "text-slate-500",
};

export default function WorkflowLogsPage() {
  const [items, setItems] = useState<ExecutionSummary[]>([]);
  const [status, setStatus] = useState("");
  const [selected, setSelected] = useState<ExecutionDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [retryingId, setRetryingId] = useState<string | null>(null);

  async function load() {
    try {
      setLoading(true);
      const result = await listExecutions({
        page_size: 50,
        status: status || undefined,
      });
      setItems(result.items ?? []);
    } catch {
      toast.error("Failed to load execution logs");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [status]);

  async function openDetail(id: string) {
    try {
      setSelected(await getExecution(id));
    } catch {
      toast.error("Failed to load execution");
    }
  }

  async function onRetry(id: string, e?: React.MouseEvent) {
    e?.stopPropagation();
    try {
      setRetryingId(id);
      await retryExecution(id);
      toast.success("Retry queued — refresh shortly for the new run");
      await load();
      if (selected?.id === id) {
        setSelected(await getExecution(id));
      }
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { message?: string } } })?.response?.data
          ?.message ?? "Retry failed";
      toast.error(msg);
    } finally {
      setRetryingId(null);
    }
  }

  const canRetry = (s: string) => s === "failed" || s === "partial";

  return (
    <div className="space-y-6">
      <PageHeader
        title="Execution logs"
        description="Append-only history of workflow runs. Failed and partial runs can be retried."
        action={
          <Link
            href="/settings/automation"
            className="rounded-xl border border-slate-200 px-3 py-2 text-sm"
          >
            ← Automation
          </Link>
        }
      />

      <div className="flex flex-wrap gap-2">
        {[
          { value: "", label: "All" },
          { value: "failed", label: "Failed" },
          { value: "partial", label: "Partial" },
          { value: "succeeded", label: "Succeeded" },
          { value: "running", label: "Running" },
        ].map((f) => (
          <button
            key={f.value || "all"}
            type="button"
            onClick={() => setStatus(f.value)}
            className={`rounded-full px-3 py-1.5 text-xs font-semibold ${
              status === f.value
                ? "bg-slate-900 text-white"
                : "bg-slate-100 text-slate-600"
            }`}
          >
            {f.label}
          </button>
        ))}
        <button
          type="button"
          onClick={load}
          className="inline-flex items-center gap-1 rounded-full border border-slate-200 px-3 py-1.5 text-xs"
        >
          <RefreshCw className="h-3 w-3" /> Refresh
        </button>
      </div>

      <div className="grid gap-4 lg:grid-cols-2">
        <div className="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-sm">
          {loading ? (
            <p className="p-6 text-sm text-slate-400">Loading…</p>
          ) : items.length === 0 ? (
            <p className="p-6 text-sm text-slate-500">
              {status === "failed"
                ? "No failed runs — nice."
                : "No executions yet."}
            </p>
          ) : (
            <ul className="divide-y divide-slate-100">
              {items.map((e) => (
                <li key={e.id}>
                  <div className="flex items-start gap-2 px-4 py-3 hover:bg-slate-50">
                    <button
                      type="button"
                      onClick={() => openDetail(e.id)}
                      className="min-w-0 flex-1 text-left"
                    >
                      <div className="flex items-center justify-between gap-2">
                        <span className="font-medium text-slate-900">
                          {e.workflow_name || e.workflow_id}
                        </span>
                        <span
                          className={`text-xs font-semibold uppercase ${statusClass[e.status] ?? "text-slate-500"}`}
                        >
                          {e.status}
                        </span>
                      </div>
                      <p className="text-xs text-slate-500">
                        {e.trigger_type} · {e.source} ·{" "}
                        {new Date(e.created_at).toLocaleString()}
                      </p>
                      {e.error_summary && (
                        <p className="mt-1 line-clamp-2 text-xs text-rose-600">
                          {e.error_summary}
                        </p>
                      )}
                    </button>
                    {canRetry(e.status) && (
                      <button
                        type="button"
                        disabled={retryingId === e.id}
                        onClick={(ev) => onRetry(e.id, ev)}
                        className="inline-flex shrink-0 items-center gap-1 rounded-lg border border-slate-200 px-2 py-1 text-xs hover:bg-white disabled:opacity-60"
                        title="Retry failed run"
                      >
                        <RotateCcw className="h-3 w-3" />
                        {retryingId === e.id ? "…" : "Retry"}
                      </button>
                    )}
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          {!selected ? (
            <p className="text-sm text-slate-500">
              Select a run to inspect steps. Use Retry on failed/partial runs to
              re-queue the same workflow against the same record.
            </p>
          ) : (
            <div className="space-y-4">
              <div className="flex items-start justify-between gap-2">
                <div>
                  <h3 className="font-semibold text-slate-900">
                    {selected.workflow_name}
                  </h3>
                  <p className="text-xs text-slate-500">
                    <span className={statusClass[selected.status]}>
                      {selected.status}
                    </span>
                    {selected.duration_ms != null
                      ? ` · ${selected.duration_ms}ms`
                      : ""}
                  </p>
                  {selected.record_id && (
                    <p className="mt-1 text-xs text-slate-500">
                      Record: {selected.record_id.slice(0, 8)}…
                    </p>
                  )}
                </div>
                {canRetry(selected.status) && (
                  <button
                    type="button"
                    disabled={retryingId === selected.id}
                    onClick={() => onRetry(selected.id)}
                    className="inline-flex items-center gap-1 rounded-lg border border-slate-200 px-2 py-1 text-xs"
                  >
                    <RotateCcw className="h-3 w-3" />
                    {retryingId === selected.id ? "Queuing…" : "Retry run"}
                  </button>
                )}
              </div>
              {selected.error_summary && (
                <div className="rounded-xl border border-rose-100 bg-rose-50 px-3 py-2 text-xs text-rose-700">
                  {selected.error_summary}
                </div>
              )}
              <ol className="space-y-2">
                {(selected.steps ?? []).map((s) => (
                  <li
                    key={s.id}
                    className="rounded-xl border border-slate-100 px-3 py-2 text-sm"
                  >
                    <div className="flex justify-between">
                      <span className="font-medium">{s.action_type}</span>
                      <span
                        className={`text-xs uppercase ${statusClass[s.status] ?? "text-slate-500"}`}
                      >
                        {s.status}
                      </span>
                    </div>
                    {s.error && (
                      <p className="mt-1 text-xs text-rose-600">{s.error}</p>
                    )}
                  </li>
                ))}
              </ol>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
