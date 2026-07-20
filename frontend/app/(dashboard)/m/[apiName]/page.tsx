"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import ModuleRecordsWorkspace from "@/features/modules/components/ModuleRecordsWorkspace";
import { resolveModuleByApiName } from "@/features/modules/paths";
import {
  getModuleFields,
  getValidationSchema,
  getViews,
} from "@/features/metadata/api";
import {
  ModuleField,
  ModuleSummary,
  SavedView,
  ValidationSchema,
} from "@/features/metadata/types";

export default function ModulePage() {
  const params = useParams<{ apiName: string }>();
  const apiName = decodeURIComponent(params.apiName ?? "");
  const router = useRouter();

  const [module, setModule] = useState<ModuleSummary | null>(null);
  const [fields, setFields] = useState<ModuleField[]>([]);
  const [schema, setSchema] = useState<ValidationSchema | null>(null);
  const [views, setViews] = useState<SavedView[]>([]);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        setLoading(true);
        setNotFound(false);
        const mod = await resolveModuleByApiName(apiName);
        if (!active) return;
        if (!mod) {
          setNotFound(true);
          return;
        }
        const [f, s, v] = await Promise.all([
          getModuleFields(mod.id),
          getValidationSchema(mod.id),
          getViews(mod.id),
        ]);
        if (!active) return;
        setModule(mod);
        setFields(f);
        setSchema(s);
        setViews(v);
      } catch {
        toast.error("Failed to load module");
        if (active) setNotFound(true);
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [apiName]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20 text-sm text-slate-400">
        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
        Loading module…
      </div>
    );
  }

  if (notFound || !module || !schema) {
    return (
      <div className="space-y-4 rounded-3xl border border-slate-200 bg-white p-8 text-center shadow-sm">
        <p className="text-slate-600">Module not found.</p>
        <button
          type="button"
          onClick={() => router.push("/dashboard")}
          className="text-sm font-semibold text-emerald-600 hover:text-emerald-700"
        >
          Back to dashboard
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-8" data-tutorial-surface="tables-page">
      {fields.length === 0 ? (
        <div className="rounded-3xl border border-slate-200 bg-white p-8 text-center text-slate-500 shadow-sm">
          This module has no fields yet. Add fields in Settings → Fields.
        </div>
      ) : (
        <ModuleRecordsWorkspace
          key={module.id}
          moduleId={module.id}
          apiName={module.api_name}
          moduleLabel={module.singular_label}
          pluralLabel={module.plural_label}
          fields={fields}
          schema={schema}
          initialViews={views}
        />
      )}
    </div>
  );
}
