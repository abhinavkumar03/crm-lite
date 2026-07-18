"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Plus } from "lucide-react";

import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";
import FormTextarea from "@/components/common/form/FormTextarea";

import {
  createTeam,
  listDepartments,
  listTeams,
} from "@/features/organization/api";
import { StructureItem } from "@/features/organization/types";
import { apiErrorMessage } from "@/features/settings/errors";

export default function TeamsSettingsPage() {
  const [items, setItems] = useState<StructureItem[]>([]);
  const [departments, setDepartments] = useState<StructureItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [departmentId, setDepartmentId] = useState("");
  const [saving, setSaving] = useState(false);

  async function reload() {
    const [t, d] = await Promise.all([listTeams(), listDepartments()]);
    setItems(t);
    setDepartments(d);
  }

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        await reload();
      } catch (err) {
        if (active) toast.error(apiErrorMessage(err, "Failed to load teams"));
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, []);

  async function handleCreate() {
    if (!name.trim()) {
      toast.error("Name is required");
      return;
    }
    try {
      setSaving(true);
      await createTeam({
        name: name.trim(),
        description: description.trim() || null,
        department_id: departmentId || null,
      });
      setName("");
      setDescription("");
      setDepartmentId("");
      await reload();
      toast.success("Team created");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to create team"));
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return (
      <div className="rounded-3xl border border-slate-200 bg-white p-8 text-slate-400 shadow-sm">
        Loading teams...
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <section className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">Teams</h2>
          <p className="text-sm text-slate-500">
            Teams can optionally sit under a department for visibility scoping.
          </p>
        </div>

        <div className="grid gap-3 sm:grid-cols-2">
          <FormInput
            label="Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <FormSelect
            label="Department"
            value={departmentId}
            onChange={(e) => setDepartmentId(e.target.value)}
          >
            <option value="">None</option>
            {departments.map((d) => (
              <option key={d.id} value={d.id}>
                {d.name}
              </option>
            ))}
          </FormSelect>
          <FormTextarea
            label="Description"
            value={description}
            rows={2}
            onChange={(e) => setDescription(e.target.value)}
          />
          <div className="flex items-end">
            <button
              type="button"
              disabled={saving}
              onClick={handleCreate}
              className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-4 py-2.5 text-sm font-semibold text-white disabled:opacity-50"
            >
              <Plus className="h-4 w-4" />
              Add team
            </button>
          </div>
        </div>
      </section>

      <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        {items.length === 0 ? (
          <p className="text-sm text-slate-500">No teams yet.</p>
        ) : (
          <ul className="divide-y divide-slate-100">
            {items.map((t) => {
              const dept = departments.find((d) => d.id === t.department_id);
              return (
                <li key={t.id} className="flex flex-col py-3">
                  <span className="text-sm font-semibold text-slate-800">
                    {t.name}
                  </span>
                  <span className="text-xs text-slate-500">
                    {dept ? `Dept: ${dept.name}` : "No department"}
                    {t.description ? ` · ${t.description}` : ""}
                  </span>
                </li>
              );
            })}
          </ul>
        )}
      </section>
    </div>
  );
}
