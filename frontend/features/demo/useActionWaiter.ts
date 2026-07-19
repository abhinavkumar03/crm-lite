"use client";

import { useEffect, useRef } from "react";

type Options = {
  active: boolean;
  /** CSS selectors that should trigger a validate attempt after click. */
  actionSelectors: string[];
  onAction: (selector: string) => void;
  /** Also poll validate periodically for resource-backed steps. */
  pollMs?: number;
  onPoll?: () => void;
};

/**
 * Listens for tutorial-relevant clicks and optionally polls for backend
 * completion after the user acts (create module/field/record/note).
 */
export function useActionWaiter({
  active,
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
    if (!active || actionSelectors.length === 0) return;

    const handler = (e: MouseEvent) => {
      const target = e.target as Element | null;
      if (!target) return;
      for (const sel of actionSelectors) {
        if (target.closest(sel)) {
          // Let the click complete; validate shortly after.
          window.setTimeout(() => onActionRef.current(sel), 400);
          return;
        }
      }
    };

    document.addEventListener("click", handler, true);
    return () => document.removeEventListener("click", handler, true);
  }, [active, actionSelectors]);

  useEffect(() => {
    if (!active || !pollMs || !onPollRef.current) return;
    const id = window.setInterval(() => onPollRef.current?.(), pollMs);
    return () => window.clearInterval(id);
  }, [active, pollMs]);
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
      return ['[data-tutorial-action="create-record"]'];
    case "add_note":
      return [
        '[data-tutorial-action="add-note"]',
        '[data-tutorial-action="open-notes-tab"]',
      ];
    case "timeline":
      return ['[data-tutorial-action="open-timeline-tab"]'];
    case "record_workspace":
      return ['[data-tour="nav-tables"]', 'a[href^="/tables/"]'];
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
