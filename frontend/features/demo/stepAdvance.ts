import type { DemoStep } from "./types";

const CREATE_STEP_KEYS = new Set([
  "create_module",
  "create_field",
  "create_record",
  "add_note",
]);

/**
 * Create / mutate steps: auto-validate after the user acts (and poll).
 * View / navigate / acknowledge steps: stay visible until Continue.
 * Only the whitelist drives auto-advance — never required_action alone
 * (mis-seeded DB metadata must not skip view steps).
 */
export function isCreateActionStep(step: DemoStep): boolean {
  return CREATE_STEP_KEYS.has(step.step_key);
}

/** View / navigate steps that stay until the user confirms Continue. */
export function isViewConfirmStep(step: DemoStep): boolean {
  if (step.step_key === "completion") return false;
  return !isCreateActionStep(step);
}

/** Only create-resource steps should silent-auto advance. */
export function shouldAutoAdvance(step: DemoStep): boolean {
  return isCreateActionStep(step);
}
