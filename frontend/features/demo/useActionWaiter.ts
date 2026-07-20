"use client";

import { useEffect, useRef } from "react";

type Options = {
  active: boolean;
  /** Bound to the step that owned this waiter — stale timers must not fire later. */
  stepKey?: string | null;
  /** CSS selectors that should trigger a validate attempt after click. */
  actionSelectors: string[];
  onAction: (selector: string, stepKey: string) => void;
  /** Also poll validate periodically for resource-backed steps. */
  pollMs?: number;
  onPoll?: (stepKey: string) => void;
};

/**
 * Listens for tutorial-relevant clicks and optionally polls for backend
 * completion after the user acts (create module/field/record/note).
 * Timeouts are cleared on step change so a create-step timer cannot
 * auto-complete the following view step.
 */
export function useActionWaiter({
  active,
  stepKey,
  actionSelectors,
  onAction,
  pollMs = 0,
  onPoll,
}: Options) {
  const onActionRef = useRef(onAction);
  const onPollRef = useRef(onPoll);
  onActionRef.current = onAction;
  onPollRef.current = onPoll;

  useEffect(() => {
    if (!active || !stepKey || actionSelectors.length === 0) return;

    const timers: number[] = [];
    const boundKey = stepKey;

    const handler = (e: MouseEvent) => {
      const target = e.target as Element | null;
      if (!target) return;
      for (const sel of actionSelectors) {
        if (target.closest(sel)) {
          timers.push(
            window.setTimeout(() => {
              onActionRef.current(sel, boundKey);
            }, 400)
          );
          return;
        }
      }
    };

    document.addEventListener("click", handler, true);
    return () => {
      document.removeEventListener("click", handler, true);
      timers.forEach((id) => window.clearTimeout(id));
    };
  }, [active, stepKey, actionSelectors]);

  useEffect(() => {
    if (!active || !stepKey || !pollMs || !onPollRef.current) return;
    const boundKey = stepKey;
    const id = window.setInterval(() => {
      onPollRef.current?.(boundKey);
    }, pollMs);
    return () => window.clearInterval(id);
  }, [active, stepKey, pollMs]);
}

/** Map demo step keys to DOM actions that hint progress. */
export function actionSelectorsForStep(stepKey: string): string[] {
  switch (stepKey) {
    case "create_module":
      return [
        '[data-tutorial-action="create-module"]',
        '[data-tutorial-action="submit-module"]',
      ];
    case "create_field":
      return [
        '[data-tutorial-action="add-field"]',
        '[data-tutorial-action="submit-field"]',
      ];
    case "create_record":
      return [
        '[data-tutorial-action="add-record"]',
        '[data-tutorial-action="create-record"]',
      ];
    case "add_note":
      // Only the Add button — not the Notes tab (view navigation).
      return ['[data-tutorial-action="add-note"]'];
    case "timeline":
      return ['[data-tutorial-action="open-timeline-tab"]'];
    case "record_workspace":
      return [
        'a[href^="/m/"]',
        '[data-tour^="nav-module-"]',
      ];
    case "product_demo_module":
      return [
        '[data-tutorial-surface="tables-records"]',
        'a[href="/m/product_demo"]',
        '[data-tour="nav-module-product_demo"]',
      ];
    case "modules_settings":
      return [
        '[data-tutorial-action="open-modules"]',
        'a[href="/settings/modules"]',
        '[data-tour="nav-settings"]',
      ];
    case "dashboard":
      return ['[data-tour="nav-dashboard"]', 'a[href="/dashboard"]'];
    case "roles_glance":
      return [
        '[data-tutorial-action="open-roles"]',
        'a[href="/settings/roles"]',
      ];
    case "validation_rules":
      return ['[data-tutorial-action="open-validation"]'];
    case "automation_settings":
      return ['[data-tutorial-action="open-automation"]'];
    case "import_engine":
      return [
        '[data-tutorial-action="open-imports"]',
        'a[href="/settings/imports"]',
      ];
    case "export_engine":
      return [
        '[data-tutorial-action="open-exports"]',
        'a[href="/settings/exports"]',
      ];
    default:
      return [];
  }
}

/** Selectors for data areas the user should look at (view/inspect steps). */
export function highlightSelectorsForStep(stepKey: string): string[] {
  switch (stepKey) {
    case "record_workspace":
      return ['[data-tutorial-surface="record-overview"]'];
    case "timeline":
      return ['[data-tutorial-surface="record-timeline"]'];
    case "product_demo_module":
      return [
        '[data-tutorial-surface="tables-records"]',
        '[data-tutorial-surface="tables-page"]',
      ];
    case "dashboard":
      return ['[data-tutorial-surface="dashboard"]'];
    case "roles_glance":
      return ['[data-tutorial-surface="roles"]'];
    case "modules_settings":
      return ['[data-tutorial-surface="modules-list"]'];
    case "automation_settings":
      return [
        '[data-tutorial-surface="automation"]',
        '[data-tutorial-action="open-automation"]',
      ];
    case "validation_rules":
      return [
        '[data-tutorial-surface="validation"]',
        '[data-tutorial-action="open-validation"]',
      ];
    case "import_engine":
      return [
        '[data-tutorial-surface="imports"]',
        '[data-tutorial-action="open-imports"]',
      ];
    case "export_engine":
      return [
        '[data-tutorial-surface="exports"]',
        '[data-tutorial-action="open-exports"]',
      ];
    case "settings_sweep":
      return ['[data-tutorial-surface="settings-home"]', '[data-tour="nav-settings"]'];
    default:
      return [];
  }
}

/**
 * @deprecated Prefer shouldAutoAdvance(step) — view steps must not poll
 * just because a resource already exists from an earlier step.
 */
export function shouldPollValidate(validatorKey: string): boolean {
  return [
    "field_exists",
    "record_exists",
    "note_exists",
    "module_exists",
  ].includes(validatorKey);
}
