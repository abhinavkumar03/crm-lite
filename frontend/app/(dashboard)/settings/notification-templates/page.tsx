"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Pencil, Plus, Trash2 } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import Modal from "@/components/common/Modal";
import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";
import FormTextarea from "@/components/common/form/FormTextarea";

import {
  createTemplate,
  deleteTemplate,
  listTemplates,
  previewTemplate,
  updateTemplate,
} from "@/features/notifications/api";
import {
  NotificationChannel,
  NotificationTemplate,
  TemplateCategory,
} from "@/features/notifications/types";

const CATEGORIES: TemplateCategory[] = [
  "sales",
  "follow_up",
  "welcome",
  "proposal",
  "invoice",
  "reminder",
  "quotation",
  "marketing",
  "support",
  "custom",
];

export default function NotificationTemplatesPage() {
  const [templates, setTemplates] = useState<NotificationTemplate[]>([]);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<NotificationTemplate | null>(null);
  const [channel, setChannel] = useState<NotificationChannel>("email");
  const [name, setName] = useState("");
  const [category, setCategory] = useState<TemplateCategory>("custom");
  const [subject, setSubject] = useState("");
  const [body, setBody] = useState("");
  const [saving, setSaving] = useState(false);
  const [preview, setPreview] = useState<{ subject: string; body: string } | null>(null);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const result = await listTemplates({ page_size: 100 });
        if (active) setTemplates(result.templates);
      } catch {
        toast.error("Failed to load templates");
      }
    })();
    return () => {
      active = false;
    };
  }, [reloadKey]);

  function openCreate() {
    setEditing(null);
    setChannel("email");
    setName("");
    setCategory("custom");
    setSubject("");
    setBody("Hello {{lead.name}},\n\n");
    setPreview(null);
    setOpen(true);
  }

  function openEdit(t: NotificationTemplate) {
    setEditing(t);
    setChannel(t.channel);
    setName(t.name);
    setCategory(t.category);
    setSubject(t.subject ?? "");
    setBody(t.body);
    setPreview(null);
    setOpen(true);
  }

  async function handlePreview() {
    if (!editing) {
      toast.error("Save the template first, then preview");
      return;
    }
    try {
      const result = await previewTemplate(editing.id, {
        data: {
          "lead.name": "Alex Sample",
          "workspace.name": "Acme Workspace",
          "owner.name": "You",
        },
      });
      setPreview({ subject: result.subject, body: result.body });
    } catch {
      toast.error("Preview failed");
    }
  }

  async function handleSave() {
    if (!name.trim() || !body.trim()) {
      toast.error("Name and body are required");
      return;
    }
    try {
      setSaving(true);
      if (editing) {
        await updateTemplate(editing.id, {
          name: name.trim(),
          category,
          subject: channel === "email" ? subject : undefined,
          body,
        });
        toast.success("Template updated");
      } else {
        await createTemplate({
          channel,
          name: name.trim(),
          category,
          subject: channel === "email" ? subject : undefined,
          body,
          variables: ["lead.name", "workspace.name", "owner.name"],
        });
        toast.success("Template created");
      }
      setOpen(false);
      setReloadKey((k) => k + 1);
    } catch {
      toast.error("Failed to save template");
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(id: string) {
    if (!confirm("Delete this template?")) return;
    try {
      await deleteTemplate(id);
      toast.success("Deleted");
      setReloadKey((k) => k + 1);
    } catch {
      toast.error("Failed to delete");
    }
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Message Templates"
        description="Reusable Email and WhatsApp templates with merge fields like {{lead.name}}."
        action={
          <button
            type="button"
            onClick={openCreate}
            className="inline-flex items-center gap-2 rounded-xl bg-emerald-600 px-4 py-2 text-sm font-semibold text-white"
          >
            <Plus size={16} /> New template
          </button>
        }
      />

      <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
        <table className="min-w-full text-left text-sm">
          <thead className="border-b border-slate-100 bg-slate-50 text-xs uppercase text-slate-500">
            <tr>
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Channel</th>
              <th className="px-4 py-3">Category</th>
              <th className="px-4 py-3">Updated</th>
              <th className="px-4 py-3" />
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {templates.length === 0 && (
              <tr>
                <td colSpan={5} className="px-4 py-8 text-center text-slate-500">
                  No templates yet.
                </td>
              </tr>
            )}
            {templates.map((t) => (
              <tr key={t.id}>
                <td className="px-4 py-3 font-medium text-slate-900">{t.name}</td>
                <td className="px-4 py-3 capitalize text-slate-600">{t.channel}</td>
                <td className="px-4 py-3 capitalize text-slate-600">
                  {t.category.replace("_", " ")}
                </td>
                <td className="px-4 py-3 text-slate-500">
                  {new Date(t.updated_at).toLocaleDateString()}
                </td>
                <td className="px-4 py-3 text-right">
                  <button
                    type="button"
                    onClick={() => openEdit(t)}
                    className="mr-2 inline-flex rounded-lg p-2 text-slate-500 hover:bg-slate-100"
                  >
                    <Pencil size={16} />
                  </button>
                  <button
                    type="button"
                    onClick={() => handleDelete(t.id)}
                    className="inline-flex rounded-lg p-2 text-red-500 hover:bg-red-50"
                  >
                    <Trash2 size={16} />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Modal
        open={open}
        onClose={() => setOpen(false)}
        title={editing ? "Edit template" : "New template"}
      >
        <div className="space-y-4">
          {!editing && (
            <FormSelect
              label="Channel"
              value={channel}
              onChange={(e) =>
                setChannel(e.target.value as NotificationChannel)
              }
            >
              <option value="email">Email</option>
              <option value="whatsapp">WhatsApp</option>
            </FormSelect>
          )}
          <FormInput
            label="Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <FormSelect
            label="Category"
            value={category}
            onChange={(e) => setCategory(e.target.value as TemplateCategory)}
          >
            {CATEGORIES.map((c) => (
              <option key={c} value={c}>
                {c.replace("_", " ")}
              </option>
            ))}
          </FormSelect>
          {channel === "email" && (
            <FormInput
              label="Subject"
              value={subject}
              onChange={(e) => setSubject(e.target.value)}
            />
          )}
          <FormTextarea
            label="Body"
            value={body}
            onChange={(e) => setBody(e.target.value)}
            rows={8}
            helperText="Use {{lead.name}}, {{workspace.name}}, {{owner.name}}, {{today}}"
          />
          {preview ? (
            <div className="rounded-xl border border-slate-200 bg-slate-50 p-3 text-sm">
              <p className="font-semibold text-slate-700">Preview</p>
              {preview.subject ? (
                <p className="mt-1 text-slate-600">Subject: {preview.subject}</p>
              ) : null}
              <pre className="mt-2 whitespace-pre-wrap text-slate-800">{preview.body}</pre>
            </div>
          ) : null}
          <div className="flex gap-2">
            {editing ? (
              <button
                type="button"
                onClick={handlePreview}
                className="flex-1 rounded-xl border border-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-700"
              >
                Preview
              </button>
            ) : null}
            <button
              type="button"
              disabled={saving}
              onClick={handleSave}
              className="flex-1 rounded-xl bg-emerald-600 px-4 py-2.5 text-sm font-semibold text-white disabled:opacity-60"
            >
              {saving ? "Saving…" : "Save template"}
            </button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
