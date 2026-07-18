"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Plus } from "lucide-react";

import FormInput from "@/components/common/form/FormInput";
import FormTextarea from "@/components/common/form/FormTextarea";

import {
  createDepartment,
  listDepartments,
} from "@/features/organization/api";
import { StructureItem } from "@/features/organization/types";
import { apiErrorMessage } from "@/features/settings/errors";

export default function DepartmentsSettingsPage() {
  const [items, setItems] = useState<StructureItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [saving, setSaving] = useState(false);

  async function reload() {
    setItems(await listDepartments());
  }

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        await reload();
      } catch (err) {
        if (active) {
          toast.error(apiErrorMessage(err, "Failed to load departments"));
        }
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
      await createDepartment({
        name: name.trim(),
        description: description.trim() || null,
      });
      setName("");
      setDescription("");
      await reload();
      toast.success("Department created");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to create department"));
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return (
      <div className="rounded-3xl border border-slate-200 bg-white p-8 text-slate-400 shadow-sm">
        Loading departments...
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <section className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">Departments</h2>
          <p className="text-sm text-slate-500">
            Org-scoped departments used for membership and record visibility.
          </p>
        </div>

        <div className="grid gap-3 sm:grid-cols-[1fr_1fr_auto]">
          <FormInput
            label="Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <FormTextarea
            label="Description"
            value={description}
            rows={1}
            onChange={(e) => setDescription(e.target.value)}
          />
          <div className="flex items-end">
            <button
              type="button"
              disabled={saving}
              onClick={handleCreate}
              className="inline-flex w-full items-center justify-center gap-2 rounded-full bg-emerald-500 px-4 py-2.5 text-sm font-semibold text-white disabled:opacity-50 sm:w-auto"
            >
              <Plus className="h-4 w-4" />
              Add
            </button>
          </div>
        </div>
      </section>

      <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        {items.length === 0 ? (
          <p className="text-sm text-slate-500">No departments yet.</p>
        ) : (
          <ul className="divide-y divide-slate-100">
            {items.map((d) => (
              <li key={d.id} className="flex flex-col py-3">
                <span className="text-sm font-semibold text-slate-800">
                  {d.name}
                </span>
                {d.description ? (
                  <span className="text-xs text-slate-500">{d.description}</span>
                ) : null}
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
