"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Plus, Power, PowerOff, Archive } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import {
  archiveWorkflow,
  disableWorkflow,
  listWorkflows,
  publishWorkflow,
} from "@/features/workflows/api";
import type { WorkflowSummary } from "@/features/workflows/types";

const statusClass: Record<string, string> = {
  active: "bg-emerald-50 text-emerald-700",
  draft: "bg-slate-100 text-slate-700",
  disabled: "bg-amber-50 text-amber-800",
  archived: "bg-rose-50 text-rose-700",
};

export default function WorkflowsListPage() {
  const router = useRouter();
  const [items, setItems] = useState<WorkflowSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [status, setStatus] = useState("");

  async function load() {
    try {
      setLoading(true);
      const result = await listWorkflows({
        page_size: 50,
        status: status || undefined,
      });
      setItems(result.items ?? []);
    } catch {
      toast.error("Failed to load workflows");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [status]);

  async function onPublish(id: string) {
    try {
      await publishWorkflow(id);
      toast.success("Workflow published");
      load();
    } catch {
      toast.error("Publish failed");
    }
  }

  async function onDisable(id: string) {
    try {
      await disableWorkflow(id);
      toast.success("Workflow disabled");
      load();
    } catch {
      toast.error("Disable failed");
    }
  }

  async function onArchive(id: string) {
    try {
      await archiveWorkflow(id);
      toast.success("Workflow archived");
      load();
    } catch {
      toast.error("Archive failed");
    }
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Workflows"
        description="Automate record events across any CRM module."
        action={
          <Link
            href="/settings/automation/workflows/new"
            className="inline-flex items-center gap-2 rounded-xl bg-slate-900 px-4 py-2 text-sm font-medium text-white hover:bg-slate-800"
          >
            <Plus className="h-4 w-4" />
            New workflow
          </Link>
        }
      />

      <div className="flex items-center gap-3">
        <select
          value={status}
          onChange={(e) => setStatus(e.target.value)}
          className="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm"
        >
          <option value="">All statuses</option>
          <option value="draft">Draft</option>
          <option value="active">Active</option>
          <option value="disabled">Disabled</option>
          <option value="archived">Archived</option>
        </select>
        <Link
          href="/settings/automation"
          className="text-sm text-slate-500 hover:text-slate-800"
        >
          ← Automation center
        </Link>
      </div>

      <div className="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-sm">
        {loading ? (
          <p className="p-6 text-sm text-slate-400">Loading…</p>
        ) : items.length === 0 ? (
          <p className="p-6 text-sm text-slate-500">
            No workflows yet. Create one or clone a template.
          </p>
        ) : (
          <table className="min-w-full text-left text-sm">
            <thead className="border-b border-slate-100 bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
              <tr>
                <th className="px-4 py-3 font-medium">Name</th>
                <th className="px-4 py-3 font-medium">Module</th>
                <th className="px-4 py-3 font-medium">Triggers</th>
                <th className="px-4 py-3 font-medium">Status</th>
                <th className="px-4 py-3 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {items.map((w) => (
                <tr key={w.id} className="border-b border-slate-50 last:border-0">
                  <td className="px-4 py-3">
                    <button
                      type="button"
                      onClick={() =>
                        router.push(`/settings/automation/workflows/${w.id}`)
                      }
                      className="font-medium text-slate-900 hover:underline"
                    >
                      {w.name}
                    </button>
                    <p className="text-xs text-slate-500">{w.description}</p>
                  </td>
                  <td className="px-4 py-3 text-slate-600">
                    {w.module_api_name ?? "—"}
                  </td>
                  <td className="px-4 py-3 text-slate-600">
                    {(w.trigger_types ?? []).join(", ") || "—"}
                  </td>
                  <td className="px-4 py-3">
                    <span
                      className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusClass[w.status] ?? "bg-slate-100"}`}
                    >
                      {w.status}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex flex-wrap gap-2">
                      {w.status !== "active" && (
                        <button
                          type="button"
                          onClick={() => onPublish(w.id)}
                          className="inline-flex items-center gap-1 rounded-lg border border-slate-200 px-2 py-1 text-xs hover:bg-slate-50"
                        >
                          <Power className="h-3 w-3" /> Publish
                        </button>
                      )}
                      {w.status === "active" && (
                        <button
                          type="button"
                          onClick={() => onDisable(w.id)}
                          className="inline-flex items-center gap-1 rounded-lg border border-slate-200 px-2 py-1 text-xs hover:bg-slate-50"
                        >
                          <PowerOff className="h-3 w-3" /> Disable
                        </button>
                      )}
                      {w.status !== "archived" && (
                        <button
                          type="button"
                          onClick={() => onArchive(w.id)}
                          className="inline-flex items-center gap-1 rounded-lg border border-slate-200 px-2 py-1 text-xs hover:bg-slate-50"
                        >
                          <Archive className="h-3 w-3" /> Archive
                        </button>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
