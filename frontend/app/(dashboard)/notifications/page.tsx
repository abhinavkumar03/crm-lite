"use client";

import { useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { RefreshCw, Send } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";
import FormTextarea from "@/components/common/form/FormTextarea";

import { listNotifications, sendNotification } from "@/features/notifications/api";
import {
  Notification,
  NotificationChannel,
  NotificationStatus,
} from "@/features/notifications/types";

const STATUS_STYLES: Record<NotificationStatus, string> = {
  queued: "bg-amber-100 text-amber-700",
  sent: "bg-emerald-100 text-emerald-700",
  failed: "bg-red-100 text-red-700",
};

function StatusBadge({ status }: { status: NotificationStatus }) {
  return (
    <span
      className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold capitalize ${STATUS_STYLES[status]}`}
    >
      {status}
    </span>
  );
}

export default function NotificationsPage() {
  const [channel, setChannel] = useState<NotificationChannel>("whatsapp");
  const [to, setTo] = useState("");
  const [subject, setSubject] = useState("");
  const [template, setTemplate] = useState("");
  const [body, setBody] = useState("Hi {{name}}, thanks for your interest!");
  const [dataText, setDataText] = useState('{\n  "name": "Dana"\n}');
  const [sending, setSending] = useState(false);

  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [statusFilter, setStatusFilter] = useState<"" | NotificationStatus>("");
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const result = await listNotifications({
          page_size: 50,
          status: statusFilter || undefined,
        });
        if (active) setNotifications(result.notifications);
      } catch {
        toast.error("Failed to load notifications");
      }
    })();
    return () => {
      active = false;
    };
  }, [statusFilter, reloadKey]);

  const parsedData = useMemo<{ ok: boolean; value: Record<string, unknown> }>(() => {
    if (!dataText.trim()) return { ok: true, value: {} };
    try {
      return { ok: true, value: JSON.parse(dataText) };
    } catch {
      return { ok: false, value: {} };
    }
  }, [dataText]);

  async function handleSend() {
    if (!to.trim()) {
      toast.error("Recipient is required");
      return;
    }
    if (!body.trim()) {
      toast.error("Message body is required");
      return;
    }
    if (!parsedData.ok) {
      toast.error("Template data is not valid JSON");
      return;
    }

    try {
      setSending(true);
      await sendNotification({
        channel,
        to: to.trim(),
        subject: channel === "email" ? subject : undefined,
        template: template || undefined,
        body,
        data: parsedData.value,
      });
      toast.success("Notification queued");
      setReloadKey((k) => k + 1);
    } catch {
      toast.error("Failed to queue notification");
    } finally {
      setSending(false);
    }
  }

  return (
    <div className="space-y-8">
      <PageHeader
        badge="Automation"
        title="WhatsApp & Email Notifications"
        description="Compose a message and dispatch it through the provider-agnostic notification pipeline. Delivery runs asynchronously in the worker; the log below reflects each message's lifecycle."
      />

      <div className="grid gap-8 lg:grid-cols-[1fr_1.2fr]">
        <div className="space-y-5 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <h3 className="text-sm font-semibold text-slate-900">Compose</h3>

          <FormSelect
            label="Channel"
            value={channel}
            onChange={(e) => setChannel(e.target.value as NotificationChannel)}
          >
            <option value="whatsapp">WhatsApp</option>
            <option value="email">Email</option>
          </FormSelect>

          <FormInput
            label={channel === "whatsapp" ? "Recipient phone" : "Recipient email"}
            placeholder={channel === "whatsapp" ? "+15551234567" : "person@example.com"}
            value={to}
            requiredMark
            onChange={(e) => setTo(e.target.value)}
          />

          {channel === "email" && (
            <FormInput
              label="Subject"
              placeholder="Welcome, {{name}}"
              value={subject}
              onChange={(e) => setSubject(e.target.value)}
            />
          )}

          <FormInput
            label="Template label (optional)"
            placeholder="lead_welcome"
            value={template}
            onChange={(e) => setTemplate(e.target.value)}
          />

          <FormTextarea
            label="Body"
            helperText="Use {{placeholder}} tokens; they are rendered from the data below."
            rows={4}
            value={body}
            requiredMark
            onChange={(e) => setBody(e.target.value)}
          />

          <FormTextarea
            label="Template data (JSON)"
            helperText={parsedData.ok ? "Values injected into the body/subject." : "Invalid JSON"}
            rows={4}
            value={dataText}
            onChange={(e) => setDataText(e.target.value)}
          />

          <button
            type="button"
            onClick={handleSend}
            disabled={sending}
            className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:opacity-50"
          >
            <Send className="h-4 w-4" />
            {sending ? "Queuing..." : "Send notification"}
          </button>
        </div>

        <div className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <div className="flex items-center justify-between gap-3">
            <h3 className="text-sm font-semibold text-slate-900">Delivery log</h3>
            <div className="flex items-center gap-2">
              <select
                value={statusFilter}
                onChange={(e) =>
                  setStatusFilter(e.target.value as "" | NotificationStatus)
                }
                className="rounded-lg border border-slate-200 bg-white px-2 py-1 text-xs text-slate-700 focus:border-emerald-400 focus:outline-none"
              >
                <option value="">All statuses</option>
                <option value="queued">Queued</option>
                <option value="sent">Sent</option>
                <option value="failed">Failed</option>
              </select>
              <button
                type="button"
                onClick={() => setReloadKey((k) => k + 1)}
                className="rounded-lg border border-slate-200 p-1.5 text-slate-600 hover:bg-slate-100"
                aria-label="Refresh"
              >
                <RefreshCw className="h-4 w-4" />
              </button>
            </div>
          </div>

          <div className="overflow-hidden rounded-2xl border border-slate-200">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50 text-left text-slate-600">
                  <th className="px-3 py-2 font-semibold">Channel</th>
                  <th className="px-3 py-2 font-semibold">Recipient</th>
                  <th className="px-3 py-2 font-semibold">Status</th>
                  <th className="px-3 py-2 font-semibold">Provider</th>
                </tr>
              </thead>
              <tbody>
                {notifications.length === 0 ? (
                  <tr>
                    <td colSpan={4} className="px-3 py-8 text-center text-slate-400">
                      No notifications yet.
                    </td>
                  </tr>
                ) : (
                  notifications.map((n) => (
                    <tr
                      key={n.id}
                      className="border-b border-slate-100 last:border-0 hover:bg-slate-50/60"
                    >
                      <td className="px-3 py-2 capitalize text-slate-700">
                        {n.channel}
                      </td>
                      <td className="px-3 py-2 text-slate-700">{n.recipient}</td>
                      <td className="px-3 py-2">
                        <StatusBadge status={n.status} />
                        {n.status === "failed" && n.error && (
                          <p className="mt-1 max-w-xs truncate text-xs text-red-400">
                            {n.error}
                          </p>
                        )}
                      </td>
                      <td className="px-3 py-2 text-slate-500">
                        {n.provider ?? "—"}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          <p className="text-xs text-slate-500">
            Statuses update as the worker processes the queue. Run the worker with
            <code className="mx-1 rounded bg-slate-100 px-1">make run-worker</code>
            to see messages move from queued to sent.
          </p>
        </div>
      </div>
    </div>
  );
}
