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
    case "create_record":
      return {
        goLabel: "Open Tutorial Lead",
        coach:
          "Open Tutorial Lead from the sidebar (or Go). Click Add Tutorial Lead, fill the form, then click Add — this step advances automatically.",
        hint:
          step.hint ||
          "Module → Add Tutorial Lead → Create record (highlighted submit).",
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
        goLabel: "Open Product Demo",
        coach:
          "Open Product Demo from the workspace sidebar. Browse fields and sample rows, then press Continue.",
        hint: step.hint || "View-only — explore the module page, then Continue.",
      };
    case "automation_settings":
      return {
        goLabel: "Open Automation settings",
        coach:
          "Glance at automation / notification controls. Press Continue when you’ve looked around.",
        hint: step.hint || "View-only — open Automation, then Continue.",
      };
    case "export_engine":
      return {
        goLabel: "Open Export (Settings)",
        coach: "Open Settings → Export, explore the page, then press Continue.",
        hint: step.hint || "View-only — Continue when ready.",
      };
    case "import_engine":
      return {
        goLabel: "Open Import (Settings)",
        coach: "Open Settings → Import, explore the page, then press Continue.",
        hint: step.hint || "View-only — Continue when ready.",
      };
    case "validation_rules":
      return {
        goLabel: "Open Validation settings",
        coach: "Glance at validation rules, then press Continue.",
        hint: step.hint || "View-only — Continue when ready.",
      };
    default:
      return {
        goLabel: step.action_label ?? "Go to page",
        hint: step.hint ?? undefined,
      };
  }
}
