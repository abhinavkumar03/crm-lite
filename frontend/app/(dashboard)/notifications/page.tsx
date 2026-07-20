"use client";

import { useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import {
  Ban,
  Mail,
  MessageCircle,
  RefreshCw,
  RotateCcw,
  Send,
} from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";
import FormTextarea from "@/components/common/form/FormTextarea";

import {
  cancelNotification,
  composeNotification,
  getNotification,
  getNotificationMetrics,
  listNotifications,
  listTemplates,
  retryNotification,
} from "@/features/notifications/api";
import {
  ComposeMode,
  Notification,
  NotificationChannel,
  NotificationMetrics,
  NotificationStatus,
  NotificationTemplate,
} from "@/features/notifications/types";

const STATUS_STYLES: Record<string, string> = {
  draft: "bg-slate-100 text-slate-700",
  scheduled: "bg-indigo-100 text-indigo-700",
  queued: "bg-amber-100 text-amber-700",
  processing: "bg-sky-100 text-sky-700",
  sent: "bg-emerald-100 text-emerald-700",
  delivered: "bg-emerald-100 text-emerald-800",
  opened: "bg-teal-100 text-teal-800",
  read: "bg-teal-100 text-teal-800",
  failed: "bg-red-100 text-red-700",
  retrying: "bg-orange-100 text-orange-700",
  cancelled: "bg-slate-200 text-slate-600",
};

function StatusBadge({ status }: { status: string }) {
  return (
    <span
      className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold capitalize ${STATUS_STYLES[status] ?? "bg-slate-100 text-slate-700"}`}
    >
      {status}
    </span>
  );
}

export default function NotificationsPage() {
  const [mode, setMode] = useState<ComposeMode>("send");
  const [channel, setChannel] = useState<NotificationChannel>("email");
  const [to, setTo] = useState("");
  const [cc, setCc] = useState("");
  const [bcc, setBcc] = useState("");
  const [subject, setSubject] = useState("");
  const [templateId, setTemplateId] = useState("");
  const [body, setBody] = useState(
    "Hello {{lead.name}},\n\nThank you for showing interest in {{workspace.name}}."
  );
  const [scheduledAt, setScheduledAt] = useState("");
  const [sending, setSending] = useState(false);

  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [templates, setTemplates] = useState<NotificationTemplate[]>([]);
  const [metrics, setMetrics] = useState<NotificationMetrics | null>(null);
  const [statusFilter, setStatusFilter] = useState<"" | NotificationStatus>("");
  const [channelFilter, setChannelFilter] = useState<"" | NotificationChannel>("");
  const [q, setQ] = useState("");
  const [selected, setSelected] = useState<Notification | null>(null);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const [list, tpls, m] = await Promise.all([
          listNotifications({
            page_size: 50,
            status: statusFilter || undefined,
            channel: channelFilter || undefined,
            q: q.trim() || undefined,
          }),
          listTemplates({ page_size: 100 }),
          getNotificationMetrics(),
        ]);
        if (!active) return;
        setNotifications(list.notifications);
        setTemplates(tpls.templates);
        setMetrics(m);
      } catch {
        toast.error("Failed to load notification center");
      }
    })();
    return () => {
      active = false;
    };
  }, [statusFilter, channelFilter, q, reloadKey]);

  const channelTemplates = useMemo(
    () => templates.filter((t) => t.channel === channel && t.is_active),
    [templates, channel]
  );

  function applyTemplate(id: string) {
    setTemplateId(id);
    const t = templates.find((x) => x.id === id);
    if (!t) return;
    if (t.subject) setSubject(t.subject);
    setBody(t.body);
  }

  async function handleCompose() {
    if (!to.trim()) {
      toast.error("Recipient is required");
      return;
    }
    if (!body.trim()) {
      toast.error("Message body is required");
      return;
    }
    if (channel === "email" && mode !== "draft" && !subject.trim()) {
      toast.error("Subject is required for email");
      return;
    }
    if (mode === "schedule" && !scheduledAt) {
      toast.error("Pick a schedule time");
      return;
    }

    try {
      setSending(true);
      await composeNotification({
        mode,
        channel,
        to: to.trim(),
        cc: cc
          .split(",")
          .map((s) => s.trim())
          .filter(Boolean),
        bcc: bcc
          .split(",")
          .map((s) => s.trim())
          .filter(Boolean),
        subject: channel === "email" ? subject : undefined,
        body,
        template_id: templateId || undefined,
        scheduled_at:
          mode === "schedule" ? new Date(scheduledAt).toISOString() : undefined,
      });
      toast.success(
        mode === "draft"
          ? "Draft saved"
          : mode === "schedule"
            ? "Message scheduled"
            : "Notification queued"
      );
      setReloadKey((k) => k + 1);
    } catch {
      toast.error("Failed to compose notification");
    } finally {
      setSending(false);
    }
  }

  async function openDetail(id: string) {
    try {
      const n = await getNotification(id);
      setSelected(n);
    } catch {
      toast.error("Failed to load notification");
    }
  }

  async function onRetry(id: string) {
    try {
      await retryNotification(id);
      toast.success("Retry queued");
      setReloadKey((k) => k + 1);
      setSelected(null);
    } catch {
      toast.error("Retry failed");
    }
  }

  async function onCancel(id: string) {
    try {
      await cancelNotification(id);
      toast.success("Cancelled");
      setReloadKey((k) => k + 1);
      setSelected(null);
    } catch {
      toast.error("Cancel failed");
    }
  }

  return (
    <div className="space-y-8" data-tour="notifications-page">
      <PageHeader
        title="Notification Center"
        description="Compose Email & WhatsApp messages, track delivery, and manage templates."
        action={
          <button
            type="button"
            onClick={() => setReloadKey((k) => k + 1)}
            className="inline-flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
          >
            <RefreshCw size={16} />
            Refresh
          </button>
        }
      />

      {metrics && (
        <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {[
            { label: "Emails today", value: metrics.emails_sent_today, icon: Mail },
            {
              label: "WhatsApp today",
              value: metrics.whatsapp_sent_today,
              icon: MessageCircle,
            },
            { label: "Failed / retrying", value: metrics.failed_count, icon: RotateCcw },
            { label: "Scheduled", value: metrics.scheduled_count, icon: Send },
            {
              label: "Delivery rate",
              value: `${Math.round(metrics.delivery_rate ?? 0)}%`,
              icon: Send,
            },
            {
              label: "Open rate",
              value: `${Math.round(metrics.open_rate ?? 0)}%`,
              icon: Mail,
            },
            {
              label: "Read rate",
              value: `${Math.round(metrics.read_rate ?? 0)}%`,
              icon: MessageCircle,
            },
            { label: "Drafts", value: metrics.draft_count ?? 0, icon: Ban },
          ].map((m) => (
            <div
              key={m.label}
              className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm"
            >
              <div className="flex items-center justify-between">
                <p className="text-sm text-slate-500">{m.label}</p>
                <m.icon size={18} className="text-slate-400" />
              </div>
              <p className="mt-2 text-2xl font-semibold text-slate-900">{m.value}</p>
            </div>
          ))}
        </section>
      )}

      <div className="grid gap-8 xl:grid-cols-[420px_1fr]">
        <section className="space-y-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <h2 className="text-lg font-semibold text-slate-900">Compose</h2>

          <FormSelect
            label="Mode"
            value={mode}
            onChange={(e) => setMode(e.target.value as ComposeMode)}
          >
            <option value="send">Send now</option>
            <option value="schedule">Schedule later</option>
            <option value="draft">Save draft</option>
          </FormSelect>

          <FormSelect
            label="Channel"
            value={channel}
            onChange={(e) => setChannel(e.target.value as NotificationChannel)}
          >
            <option value="email">Email</option>
            <option value="whatsapp">WhatsApp</option>
          </FormSelect>

          <FormSelect
            label="Template"
            value={templateId}
            onChange={(e) => applyTemplate(e.target.value)}
          >
            <option value="">None</option>
            {channelTemplates.map((t) => (
              <option key={t.id} value={t.id}>
                {t.name} ({t.category})
              </option>
            ))}
          </FormSelect>

          <FormInput
            label={channel === "email" ? "To" : "Recipient"}
            value={to}
            onChange={(e) => setTo(e.target.value)}
            placeholder={channel === "email" ? "name@company.com" : "+15551234567"}
          />

          {channel === "email" && (
            <>
              <FormInput
                label="CC"
                value={cc}
                onChange={(e) => setCc(e.target.value)}
                placeholder="comma-separated"
              />
              <FormInput
                label="BCC"
                value={bcc}
                onChange={(e) => setBcc(e.target.value)}
                placeholder="comma-separated"
              />
              <FormInput
                label="Subject"
                value={subject}
                onChange={(e) => setSubject(e.target.value)}
              />
            </>
          )}

          <FormTextarea
            label="Message"
            value={body}
            onChange={(e) => setBody(e.target.value)}
            rows={8}
          />

          {mode === "schedule" && (
            <FormInput
              label="Schedule at"
              type="datetime-local"
              value={scheduledAt}
              onChange={(e) => setScheduledAt(e.target.value)}
            />
          )}

          <button
            type="button"
            data-tutorial-action="send-notification"
            disabled={sending}
            onClick={handleCompose}
            className="inline-flex w-full items-center justify-center gap-2 rounded-xl bg-emerald-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-emerald-700 disabled:opacity-60"
          >
            <Send size={16} />
            {sending
              ? "Working…"
              : mode === "draft"
                ? "Save draft"
                : mode === "schedule"
                  ? "Schedule"
                  : "Send"}
          </button>
        </section>

        <section className="space-y-4">
          <div className="grid gap-3 sm:grid-cols-3">
            <FormInput
              label="Search"
              value={q}
              onChange={(e) => setQ(e.target.value)}
              placeholder="Recipient, subject…"
            />
            <FormSelect
              label="Status"
              value={statusFilter}
              onChange={(e) =>
                setStatusFilter(e.target.value as "" | NotificationStatus)
              }
            >
              <option value="">All statuses</option>
              {(
                [
                  "draft",
                  "scheduled",
                  "queued",
                  "sent",
                  "delivered",
                  "opened",
                  "read",
                  "failed",
                  "retrying",
                  "cancelled",
                ] as NotificationStatus[]
              ).map((s) => (
                <option key={s} value={s}>
                  {s}
                </option>
              ))}
            </FormSelect>
            <FormSelect
              label="Channel"
              value={channelFilter}
              onChange={(e) =>
                setChannelFilter(e.target.value as "" | NotificationChannel)
              }
            >
              <option value="">All channels</option>
              <option value="email">Email</option>
              <option value="whatsapp">WhatsApp</option>
            </FormSelect>
          </div>

          <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
            <ul className="divide-y divide-slate-100">
              {notifications.length === 0 && (
                <li className="p-8 text-center text-sm text-slate-500">
                  No notifications yet. Compose one to get started.
                </li>
              )}
              {notifications.map((n) => (
                <li key={n.id}>
                  <button
                    type="button"
                    onClick={() => openDetail(n.id)}
                    className="flex w-full items-start justify-between gap-4 px-4 py-3 text-left hover:bg-slate-50"
                  >
                    <div className="min-w-0">
                      <div className="flex items-center gap-2">
                        {n.channel === "email" ? (
                          <Mail size={14} className="text-slate-400" />
                        ) : (
                          <MessageCircle size={14} className="text-slate-400" />
                        )}
                        <p className="truncate text-sm font-medium text-slate-900">
                          {n.subject || n.body.slice(0, 60)}
                        </p>
                      </div>
                      <p className="mt-1 truncate text-xs text-slate-500">
                        {n.recipient} · {new Date(n.created_at).toLocaleString()}
                      </p>
                    </div>
                    <StatusBadge status={n.status} />
                  </button>
                </li>
              ))}
            </ul>
          </div>
        </section>
      </div>

      {selected && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 p-4">
          <div className="max-h-[90vh] w-full max-w-lg overflow-y-auto rounded-2xl bg-white p-6 shadow-xl">
            <div className="flex items-start justify-between gap-4">
              <div>
                <h3 className="text-lg font-semibold text-slate-900">
                  {selected.subject || "Message"}
                </h3>
                <p className="mt-1 text-sm text-slate-500">{selected.recipient}</p>
              </div>
              <StatusBadge status={selected.status} />
            </div>
            <pre className="mt-4 whitespace-pre-wrap rounded-xl bg-slate-50 p-3 text-sm text-slate-700">
              {selected.body}
            </pre>
            {(selected.error || selected.last_error) && (
              <p className="mt-3 text-sm text-red-600">
                {selected.last_error || selected.error}
              </p>
            )}
            {selected.events && selected.events.length > 0 && (
              <ul className="mt-4 space-y-2 border-t border-slate-100 pt-4">
                {selected.events.map((e) => (
                  <li
                    key={e.id}
                    className="flex justify-between text-xs text-slate-500"
                  >
                    <span className="capitalize">{e.event}</span>
                    <span>{new Date(e.created_at).toLocaleString()}</span>
                  </li>
                ))}
              </ul>
            )}
            <div className="mt-6 flex flex-wrap gap-2">
              {(selected.status === "failed" ||
                selected.status === "retrying") && (
                <button
                  type="button"
                  onClick={() => onRetry(selected.id)}
                  className="inline-flex items-center gap-2 rounded-xl bg-amber-600 px-3 py-2 text-sm font-medium text-white"
                >
                  <RotateCcw size={14} /> Retry
                </button>
              )}
              {selected.status === "scheduled" && (
                <button
                  type="button"
                  onClick={() => onCancel(selected.id)}
                  className="inline-flex items-center gap-2 rounded-xl bg-slate-700 px-3 py-2 text-sm font-medium text-white"
                >
                  <Ban size={14} /> Cancel
                </button>
              )}
              <button
                type="button"
                onClick={() => setSelected(null)}
                className="ml-auto rounded-xl border border-slate-200 px-3 py-2 text-sm font-medium text-slate-700"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
