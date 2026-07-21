"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { toast } from "sonner";
import {
  Save,
  MessageCircle,
  ArrowRight,
  Workflow,
  ScrollText,
  LayoutTemplate,
} from "lucide-react";

import FormSelect from "@/components/common/form/FormSelect";
import Toggle from "@/components/common/form/Toggle";

import { getSettings, updateSettings } from "@/features/settings/api";
import { AutomationSettings } from "@/features/settings/types";
import { getWorkflowMetrics } from "@/features/workflows/api";
import type { WorkflowMetrics } from "@/features/workflows/types";

const hubLinks = [
  {
    href: "/settings/automation/workflows",
    title: "Workflows",
    description: "Create and publish automation rules",
    icon: Workflow,
  },
  {
    href: "/settings/automation/logs",
    title: "Execution logs",
    description: "Inspect runs, failures, and retries",
    icon: ScrollText,
  },
  {
    href: "/settings/automation/templates",
    title: "Templates",
    description: "Clone starter workflows",
    icon: LayoutTemplate,
  },
];

export default function AutomationSettingsPage() {
  const [automation, setAutomation] = useState<AutomationSettings | null>(null);
  const [metrics, setMetrics] = useState<WorkflowMetrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const [settings, m] = await Promise.all([
          getSettings(),
          getWorkflowMetrics().catch(() => null),
        ]);
        if (active) {
          setAutomation(settings.automation);
          setMetrics(m);
        }
      } catch {
        toast.error("Failed to load settings");
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, []);

  function patch(p: Partial<AutomationSettings>) {
    setAutomation((prev) => (prev ? { ...prev, ...p } : prev));
  }

  async function handleSave() {
    if (!automation) return;
    try {
      setSaving(true);
      const updated = await updateSettings({ automation });
      setAutomation(updated.automation);
      toast.success("Automation settings saved");
    } catch {
      toast.error("Failed to save settings");
    } finally {
      setSaving(false);
    }
  }

  if (loading || !automation) {
    return (
      <div className="rounded-3xl border border-slate-200 bg-white p-8 text-slate-400 shadow-sm">
        Loading settings...
      </div>
    );
  }

  return (
    <div className="space-y-6" data-tutorial-surface="automation">
      <section className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">
            Automation center
          </h2>
          <p className="text-sm text-slate-500">
            Define metadata-driven workflows across any module, then tune
            notification preferences.
          </p>
        </div>

        {metrics && (
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
            {[
              ["Active", metrics.active_workflows],
              ["Draft", metrics.draft_workflows],
              ["Executed today", metrics.executed_today],
              ["Failed today", metrics.failed_today],
            ].map(([label, value]) => (
              <div
                key={String(label)}
                className="rounded-2xl border border-slate-100 bg-slate-50 px-4 py-3"
              >
                <p className="text-xs uppercase tracking-wide text-slate-500">
                  {label}
                </p>
                <p className="mt-1 text-2xl font-semibold text-slate-900">
                  {value}
                </p>
              </div>
            ))}
          </div>
        )}

        <div className="grid gap-3 md:grid-cols-3">
          {hubLinks.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className="group flex flex-col gap-2 rounded-2xl border border-slate-200 p-4 transition hover:border-slate-300 hover:bg-slate-50"
            >
              <item.icon className="h-5 w-5 text-slate-700" />
              <div className="flex items-center justify-between gap-2">
                <span className="font-medium text-slate-900">{item.title}</span>
                <ArrowRight className="h-4 w-4 text-slate-400 transition group-hover:translate-x-0.5" />
              </div>
              <p className="text-sm text-slate-500">{item.description}</p>
            </Link>
          ))}
        </div>
      </section>

      <section className="space-y-5 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">
            Notification preferences
          </h2>
          <p className="text-sm text-slate-500">
            Master switches for the outbound notification pipeline. Delivery
            credentials live under Communication Providers.
          </p>
        </div>

        <div className="rounded-2xl border border-slate-200 p-4">
          <Toggle
            label="Enable notifications"
            description="Master switch for outbound WhatsApp & email."
            checked={automation.notifications_enabled}
            onChange={(v) => patch({ notifications_enabled: v })}
          />
        </div>

        <FormSelect
          label="Default channel"
          value={automation.default_channel}
          onChange={(e) =>
            patch({
              default_channel: e.target.value as AutomationSettings["default_channel"],
            })
          }
          options={[
            { value: "email", label: "Email" },
            { value: "whatsapp", label: "WhatsApp" },
          ]}
        />

        <div className="rounded-2xl border border-slate-200 p-4">
          <Toggle
            label="Daily digest"
            description="Summarize queued notifications once per day."
            checked={automation.daily_digest}
            onChange={(v) => patch({ daily_digest: v })}
          />
        </div>

        <div className="flex flex-wrap items-center gap-3">
          <button
            type="button"
            onClick={handleSave}
            disabled={saving}
            className="inline-flex items-center gap-2 rounded-xl bg-slate-900 px-4 py-2 text-sm font-medium text-white hover:bg-slate-800 disabled:opacity-60"
          >
            <Save className="h-4 w-4" />
            {saving ? "Saving…" : "Save preferences"}
          </button>
          <Link
            href="/notifications"
            className="inline-flex items-center gap-2 text-sm font-medium text-slate-700 hover:text-slate-900"
          >
            <MessageCircle className="h-4 w-4" />
            Open notification center
          </Link>
        </div>
      </section>
    </div>
  );
}
