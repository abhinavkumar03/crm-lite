import type { DemoStep } from "./types";

/** Client-side coaching when DB seed copy is thin or route-only. */
export function coachForStep(step: DemoStep): {
  goLabel: string;
  coach?: string;
  hint?: string;
} {
  switch (step.step_key) {
    case "create_module":
      return {
        goLabel: "Open Modules settings",
        coach:
          "Click the highlighted Create module button (top-right). The form is pre-filled — then click Create module in the dialog.",
        hint:
          step.hint ||
          "Settings → Modules → Create module → Create module in the dialog.",
      };
    case "create_field":
      return {
        goLabel: "Open Fields settings",
        coach:
          "Click the highlighted New field button. Values are pre-filled — then click Create field in the dialog.",
        hint:
          step.hint ||
          "Settings → Fields → New field → Create field in the dialog.",
      };
    case "record_workspace":
      return {
        goLabel: "Open Tutorial Lead workspace",
        coach:
          "We’ll open a Tutorial Lead record. Skim the Overview tabs and fields — then press Continue.",
        hint: step.hint || "This is a view step — Continue when you’ve looked around.",
      };
    case "add_note":
      return {
        goLabel: "Open Notes tab",
        coach:
          "You’re taken to Notes with a follow-up already written. Click the highlighted Add button — this step advances automatically.",
        hint:
          step.hint ||
          "Open the Notes tab on a Tutorial Lead and click Add.",
      };
    case "timeline":
      return {
        goLabel: "Open Timeline tab",
        coach:
          "Review the activity stream (creates, notes, updates). When you’ve seen it, press Continue.",
        hint: step.hint || "View-only — Continue after inspecting the Timeline.",
      };
    case "product_demo_module":
      return {
        goLabel: "Open Product Demo table",
        coach:
          "Browse the seeded Product Demo module (fields + sample rows). Press Continue when ready.",
        hint: step.hint || "View-only — explore the table, then Continue.",
      };
    default:
      return {
        goLabel: step.action_label ?? "Go to page",
        hint: step.hint ?? undefined,
      };
  }
}
