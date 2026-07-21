"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { toast } from "sonner";
import { Plus, Trash2, Save, Power, RotateCcw } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";
import FormTextarea from "@/components/common/form/FormTextarea";

import {
  createWorkflow,
  getBuilderMetadata,
  getWorkflow,
  listWorkflowVersions,
  publishWorkflow,
  rollbackWorkflow,
  updateWorkflow,
} from "@/features/workflows/api";
import type {
  ActionInput,
  BuilderMetadata,
  ConditionInput,
  TriggerInput,
  VersionSummary,
} from "@/features/workflows/types";

const emptyCondition = (): ConditionInput => ({
  node_type: "group",
  logic: "and",
  children: [],
});

export default function WorkflowEditorPage() {
  const params = useParams();
  const router = useRouter();
  const id = typeof params.id === "string" ? params.id : "";
  const isNew = id === "new" || !id;

  const [meta, setMeta] = useState<BuilderMetadata | null>(null);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [moduleId, setModuleId] = useState("");
  const [onError, setOnError] = useState<"stop" | "continue">("stop");
  const [triggers, setTriggers] = useState<TriggerInput[]>([
    { type: "record_created", config: {} },
  ]);
  const [conditions, setConditions] = useState<ConditionInput>(emptyCondition());
  const [actions, setActions] = useState<ActionInput[]>([]);
  const [versions, setVersions] = useState<VersionSummary[]>([]);
  const [status, setStatus] = useState<string>("draft");
  const [publishedVersionId, setPublishedVersionId] = useState<string | null>(
    null
  );
  const [draftVersionId, setDraftVersionId] = useState<string | null>(null);
  const [changelog, setChangelog] = useState("");
  const [showPublishDialog, setShowPublishDialog] = useState(false);
  const [version, setVersion] = useState(0);
  const [saving, setSaving] = useState(false);
  const [loading, setLoading] = useState(!isNew);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const m = await getBuilderMetadata();
        if (!active) return;
        setMeta(m);
        if (!isNew) {
          const wf = await getWorkflow(id);
          if (!active) return;
          setName(wf.name);
          setDescription(wf.description ?? "");
          setModuleId(wf.module_id ?? "");
          setOnError(wf.on_action_error ?? "stop");
          setStatus(wf.status);
          setPublishedVersionId(wf.published_version_id ?? null);
          setDraftVersionId(wf.draft_version_id ?? null);
          setVersion(wf.version ?? 0);
          setTriggers(
            (wf.triggers ?? []).map((t) => ({
              type: t.type,
              config: t.config ?? {},
            }))
          );
          if (wf.conditions) setConditions(wf.conditions);
          setActions(
            (wf.actions ?? []).map((a) => ({
              type: a.type,
              config: a.config ?? {},
              max_retries: a.max_retries,
              continue_on_error: a.continue_on_error,
            }))
          );
          const vers = await listWorkflowVersions(id);
          if (active) setVersions(vers);
        }
      } catch {
        toast.error("Failed to load workflow builder");
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [id, isNew]);

  const moduleFields =
    meta?.modules.find((m) => m.id === moduleId)?.fields ?? [];

  async function save(publish = false, publishChangelog = "") {
    if (!name.trim()) {
      toast.error("Name is required");
      return;
    }
    const payload = {
      name: name.trim(),
      description,
      module_id: moduleId || null,
      on_action_error: onError,
      triggers,
      conditions:
        (conditions.children?.length ?? 0) > 0 || conditions.node_type === "predicate"
          ? conditions
          : null,
      actions,
    };
    try {
      setSaving(true);
      let wfId = id;
      if (isNew) {
        const created = await createWorkflow(payload);
        wfId = created.id;
        toast.success("Draft created");
      } else {
        await updateWorkflow(id, payload);
        toast.success("Draft saved");
      }
      if (publish) {
        const published = await publishWorkflow(wfId, publishChangelog);
        setStatus(published.status);
        setPublishedVersionId(published.published_version_id ?? null);
        setDraftVersionId(published.draft_version_id ?? null);
        setVersion(published.version ?? 0);
        setVersions(await listWorkflowVersions(wfId));
        toast.success(`Published v${published.version}`);
        setShowPublishDialog(false);
        setChangelog("");
      }
      if (isNew || publish) {
        router.push(`/settings/automation/workflows/${wfId}`);
      } else {
        const wf = await getWorkflow(wfId);
        setStatus(wf.status);
        setDraftVersionId(wf.draft_version_id ?? null);
        setVersion(wf.version ?? 0);
        setVersions(await listWorkflowVersions(wfId));
      }
    } catch (e: unknown) {
      const msg =
        (e as { response?: { data?: { message?: string } } })?.response?.data
          ?.message ?? "Save failed";
      toast.error(msg);
    } finally {
      setSaving(false);
    }
  }

  async function onRollback(versionId: string, versionNum: number) {
    if (
      !confirm(
        `Rollback to v${versionNum}? This publishes a new version cloned from that snapshot and activates the workflow.`
      )
    ) {
      return;
    }
    try {
      const wf = await rollbackWorkflow(id, versionId);
      toast.success(`Rolled back — now live as v${wf.version}`);
      setName(wf.name);
      setStatus(wf.status);
      setPublishedVersionId(wf.published_version_id ?? null);
      setDraftVersionId(wf.draft_version_id ?? null);
      setVersion(wf.version ?? 0);
      setTriggers(wf.triggers.map((t) => ({ type: t.type, config: t.config })));
      setActions(wf.actions.map((a) => ({ type: a.type, config: a.config })));
      if (wf.conditions) setConditions(wf.conditions);
      setVersions(await listWorkflowVersions(id));
    } catch {
      toast.error("Rollback failed");
    }
  }

  if (loading || !meta) {
    return (
      <div className="rounded-3xl border border-slate-200 bg-white p-8 text-slate-400">
        Loading builder…
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title={isNew ? "New workflow" : "Edit workflow"}
        description={
          isNew
            ? "Create a draft, then publish when ready."
            : `Status: ${status} · editing v${version}${draftVersionId ? " (draft)" : ""}${publishedVersionId ? " · has published version" : ""}`
        }
        action={
          <div className="flex flex-wrap gap-2">
            <Link
              href="/settings/automation/workflows"
              className="rounded-xl border border-slate-200 px-3 py-2 text-sm"
            >
              Back
            </Link>
            <button
              type="button"
              disabled={saving}
              onClick={() => save(false)}
              className="inline-flex items-center gap-2 rounded-xl border border-slate-200 px-3 py-2 text-sm"
            >
              <Save className="h-4 w-4" /> Save draft
            </button>
            <button
              type="button"
              disabled={saving}
              onClick={() => setShowPublishDialog(true)}
              className="inline-flex items-center gap-2 rounded-xl bg-slate-900 px-3 py-2 text-sm text-white"
            >
              <Power className="h-4 w-4" /> Publish…
            </button>
          </div>
        }
      />

      {!isNew && (
        <div className="flex flex-wrap gap-2 text-xs">
          <span
            className={`rounded-full px-2.5 py-1 font-medium ${
              status === "active"
                ? "bg-emerald-50 text-emerald-700"
                : status === "draft"
                  ? "bg-slate-100 text-slate-700"
                  : "bg-amber-50 text-amber-800"
            }`}
          >
            {status}
          </span>
          {draftVersionId && (
            <span className="rounded-full bg-sky-50 px-2.5 py-1 font-medium text-sky-800">
              Unpublished draft changes
            </span>
          )}
          {publishedVersionId && (
            <span className="rounded-full bg-violet-50 px-2.5 py-1 font-medium text-violet-800">
              Live version published
            </span>
          )}
        </div>
      )}

      {showPublishDialog && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 p-4">
          <div className="w-full max-w-md space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-xl">
            <h3 className="text-lg font-semibold text-slate-900">
              Publish workflow
            </h3>
            <p className="text-sm text-slate-500">
              Saves the current definition as a new published version and sets
              status to active. Only published versions execute.
            </p>
            <FormTextarea
              label="Changelog (optional)"
              value={changelog}
              onChange={(e) => setChangelog(e.target.value)}
              placeholder="What changed in this version?"
            />
            <div className="flex justify-end gap-2">
              <button
                type="button"
                className="rounded-xl border border-slate-200 px-3 py-2 text-sm"
                onClick={() => setShowPublishDialog(false)}
              >
                Cancel
              </button>
              <button
                type="button"
                disabled={saving}
                className="rounded-xl bg-slate-900 px-3 py-2 text-sm text-white"
                onClick={() => save(true, changelog)}
              >
                {saving ? "Publishing…" : "Confirm publish"}
              </button>
            </div>
          </div>
        </div>
      )}

      <section className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <FormInput label="Name" value={name} onChange={(e) => setName(e.target.value)} />
        <FormTextarea
          label="Description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />
        <FormSelect
          label="Module"
          value={moduleId}
          onChange={(e) => setModuleId(e.target.value)}
          options={[
            { value: "", label: "Select module" },
            ...meta.modules.map((m) => ({
              value: m.id,
              label: `${m.label} (${m.api_name})`,
            })),
          ]}
        />
        <FormSelect
          label="On action error"
          value={onError}
          onChange={(e) => setOnError(e.target.value as "stop" | "continue")}
          options={[
            { value: "stop", label: "Stop workflow" },
            { value: "continue", label: "Continue next actions" },
          ]}
        />
      </section>

      <section className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div className="flex items-center justify-between">
          <h3 className="font-semibold text-slate-900">Triggers</h3>
          <button
            type="button"
            className="inline-flex items-center gap-1 text-sm text-slate-700"
            onClick={() =>
              setTriggers((t) => [...t, { type: "record_updated", config: {} }])
            }
          >
            <Plus className="h-4 w-4" /> Add
          </button>
        </div>
        {triggers.map((t, i) => (
          <div key={i} className="grid gap-3 rounded-2xl border border-slate-100 p-4 md:grid-cols-3">
            <FormSelect
              label="Type"
              value={t.type}
              onChange={(e) => {
                const next = [...triggers];
                next[i] = { ...next[i], type: e.target.value };
                setTriggers(next);
              }}
              options={meta.triggers.map((tr) => ({
                value: tr.type,
                label: tr.label,
              }))}
            />
            {t.type === "field_updated" && (
              <FormSelect
                label="Field"
                value={String(t.config?.field_api_name ?? "")}
                onChange={(e) => {
                  const next = [...triggers];
                  next[i] = {
                    ...next[i],
                    config: { ...next[i].config, field_api_name: e.target.value },
                  };
                  setTriggers(next);
                }}
                options={[
                  { value: "", label: "Select field" },
                  ...moduleFields.map((f) => ({
                    value: f.api_name,
                    label: f.label,
                  })),
                ]}
              />
            )}
            {(t.type === "scheduled" || t.type === "date_based") && (
              <>
                {t.type === "scheduled" && (
                  <>
                    <FormInput
                      label="Hour (UTC)"
                      value={String(t.config?.hour ?? "")}
                      onChange={(e) => {
                        const next = [...triggers];
                        next[i] = {
                          ...next[i],
                          config: {
                            ...next[i].config,
                            hour: Number(e.target.value) || 0,
                          },
                        };
                        setTriggers(next);
                      }}
                    />
                    <FormInput
                      label="Minute"
                      value={String(t.config?.minute ?? "0")}
                      onChange={(e) => {
                        const next = [...triggers];
                        next[i] = {
                          ...next[i],
                          config: {
                            ...next[i].config,
                            minute: Number(e.target.value) || 0,
                          },
                        };
                        setTriggers(next);
                      }}
                    />
                    <FormInput
                      label="Batch size"
                      value={String(t.config?.batch_size ?? "100")}
                      onChange={(e) => {
                        const next = [...triggers];
                        next[i] = {
                          ...next[i],
                          config: {
                            ...next[i].config,
                            batch_size: Number(e.target.value) || 100,
                          },
                        };
                        setTriggers(next);
                      }}
                    />
                  </>
                )}
                {t.type === "date_based" && (
                  <>
                    <FormSelect
                      label="Date field"
                      value={String(t.config?.field_api_name ?? "")}
                      onChange={(e) => {
                        const next = [...triggers];
                        next[i] = {
                          ...next[i],
                          config: {
                            ...next[i].config,
                            field_api_name: e.target.value,
                          },
                        };
                        setTriggers(next);
                      }}
                      options={[
                        { value: "", label: "Select field" },
                        ...moduleFields.map((f) => ({
                          value: f.api_name,
                          label: f.label,
                        })),
                      ]}
                    />
                    <FormInput
                      label="Offset days"
                      value={String(t.config?.offset_days ?? "0")}
                      onChange={(e) => {
                        const next = [...triggers];
                        next[i] = {
                          ...next[i],
                          config: {
                            ...next[i].config,
                            offset_days: Number(e.target.value) || 0,
                          },
                        };
                        setTriggers(next);
                      }}
                    />
                  </>
                )}
              </>
            )}
            <button
              type="button"
              className="self-end text-rose-600"
              onClick={() => setTriggers(triggers.filter((_, j) => j !== i))}
            >
              <Trash2 className="h-4 w-4" />
            </button>
          </div>
        ))}
      </section>

      <section className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div className="flex items-center justify-between">
          <h3 className="font-semibold text-slate-900">Conditions</h3>
          <button
            type="button"
            className="inline-flex items-center gap-1 text-sm"
            onClick={() =>
              setConditions((c) => ({
                ...c,
                children: [
                  ...(c.children ?? []),
                  {
                    node_type: "predicate",
                    field_api_name: moduleFields[0]?.api_name ?? "",
                    operator: "eq",
                    value: "",
                  },
                ],
              }))
            }
          >
            <Plus className="h-4 w-4" /> Add predicate
          </button>
        </div>
        <FormSelect
          label="Group logic"
          value={conditions.logic ?? "and"}
          onChange={(e) =>
            setConditions((c) => ({
              ...c,
              logic: e.target.value as "and" | "or",
            }))
          }
          options={[
            { value: "and", label: "AND" },
            { value: "or", label: "OR" },
          ]}
        />
        {(conditions.children ?? []).map((child, i) => (
          <div
            key={i}
            className="grid gap-3 rounded-2xl border border-slate-100 p-4 md:grid-cols-4"
          >
            <FormSelect
              label="Field"
              value={child.field_api_name ?? ""}
              onChange={(e) => {
                const children = [...(conditions.children ?? [])];
                children[i] = { ...children[i], field_api_name: e.target.value };
                setConditions({ ...conditions, children });
              }}
              options={moduleFields.map((f) => ({
                value: f.api_name,
                label: f.label,
              }))}
            />
            <FormSelect
              label="Operator"
              value={child.operator ?? "eq"}
              onChange={(e) => {
                const children = [...(conditions.children ?? [])];
                children[i] = { ...children[i], operator: e.target.value };
                setConditions({ ...conditions, children });
              }}
              options={meta.operators.map((o) => ({
                value: o.key,
                label: o.label,
              }))}
            />
            <FormInput
              label="Value"
              value={String(child.value ?? "")}
              onChange={(e) => {
                const children = [...(conditions.children ?? [])];
                children[i] = { ...children[i], value: e.target.value };
                setConditions({ ...conditions, children });
              }}
            />
            <button
              type="button"
              className="self-end text-rose-600"
              onClick={() =>
                setConditions({
                  ...conditions,
                  children: (conditions.children ?? []).filter((_, j) => j !== i),
                })
              }
            >
              <Trash2 className="h-4 w-4" />
            </button>
          </div>
        ))}
      </section>

      <section className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div className="flex items-center justify-between">
          <h3 className="font-semibold text-slate-900">Actions</h3>
          <button
            type="button"
            className="inline-flex items-center gap-1 text-sm"
            onClick={() =>
              setActions((a) => [
                ...a,
                { type: "create_activity", config: { description: "Workflow ran" } },
              ])
            }
          >
            <Plus className="h-4 w-4" /> Add action
          </button>
        </div>
        {actions.map((a, i) => (
          <div key={i} className="space-y-3 rounded-2xl border border-slate-100 p-4">
            <div className="grid gap-3 md:grid-cols-2">
              <FormSelect
                label="Action"
                value={a.type}
                onChange={(e) => {
                  const next = [...actions];
                  next[i] = { ...next[i], type: e.target.value, config: {} };
                  setActions(next);
                }}
                options={meta.actions
                  .filter((x) => x.mvp)
                  .map((x) => ({ value: x.type, label: x.label }))}
              />
              <button
                type="button"
                className="self-end justify-self-end text-rose-600"
                onClick={() => setActions(actions.filter((_, j) => j !== i))}
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>
            {(a.type === "send_email" || a.type === "send_whatsapp") && (
              <>
                <FormInput
                  label="To"
                  value={String(a.config?.to ?? "")}
                  onChange={(e) => {
                    const next = [...actions];
                    next[i] = {
                      ...next[i],
                      config: { ...next[i].config, to: e.target.value },
                    };
                    setActions(next);
                  }}
                />
                <FormInput
                  label="Subject / body"
                  value={String(a.config?.subject ?? a.config?.body ?? "")}
                  onChange={(e) => {
                    const next = [...actions];
                    next[i] = {
                      ...next[i],
                      config: {
                        ...next[i].config,
                        subject: e.target.value,
                        body: e.target.value || " ",
                      },
                    };
                    setActions(next);
                  }}
                />
              </>
            )}
            {a.type === "assign_owner" && (
              <FormSelect
                label="Owner"
                value={String(a.config?.owner_id ?? "")}
                onChange={(e) => {
                  const next = [...actions];
                  next[i] = {
                    ...next[i],
                    config: { ...next[i].config, owner_id: e.target.value },
                  };
                  setActions(next);
                }}
                options={[
                  { value: "", label: "Select user" },
                  ...meta.users.map((u) => ({
                    value: u.id,
                    label: `${u.name} (${u.email})`,
                  })),
                ]}
              />
            )}
            {a.type === "update_record" && (
              <FormTextarea
                label="Fields JSON"
                value={JSON.stringify(a.config?.fields ?? {}, null, 2)}
                onChange={(e) => {
                  try {
                    const fields = JSON.parse(e.target.value || "{}");
                    const next = [...actions];
                    next[i] = { ...next[i], config: { fields } };
                    setActions(next);
                  } catch {
                    /* ignore partial JSON */
                  }
                }}
              />
            )}
            {a.type === "create_note" && (
              <FormTextarea
                label="Note body"
                value={String(a.config?.body ?? "")}
                onChange={(e) => {
                  const next = [...actions];
                  next[i] = {
                    ...next[i],
                    config: { ...next[i].config, body: e.target.value },
                  };
                  setActions(next);
                }}
              />
            )}
            {a.type === "create_activity" && (
              <FormInput
                label="Description"
                value={String(a.config?.description ?? "")}
                onChange={(e) => {
                  const next = [...actions];
                  next[i] = {
                    ...next[i],
                    config: { ...next[i].config, description: e.target.value },
                  };
                  setActions(next);
                }}
              />
            )}
            {a.type === "delay" && (
              <FormInput
                label="Seconds"
                value={String(a.config?.seconds ?? "60")}
                onChange={(e) => {
                  const next = [...actions];
                  next[i] = {
                    ...next[i],
                    config: { seconds: Number(e.target.value) || 60 },
                  };
                  setActions(next);
                }}
              />
            )}
            {a.type === "webhook" && (
              <FormInput
                label="URL"
                value={String(a.config?.url ?? "")}
                onChange={(e) => {
                  const next = [...actions];
                  next[i] = {
                    ...next[i],
                    config: { ...next[i].config, url: e.target.value },
                  };
                  setActions(next);
                }}
              />
            )}
            {a.type === "invoke_workflow" && (
              <FormInput
                label="Workflow ID"
                value={String(a.config?.workflow_id ?? "")}
                onChange={(e) => {
                  const next = [...actions];
                  next[i] = {
                    ...next[i],
                    config: { workflow_id: e.target.value },
                  };
                  setActions(next);
                }}
              />
            )}
          </div>
        ))}
      </section>

      {!isNew && versions.length > 0 && (
        <section className="space-y-3 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <div>
            <h3 className="font-semibold text-slate-900">Version history</h3>
            <p className="text-xs text-slate-500">
              Every publish creates an immutable version. Rollback clones a past
              snapshot into a new published version.
            </p>
          </div>
          <ul className="space-y-2 text-sm">
            {versions.map((v) => {
              const isLive = publishedVersionId === v.id;
              const stateClass =
                v.state === "published"
                  ? "bg-emerald-50 text-emerald-700"
                  : v.state === "rolled_back"
                    ? "bg-amber-50 text-amber-800"
                    : "bg-slate-100 text-slate-600";
              return (
                <li
                  key={v.id}
                  className="flex flex-wrap items-center justify-between gap-2 rounded-xl border border-slate-100 px-3 py-2"
                >
                  <div className="min-w-0 space-y-0.5">
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="font-medium text-slate-900">
                        v{v.version}
                      </span>
                      <span
                        className={`rounded-full px-2 py-0.5 text-[10px] font-semibold uppercase ${stateClass}`}
                      >
                        {v.state}
                      </span>
                      {isLive && (
                        <span className="rounded-full bg-violet-50 px-2 py-0.5 text-[10px] font-semibold uppercase text-violet-800">
                          live
                        </span>
                      )}
                    </div>
                    <p className="text-xs text-slate-500">
                      {v.published_at
                        ? new Date(v.published_at).toLocaleString()
                        : new Date(v.created_at).toLocaleString()}
                      {v.changelog ? ` · ${v.changelog}` : ""}
                    </p>
                  </div>
                  {(v.state === "published" || v.state === "rolled_back") &&
                    !isLive && (
                      <button
                        type="button"
                        onClick={() => onRollback(v.id, v.version)}
                        className="inline-flex items-center gap-1 rounded-lg border border-slate-200 px-2 py-1 text-xs text-slate-700 hover:bg-slate-50"
                      >
                        <RotateCcw className="h-3 w-3" /> Rollback
                      </button>
                    )}
                </li>
              );
            })}
          </ul>
        </section>
      )}
    </div>
  );
}
