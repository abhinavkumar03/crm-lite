"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Mail, MessageCircle, Plus, Send } from "lucide-react";

import Modal from "@/components/common/Modal";
import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";
import FormTextarea from "@/components/common/form/FormTextarea";

import {
  composeNotification,
  listNotifications,
  listTemplates,
} from "@/features/notifications/api";
import {
  Notification,
  NotificationChannel,
  NotificationTemplate,
} from "@/features/notifications/types";

type Props = {
  moduleId: string;
  recordId: string;
  channel: NotificationChannel;
  defaultTo?: string;
  recordLabel?: string;
};

const STATUS_STYLES: Record<string, string> = {
  draft: "bg-slate-100 text-slate-700",
  scheduled: "bg-indigo-100 text-indigo-700",
  queued: "bg-amber-100 text-amber-700",
  sent: "bg-emerald-100 text-emerald-700",
  delivered: "bg-emerald-100 text-emerald-800",
  opened: "bg-teal-100 text-teal-800",
  read: "bg-teal-100 text-teal-800",
  failed: "bg-red-100 text-red-700",
  retrying: "bg-orange-100 text-orange-700",
  cancelled: "bg-slate-200 text-slate-600",
};

export default function RecordCommunicationsPanel({
  moduleId,
  recordId,
  channel,
  defaultTo = "",
  recordLabel = "",
}: Props) {
  const [items, setItems] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [to, setTo] = useState(defaultTo);
  const [subject, setSubject] = useState("");
  const [body, setBody] = useState(
    channel === "email"
      ? `Hello {{lead.name}},\n\nFollowing up from our CRM.`
      : `Hi {{lead.name}}, can we schedule a quick call?`
  );
  const [templateId, setTemplateId] = useState("");
  const [templates, setTemplates] = useState<NotificationTemplate[]>([]);
  const [sending, setSending] = useState(false);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    setTo(defaultTo);
  }, [defaultTo]);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        setLoading(true);
        const [result, tpls] = await Promise.all([
          listNotifications({
            page_size: 50,
            channel,
            module_id: moduleId,
            entity_id: recordId,
            entity_type: "RECORD",
          }),
          listTemplates({ page_size: 100, channel }),
        ]);
        if (active) {
          setItems(result.notifications);
          setTemplates(tpls.templates.filter((t) => t.is_active));
        }
      } catch {
        toast.error("Failed to load messages");
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [moduleId, recordId, channel, reloadKey]);

  async function handleSend() {
    if (!to.trim() || !body.trim()) {
      toast.error("Recipient and message are required");
      return;
    }
    if (channel === "email" && !subject.trim()) {
      toast.error("Subject is required");
      return;
    }
    try {
      setSending(true);
      await composeNotification({
        mode: "send",
        channel,
        to: to.trim(),
        subject: channel === "email" ? subject : undefined,
        body,
        template_id: templateId || undefined,
        entity_type: "RECORD",
        entity_id: recordId,
        module_id: moduleId,
        data: recordLabel ? { "lead.name": recordLabel, name: recordLabel } : {},
      });
      toast.success("Message queued");
      setOpen(false);
      setReloadKey((k) => k + 1);
    } catch {
      toast.error("Failed to send");
    } finally {
      setSending(false);
    }
  }

  const title = channel === "email" ? "Emails" : "WhatsApp";
  const Icon = channel === "email" ? Mail : MessageCircle;

  return (
    <div className="space-y-4" data-tutorial-surface={`record-${channel}`}>
      <div className="flex items-center justify-between">
        <h2 className="flex items-center gap-2 text-lg font-semibold text-slate-900">
          <Icon size={18} /> {title}
        </h2>
        <button
          type="button"
          onClick={() => setOpen(true)}
          className="inline-flex items-center gap-2 rounded-xl bg-emerald-600 px-3 py-2 text-sm font-medium text-white hover:bg-emerald-700"
          data-tutorial-action={`compose-${channel}`}
        >
          <Plus size={16} /> Compose
        </button>
      </div>

      {loading ? (
        <p className="text-sm text-slate-500">Loading…</p>
      ) : items.length === 0 ? (
        <p className="rounded-2xl border border-dashed border-slate-200 p-8 text-center text-sm text-slate-500">
          No {title.toLowerCase()} yet for this record.
        </p>
      ) : (
        <ul className="space-y-3">
          {items.map((n) => (
            <li
              key={n.id}
              className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm"
            >
              <div className="flex items-start justify-between gap-3">
                <div className="min-w-0">
                  <p className="truncate text-sm font-medium text-slate-900">
                    {n.subject || n.body.slice(0, 80)}
                  </p>
                  <p className="mt-1 text-xs text-slate-500">
                    {n.recipient} · {new Date(n.created_at).toLocaleString()}
                  </p>
                </div>
                <span
                  className={`inline-flex shrink-0 rounded-full px-2.5 py-0.5 text-xs font-semibold capitalize ${STATUS_STYLES[n.status] ?? "bg-slate-100 text-slate-700"}`}
                >
                  {n.status}
                </span>
              </div>
              <p className="mt-3 whitespace-pre-wrap text-sm text-slate-600">
                {n.body}
              </p>
            </li>
          ))}
        </ul>
      )}

      <Modal
        open={open}
        onClose={() => setOpen(false)}
        title={`Compose ${title}`}
      >
        <div className="space-y-4">
          <FormSelect label="Channel" value={channel} disabled>
            <option value={channel}>{title}</option>
          </FormSelect>
          <FormInput
            label="To"
            value={to}
            onChange={(e) => setTo(e.target.value)}
          />
          {templates.length > 0 ? (
            <FormSelect
              label="Template"
              value={templateId}
              onChange={(e) => {
                const id = e.target.value;
                setTemplateId(id);
                const t = templates.find((x) => x.id === id);
                if (t) {
                  if (t.subject) setSubject(t.subject);
                  setBody(t.body);
                }
              }}
            >
              <option value="">Custom message</option>
              {templates.map((t) => (
                <option key={t.id} value={t.id}>
                  {t.name}
                </option>
              ))}
            </FormSelect>
          ) : null}
          {channel === "email" && (
            <FormInput
              label="Subject"
              value={subject}
              onChange={(e) => setSubject(e.target.value)}
            />
          )}
          <FormTextarea
            label="Message"
            value={body}
            onChange={(e) => setBody(e.target.value)}
            rows={6}
          />
          <button
            type="button"
            disabled={sending}
            onClick={handleSend}
            className="inline-flex w-full items-center justify-center gap-2 rounded-xl bg-emerald-600 px-4 py-2.5 text-sm font-semibold text-white disabled:opacity-60"
          >
            <Send size={16} />
            {sending ? "Sending…" : "Send"}
          </button>
        </div>
      </Modal>
    </div>
  );
}
