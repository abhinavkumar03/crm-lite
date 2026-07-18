"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { toast } from "sonner";
import { Save, MessageCircle, ArrowRight } from "lucide-react";

import FormSelect from "@/components/common/form/FormSelect";
import Toggle from "@/components/common/form/Toggle";

import { getSettings, updateSettings } from "@/features/settings/api";
import { AutomationSettings } from "@/features/settings/types";

export default function AutomationSettingsPage() {
  const [automation, setAutomation] = useState<AutomationSettings | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const data = await getSettings();
        if (active) setAutomation(data.automation);
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
    <div className="space-y-6">
      <section className="space-y-5 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">Automation</h2>
          <p className="text-sm text-slate-500">
            Behavioural preferences for the notification pipeline. Provider
            credentials are configured via environment variables.
          </p>
        </div>

        <div className="rounded-2xl border border-slate-200 p-4">
          <Toggle
            label="Enable notifications"
            description="Master switch for outbound WhatsApp & email automations."
            checked={automation.notifications_enabled}
            onChange={(v) => patch({ notifications_enabled: v })}
          />
        </div>

        <FormSelect
          label="Default channel"
          value={automation.default_channel}
          onChange={(e) =>
            patch({ default_channel: e.target.value as "whatsapp" | "email" })
          }
        >
          <option value="whatsapp">WhatsApp</option>
          <option value="email">Email</option>
        </FormSelect>

        <div className="rounded-2xl border border-slate-200 p-4">
          <Toggle
            label="Daily digest"
            description="Send a once-a-day summary instead of per-event messages."
            checked={automation.daily_digest}
            onChange={(v) => patch({ daily_digest: v })}
          />
        </div>

        <div className="flex justify-end">
          <button
            type="button"
            onClick={handleSave}
            disabled={saving}
            className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:opacity-50"
          >
            <Save className="h-4 w-4" />
            {saving ? "Saving..." : "Save changes"}
          </button>
        </div>
      </section>

      <Link
        href="/notifications"
        className="flex items-center justify-between gap-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm transition hover:border-emerald-300 hover:shadow-md"
      >
        <div className="flex items-center gap-4">
          <span className="flex h-11 w-11 items-center justify-center rounded-2xl bg-emerald-50 text-emerald-600">
            <MessageCircle className="h-5 w-5" />
          </span>
          <div>
            <p className="font-semibold text-slate-800">Compose a notification</p>
            <p className="text-sm text-slate-500">
              Send a WhatsApp/email message and watch delivery in the log.
            </p>
          </div>
        </div>
        <ArrowRight className="h-5 w-5 text-slate-400" />
      </Link>
    </div>
  );
}
