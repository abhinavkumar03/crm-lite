"use client";

import { useEffect, useMemo, useState } from "react";

import { usePathname } from "next/navigation";

import { SpotlightLayer, useTargetRect } from "@/features/guided";

import { useDemo } from "./DemoProvider";
import { shouldAutoAdvance } from "./stepAdvance";
import { actionSelectorsForStep, useActionWaiter } from "./useActionWaiter";

function firstMatchingSelector(candidates: string[]): string | null {
  if (typeof document === "undefined") return null;
  for (const sel of candidates) {
    try {
      if (document.querySelector(sel)) return sel;
    } catch {
      // invalid selector — skip
    }
  }
  return null;
}

function normalizeSelector(raw: string): string {
  const t = raw.trim();
  if (t.startsWith("[") || t.startsWith(".") || t.startsWith("#")) return t;
  if (t.includes("=")) return `[${t}]`;
  return `[data-tutorial-action="${t}"], [data-tour="${t}"]`;
}

/**
 * Click-through spotlight for the active demo step. Focuses the target
 * element while allowing interaction inside the hole.
 */
export default function DemoGuideOverlay() {
  const demo = useDemo();
  const pathname = usePathname();
  const [tick, setTick] = useState(0);
  const [modalOpen, setModalOpen] = useState(false);

  const step = demo?.currentStep ?? null;
  const running = demo?.mode === "running" && !!step;

  // Re-resolve targets when modals mount (New module / New field).
  useEffect(() => {
    if (!running) return;
    const id = window.setInterval(() => {
      setTick((n) => n + 1);
      setModalOpen(!!document.querySelector('[data-demo-modal="true"]'));
    }, 250);
    setModalOpen(!!document.querySelector('[data-demo-modal="true"]'));
    return () => window.clearInterval(id);
  }, [running, step?.step_key]);

  const candidates = useMemo(() => {
    if (!step) return [] as string[];
    const fromActions = actionSelectorsForStep(step.step_key);
    const fromTarget = step.target_selector
      ? [normalizeSelector(step.target_selector)]
      : [];
    const submit = fromActions.filter((s) => s.includes("submit-"));
    const openers = fromActions.filter((s) => !s.includes("submit-"));

    // Modal open → highlight Create/submit. Otherwise → highlight New/open button.
    const ordered = modalOpen
      ? [...submit, ...openers, ...fromTarget]
      : [...openers, ...submit, ...fromTarget];

    return Array.from(new Set(ordered));
  }, [step, modalOpen]);

  const selector = useMemo(
    () => (running ? firstMatchingSelector(candidates) : null),
    // tick forces re-query after route/modal paint
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [running, candidates, pathname, tick]
  );

  // Keep measuring while modal is open so submit buttons get a hole.
  const rect = useTargetRect(
    selector,
    running,
    `${step?.step_key ?? ""}:${selector ?? ""}:${tick}:${modalOpen ? "m" : "p"}`,
    pathname
  );

  const autoAdvance = !!step && shouldAutoAdvance(step);

  // Only wire click→validate for create/submit controls (not nav / tab views).
  const autoActionSelectors = useMemo(() => {
    if (!autoAdvance) return [] as string[];
    return candidates.filter(
      (s) =>
        s.includes("submit-") ||
        s.includes("create-") ||
        s.includes('="add-note"') ||
        s.includes("add-field") ||
        s.includes("create-record")
    );
  }, [autoAdvance, candidates]);

  useActionWaiter({
    active: running && !!demo && !demo.busy && autoAdvance,
    actionSelectors: autoActionSelectors,
    onAction: () => {
      window.setTimeout(() => {
        void demo?.validate({ silent: true });
      }, 600);
    },
    pollMs: running && autoAdvance ? 2500 : 0,
    onPoll: () => {
      if (!demo?.busy) void demo?.validate({ silent: true });
    },
  });

  if (!running || !step) return null;

  // Create-record uses a full dynamic form — spotlight would fight the inputs.
  if (step.step_key === "create_record") {
    return null;
  }

  if (!selector && step.validator_key === "none") return null;

  // Modal is z-[70]; raise spotlight above it so Create module/field stays focused.
  // Do not capture outside the hole while a modal is open — blockers steal scroll.
  const zIndex = modalOpen ? 80 : 50;

  return (
    <SpotlightLayer
      rect={rect}
      mode="guide"
      pulse={!!rect}
      zIndex={zIndex}
      captureOutside={!modalOpen}
    />
  );
}
