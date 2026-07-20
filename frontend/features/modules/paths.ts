import { getModules } from "@/features/metadata/api";
import type { ModuleSummary } from "@/features/metadata/types";

/** Resolve a sidebar/route api_name to a module summary. */
export async function resolveModuleByApiName(
  apiName: string
): Promise<ModuleSummary | null> {
  const normalized = apiName.trim().toLowerCase();
  if (!normalized) return null;
  const modules = await getModules();
  return (
    modules.find((m) => m.api_name.toLowerCase() === normalized) ?? null
  );
}

export function moduleListPath(apiName: string): string {
  return `/m/${encodeURIComponent(apiName)}`;
}

export function moduleRecordPath(apiName: string, recordId: string): string {
  return `/m/${encodeURIComponent(apiName)}/${recordId}`;
}
