import { getModules, listRecords } from "@/features/metadata/api";
import { moduleRecordPath } from "@/features/modules/paths";

export type TutorialWorkspaceTarget = {
  moduleId: string;
  apiName: string;
  recordId: string;
  tab: "overview" | "notes" | "timeline";
};

/**
 * Resolves a Tutorial Lead (or fallback) record for workspace-guided demo steps.
 */
export async function resolveTutorialWorkspace(
  tab: TutorialWorkspaceTarget["tab"] = "overview",
  moduleApiName = "tutorial_lead"
): Promise<TutorialWorkspaceTarget | null> {
  const modules = await getModules();
  const mod =
    modules.find((m) => m.api_name === moduleApiName) ??
    modules.find((m) => m.storage_strategy === "dynamic");
  if (!mod) return null;

  const result = await listRecords(mod.id, {
    page: 1,
    page_size: 1,
    sort: "created_at",
    order: "desc",
  });
  const record = result.records?.[0];
  if (!record) return null;

  return {
    moduleId: mod.id,
    apiName: mod.api_name,
    recordId: record.id,
    tab,
  };
}

export function workspacePath(target: TutorialWorkspaceTarget): string {
  return `${moduleRecordPath(target.apiName, target.recordId)}?tab=${target.tab}`;
}

/** Map demo step keys to the workspace tab they should open. */
export function workspaceTabForStep(
  stepKey: string
): TutorialWorkspaceTarget["tab"] | null {
  switch (stepKey) {
    case "record_workspace":
      return "overview";
    case "add_note":
      return "notes";
    case "timeline":
      return "timeline";
    default:
      return null;
  }
}
