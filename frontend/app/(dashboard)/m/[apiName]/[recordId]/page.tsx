"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import { toast } from "sonner";
import {
  ArrowLeft,
  Building2,
  Loader2,
  Paperclip,
  Pencil,
  Plus,
  Trash2,
  X,
} from "lucide-react";

import DynamicForm from "@/features/metadata/components/DynamicForm";
import {
  deleteRecord,
  getModuleFields,
  getModules,
  getRecord,
  getValidationSchema,
  listRecords,
  updateRecord,
} from "@/features/metadata/api";
import {
  FormValues,
  ModuleField,
  ModuleSummary,
  RecordResponse,
  ValidationSchema,
} from "@/features/metadata/types";
import { errorListToMap } from "@/features/metadata/lib/validation";

import {
  createRecordNote,
  deleteRecordAttachment,
  deleteRecordNote,
  getDetailLayout,
  getFormLayout,
  listRecordActivities,
  listRecordAttachments,
  listRecordNotes,
  listRelatedDescriptors,
  uploadRecordAttachment,
} from "@/features/workspace/api";
import {
  DetailLayout,
  FormLayout,
  RelatedDescriptor,
  WorkspaceActivity,
  WorkspaceAttachment,
  WorkspaceNote,
} from "@/features/workspace/types";
import { useDemo } from "@/features/demo/DemoProvider";
import RecordCommunicationsPanel from "@/features/notifications/components/RecordCommunicationsPanel";

type Tab =
  | "overview"
  | "notes"
  | "attachments"
  | "emails"
  | "whatsapp"
  | "timeline"
  | "related";

const VALID_TABS: Tab[] = [
  "overview",
  "notes",
  "attachments",
  "emails",
  "whatsapp",
  "timeline",
  "related",
];

function recordTitle(rec: RecordResponse, fields: ModuleField[]): string {
  const preferred = fields.find(
    (f) =>
      f.is_required &&
      (f.field_type === "text" || f.api_name === "name" || f.api_name === "title")
  );
  const nameField =
    fields.find((f) => f.api_name === "name") ||
    fields.find((f) => f.api_name === "title") ||
    preferred ||
    fields[0];
  if (!nameField) return rec.id.slice(0, 8);
  const raw = rec.data[nameField.api_name];
  if (raw == null || raw === "") return rec.id.slice(0, 8);
  return String(raw);
}

function formatValue(field: ModuleField, rec: RecordResponse): string {
  const rel = rec.relations?.[field.api_name];
  if (rel?.label) return rel.label;
  const v = rec.data[field.api_name];
  if (v == null || v === "") return "—";
  if (typeof v === "boolean") return v ? "Yes" : "No";
  if (Array.isArray(v)) return v.join(", ");
  return String(v);
}

export default function ModuleRecordPage() {
  const params = useParams<{ apiName: string; recordId: string }>();
  const apiName = decodeURIComponent(params.apiName ?? "");
  const recordId = params.recordId;
  const router = useRouter();
  const searchParams = useSearchParams();
  const demo = useDemo();
  const tutorialNote =
    demo?.mode === "running" && demo.currentStep?.step_key === "add_note";
  const tutorialTimeline =
    demo?.mode === "running" && demo.currentStep?.step_key === "timeline";

  const [moduleId, setModuleId] = useState("");
  const [module, setModule] = useState<ModuleSummary | null>(null);
  const [fields, setFields] = useState<ModuleField[]>([]);
  const [schema, setSchema] = useState<ValidationSchema | null>(null);
  const [layout, setLayout] = useState<DetailLayout | null>(null);
  const [formLayout, setFormLayout] = useState<FormLayout | null>(null);
  const [record, setRecord] = useState<RecordResponse | null>(null);
  const [tab, setTab] = useState<Tab>("overview");
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const [editErrors, setEditErrors] = useState<Record<string, string>>({});
  const [saving, setSaving] = useState(false);

  const [notes, setNotes] = useState<WorkspaceNote[]>([]);
  const [noteBody, setNoteBody] = useState("");
  const [attachments, setAttachments] = useState<WorkspaceAttachment[]>([]);
  const [activities, setActivities] = useState<WorkspaceActivity[]>([]);
  const [related, setRelated] = useState<RelatedDescriptor[]>([]);
  const [relatedRows, setRelatedRows] = useState<
    Record<string, RecordResponse[]>
  >({});

  // Deep-link / demo navigation: ?tab=notes|timeline|… and ?edit=1
  useEffect(() => {
    const raw = searchParams.get("tab");
    if (raw && VALID_TABS.includes(raw as Tab)) {
      setTab(raw as Tab);
    }
    if (searchParams.get("edit") === "1") {
      setEditing(true);
    }
  }, [searchParams]);

  useEffect(() => {
    if (tutorialNote) {
      setTab("notes");
      setNoteBody((prev) =>
        prev.trim()
          ? prev
          : "Followed up with the prospect — demo walkthrough note."
      );
    } else if (tutorialTimeline) {
      setTab("timeline");
    }
  }, [tutorialNote, tutorialTimeline]);

  const reloadRecord = useCallback(async () => {
    if (!moduleId) return;
    const rec = await getRecord(moduleId, recordId, true);
    setRecord(rec);
  }, [moduleId, recordId]);

  const reloadSide = useCallback(async () => {
    if (!moduleId) return;
    const [n, a, t, r] = await Promise.all([
      listRecordNotes(moduleId, recordId),
      listRecordAttachments(moduleId, recordId),
      listRecordActivities(moduleId, recordId),
      listRelatedDescriptors(moduleId),
    ]);
    setNotes(n);
    setAttachments(a);
    setActivities(t);
    setRelated(r);
  }, [moduleId, recordId]);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        setLoading(true);
        const mods = await getModules();
        const mod =
          mods.find((m) => m.api_name.toLowerCase() === apiName.toLowerCase()) ??
          null;
        if (!mod) {
          toast.error("Module not found");
          router.replace("/dashboard");
          return;
        }
        if (!active) return;
        setModule(mod);
        setModuleId(mod.id);
        const [f, s, lay, formLay] = await Promise.all([
          getModuleFields(mod.id),
          getValidationSchema(mod.id),
          getDetailLayout(mod.id),
          getFormLayout(mod.id, "edit"),
        ]);
        if (!active) return;
        setFields(f);
        setSchema(s);
        setLayout(lay);
        setFormLayout(formLay);
        const rec = await getRecord(mod.id, recordId, true);
        if (!active) return;
        setRecord(rec);
        const [n, a, t, r] = await Promise.all([
          listRecordNotes(mod.id, recordId),
          listRecordAttachments(mod.id, recordId),
          listRecordActivities(mod.id, recordId),
          listRelatedDescriptors(mod.id),
        ]);
        if (!active) return;
        setNotes(n);
        setAttachments(a);
        setActivities(t);
        setRelated(r);
      } catch {
        toast.error("Failed to load record workspace");
        router.replace(apiName ? `/m/${encodeURIComponent(apiName)}` : "/dashboard");
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [apiName, recordId, router]);

  useEffect(() => {
    if (tab !== "related" || related.length === 0) return;
    let active = true;
    (async () => {
      const entries = await Promise.all(
        related.map(async (d) => {
          try {
            const params: Record<string, string | number | boolean> = {
              page_size: 25,
              expand: true,
            };
            params[`filter.${d.lookup_field_api_name}`] = recordId;
            const result = await listRecords(
              d.child_module_id,
              params as Parameters<typeof listRecords>[1]
            );
            return [d.child_module_id, result.records] as [string, RecordResponse[]];
          } catch {
            return [d.child_module_id, [] as RecordResponse[]] as [
              string,
              RecordResponse[],
            ];
          }
        })
      );
      if (!active) return;
      const map: Record<string, RecordResponse[]> = {};
      for (const [id, rows] of entries) map[id] = rows;
      setRelatedRows(map);
    })();
    return () => {
      active = false;
    };
  }, [tab, related, recordId]);

  const title = useMemo(
    () => (record ? recordTitle(record, fields) : "Record"),
    [record, fields]
  );

  const sections = layout?.config?.sections?.length
    ? layout.config.sections
    : [
        {
          key: "general",
          label: "General Information",
          fields: fields.filter((f) => f.is_visible).map((f) => f.api_name),
        },
      ];

  const fieldMap = useMemo(
    () => new Map(fields.map((f) => [f.api_name, f])),
    [fields]
  );

  async function handleSave(values: FormValues) {
    if (!record || !moduleId) return;
    try {
      setSaving(true);
      setEditErrors({});
      await updateRecord(moduleId, recordId, values);
      await reloadRecord();
      await reloadSide();
      setEditing(false);
      toast.success("Record updated");
    } catch (err: unknown) {
      const errors = (err as { response?: { data?: { errors?: unknown } } })
        ?.response?.data?.errors;
      if (Array.isArray(errors)) {
        setEditErrors(errorListToMap(errors as { field: string; message: string }[]));
      }
      toast.error("Failed to update record");
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete() {
    if (!moduleId || !confirm("Delete this record?")) return;
    try {
      await deleteRecord(moduleId, recordId);
      toast.success("Record deleted");
      router.replace(`/m/${encodeURIComponent(apiName)}`);
    } catch {
      toast.error("Failed to delete");
    }
  }

  async function handleAddNote() {
    if (!moduleId || !noteBody.trim()) return;
    try {
      await createRecordNote(moduleId, recordId, noteBody.trim());
      setNoteBody("");
      await reloadSide();
      toast.success("Note added");
      if (tutorialNote) {
        await demo?.validate({ silent: true, stepKey: "add_note" });
      }
    } catch {
      toast.error("Failed to add note");
    }
  }

  async function handleUpload(file: File | null) {
    if (!moduleId || !file) return;
    try {
      await uploadRecordAttachment(moduleId, recordId, file);
      await reloadSide();
      toast.success("Attachment uploaded");
    } catch {
      toast.error("Upload failed");
    }
  }

  if (loading || !record || !schema) {
    return (
      <div className="flex items-center justify-center py-20 text-sm text-slate-400">
        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
        Loading workspace…
      </div>
    );
  }

  const tabs: { id: Tab; label: string }[] = [
    { id: "overview", label: "Overview" },
    { id: "notes", label: "Notes" },
    { id: "attachments", label: "Attachments" },
    { id: "emails", label: "Emails" },
    { id: "whatsapp", label: "WhatsApp" },
    { id: "timeline", label: "Timeline" },
    { id: "related", label: "Related" },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-3 text-sm text-slate-500">
        <Link
          href={`/m/${encodeURIComponent(apiName)}`}
          className="inline-flex items-center gap-1 font-medium text-emerald-600 hover:text-emerald-700"
        >
          <ArrowLeft className="h-4 w-4" />
          {module?.plural_label ?? "Records"}
        </Link>
        <span>/</span>
        <span className="font-semibold text-slate-800">{title}</span>
      </div>

      <header className="flex flex-wrap items-start justify-between gap-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div className="flex items-start gap-4">
          <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-emerald-50 text-emerald-600">
            <Building2 className="h-7 w-7" />
          </div>
          <div>
            <h1 className="text-2xl font-semibold text-slate-900">{title}</h1>
            <p className="mt-1 text-sm text-slate-500">
              {module?.singular_label ?? "Record"} · Visibility{" "}
              <span className="capitalize">
                {record.visibility || "organization"}
              </span>
            </p>
            <p className="mt-1 text-xs text-slate-400">
              Created {new Date(record.created_at).toLocaleString()} · Updated{" "}
              {new Date(record.updated_at).toLocaleString()}
            </p>
          </div>
        </div>
        <div className="flex flex-wrap gap-2">
          <button
            type="button"
            onClick={() => setEditing((v) => !v)}
            className="inline-flex items-center gap-2 rounded-full border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-700 hover:bg-slate-50"
          >
            {editing ? <X className="h-4 w-4" /> : <Pencil className="h-4 w-4" />}
            {editing ? "Cancel" : "Edit"}
          </button>
          <button
            type="button"
            onClick={handleDelete}
            className="inline-flex items-center gap-2 rounded-full border border-red-200 px-4 py-2 text-sm font-semibold text-red-600 hover:bg-red-50"
          >
            <Trash2 className="h-4 w-4" />
            Delete
          </button>
        </div>
      </header>

      <div className="flex flex-wrap gap-2 border-b border-slate-200 pb-2">
        {tabs.map((t) => (
          <button
            key={t.id}
            type="button"
            data-tutorial-action={
              t.id === "notes"
                ? "open-notes-tab"
                : t.id === "timeline"
                  ? "open-timeline-tab"
                  : undefined
            }
            onClick={() => setTab(t.id)}
            className={`rounded-full px-4 py-2 text-sm font-semibold transition ${
              tab === t.id
                ? "bg-emerald-500 text-white"
                : "text-slate-600 hover:bg-slate-100"
            }`}
          >
            {t.label}
          </button>
        ))}
      </div>

      {tab === "overview" && (
        <div className="space-y-5" data-tutorial-surface="record-overview">
          {editing ? (
            <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
              <DynamicForm
                fields={fields}
                schema={schema}
                formLayout={formLayout}
                initialValues={record.data as FormValues}
                externalErrors={editErrors}
                submitText={saving ? "Saving…" : "Save changes"}
                sectionTitle="Edit record"
                onSubmit={handleSave}
                onCancel={() => setEditing(false)}
              />
            </div>
          ) : (
            sections.map((section) => (
              <section
                key={section.key}
                className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm"
              >
                <h2 className="mb-4 text-lg font-semibold text-slate-900">
                  {section.label}
                </h2>
                <dl className="grid gap-4 sm:grid-cols-2">
                  {section.fields.map((apiName) => {
                    if (
                      ["owner_id", "assigned_to", "visibility", "created_at", "updated_at"].includes(
                        apiName
                      )
                    ) {
                      const sysVal =
                        apiName === "owner_id"
                          ? record.owner_id ?? "—"
                          : apiName === "assigned_to"
                            ? record.assigned_to ?? "—"
                            : apiName === "visibility"
                              ? record.visibility ?? "—"
                              : apiName === "created_at"
                                ? new Date(record.created_at).toLocaleString()
                                : new Date(record.updated_at).toLocaleString();
                      return (
                        <div key={apiName}>
                          <dt className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                            {apiName.replace(/_/g, " ")}
                          </dt>
                          <dd className="mt-1 text-sm text-slate-800">{sysVal}</dd>
                        </div>
                      );
                    }
                    const field = fieldMap.get(apiName);
                    if (!field) return null;
                    return (
                      <div key={apiName}>
                        <dt className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                          {field.label}
                        </dt>
                        <dd className="mt-1 text-sm text-slate-800">
                          {formatValue(field, record)}
                        </dd>
                      </div>
                    );
                  })}
                </dl>
              </section>
            ))
          )}
        </div>
      )}

      {tab === "notes" && (
        <div
          className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm"
          data-tutorial-surface="add-note"
        >
          {tutorialNote && (
            <p className="rounded-2xl bg-emerald-50 px-4 py-3 text-sm text-emerald-800">
              A follow-up note is pre-filled. Click <strong>Add</strong> to
              continue the walkthrough.
            </p>
          )}
          <div className="flex gap-2">
            <textarea
              value={noteBody}
              onChange={(e) => setNoteBody(e.target.value)}
              rows={3}
              placeholder="Write a note…"
              readOnly={tutorialNote}
              className="w-full rounded-2xl border border-slate-200 px-4 py-3 text-sm outline-none focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100 read-only:bg-slate-50"
            />
            <button
              type="button"
              data-tutorial-action="add-note"
              onClick={handleAddNote}
              className="inline-flex h-fit items-center gap-2 rounded-full bg-emerald-500 px-4 py-2.5 text-sm font-semibold text-white"
            >
              <Plus className="h-4 w-4" />
              Add
            </button>
          </div>
          <ul className="divide-y divide-slate-100">
            {notes.length === 0 ? (
              <li className="py-6 text-sm text-slate-400">No notes yet.</li>
            ) : (
              notes.map((n) => (
                <li key={n.id} className="flex justify-between gap-3 py-4">
                  <div>
                    {n.title ? (
                      <p className="text-sm font-semibold text-slate-800">{n.title}</p>
                    ) : null}
                    <p className="whitespace-pre-wrap text-sm text-slate-700">{n.body}</p>
                    <p className="mt-1 text-xs text-slate-400">
                      {n.author_name} · {new Date(n.created_at).toLocaleString()}
                    </p>
                  </div>
                  <button
                    type="button"
                    onClick={async () => {
                      await deleteRecordNote(moduleId, recordId, n.id);
                      await reloadSide();
                    }}
                    className="text-slate-300 hover:text-red-500"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </li>
              ))
            )}
          </ul>
        </div>
      )}

      {tab === "attachments" && (
        <div className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <label className="inline-flex cursor-pointer items-center gap-2 rounded-full border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-700 hover:bg-slate-50">
            <Paperclip className="h-4 w-4" />
            Upload file
            <input
              type="file"
              className="hidden"
              onChange={(e) => handleUpload(e.target.files?.[0] ?? null)}
            />
          </label>
          <ul className="divide-y divide-slate-100">
            {attachments.length === 0 ? (
              <li className="py-6 text-sm text-slate-400">No attachments.</li>
            ) : (
              attachments.map((a) => (
                <li key={a.id} className="flex items-center justify-between gap-3 py-3">
                  <a
                    href={a.file_url}
                    target="_blank"
                    rel="noreferrer"
                    className="text-sm font-medium text-emerald-600 hover:underline"
                  >
                    {a.file_name}
                  </a>
                  <div className="flex items-center gap-3">
                    <span className="text-xs text-slate-400">
                      {a.uploader_name} · {new Date(a.created_at).toLocaleString()}
                    </span>
                    <button
                      type="button"
                      onClick={async () => {
                        await deleteRecordAttachment(moduleId, recordId, a.id);
                        await reloadSide();
                      }}
                      className="text-slate-300 hover:text-red-500"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </li>
              ))
            )}
          </ul>
        </div>
      )}

      {tab === "emails" && module && (
        <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <RecordCommunicationsPanel
            moduleId={module.id}
            recordId={record.id}
            channel="email"
            defaultTo={String(record.data?.email ?? "")}
            recordLabel={title}
          />
        </div>
      )}

      {tab === "whatsapp" && module && (
        <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <RecordCommunicationsPanel
            moduleId={module.id}
            recordId={record.id}
            channel="whatsapp"
            defaultTo={String(
              record.data?.phone ?? record.data?.mobile ?? ""
            )}
            recordLabel={title}
          />
        </div>
      )}

      {tab === "timeline" && (
        <div
          className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm"
          data-tutorial-surface="record-timeline"
        >
          <ul className="space-y-4">
            {activities.length === 0 ? (
              <li className="text-sm text-slate-400">No activity yet.</li>
            ) : (
              activities.map((a) => (
                <li key={a.id} className="flex gap-3 border-l-2 border-emerald-200 pl-4">
                  <div>
                    <p className="text-sm font-semibold text-slate-800">{a.description}</p>
                    <p className="text-xs text-slate-400">
                      {a.action} · {a.actor_name} ·{" "}
                      {new Date(a.created_at).toLocaleString()}
                    </p>
                  </div>
                </li>
              ))
            )}
          </ul>
        </div>
      )}

      {tab === "related" && (
        <div className="space-y-5">
          {related.length === 0 ? (
            <div className="rounded-3xl border border-slate-200 bg-white p-6 text-sm text-slate-400 shadow-sm">
              No related modules (lookup fields pointing here).
            </div>
          ) : (
            related.map((d) => {
              const rows = relatedRows[d.child_module_id] ?? [];
              return (
                <section
                  key={d.child_module_id}
                  className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm"
                >
                  <div className="mb-3 flex items-center justify-between">
                    <h2 className="text-lg font-semibold text-slate-900">
                      {d.child_module_name}
                    </h2>
                    <Link
                      href={`/m/${encodeURIComponent(d.child_api_name)}?create=1`}
                      className="text-sm font-semibold text-emerald-600"
                    >
                      Add record
                    </Link>
                  </div>
                  {rows.length === 0 ? (
                    <p className="text-sm text-slate-400">No related records.</p>
                  ) : (
                    <ul className="divide-y divide-slate-100">
                      {rows.map((r) => (
                        <li key={r.id}>
                          <Link
                            href={`/m/${encodeURIComponent(d.child_api_name)}/${r.id}`}
                            className="block py-2 text-sm font-medium text-slate-700 hover:text-emerald-600"
                          >
                            {String(
                              r.data.name ??
                                r.data.title ??
                                r.id.slice(0, 8)
                            )}
                          </Link>
                        </li>
                      ))}
                    </ul>
                  )}
                </section>
              );
            })
          )}
        </div>
      )}
    </div>
  );
}
