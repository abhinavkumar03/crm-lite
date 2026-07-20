"use client";

import { Suspense, useEffect } from "react";
import { useParams, useRouter, useSearchParams } from "next/navigation";

import { getModules } from "@/features/metadata/api";
import { moduleRecordPath } from "@/features/modules/paths";

function RedirectInner() {
  const params = useParams<{ moduleId: string; recordId: string }>();
  const searchParams = useSearchParams();
  const router = useRouter();

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const mods = await getModules();
        if (!active) return;
        const mod = mods.find((m) => m.id === params.moduleId);
        if (!mod) {
          router.replace("/dashboard");
          return;
        }
        const qs = searchParams.toString();
        const base = moduleRecordPath(mod.api_name, params.recordId);
        router.replace(qs ? `${base}?${qs}` : base);
      } catch {
        router.replace("/dashboard");
      }
    })();
    return () => {
      active = false;
    };
  }, [params.moduleId, params.recordId, router, searchParams]);

  return (
    <div className="py-16 text-center text-sm text-slate-400">
      Opening record…
    </div>
  );
}

/** Legacy UUID record URL → /m/{apiName}/{recordId}. */
export default function LegacyRecordRedirectPage() {
  return (
    <Suspense
      fallback={
        <div className="py-16 text-center text-sm text-slate-400">
          Opening record…
        </div>
      }
    >
      <RedirectInner />
    </Suspense>
  );
}
