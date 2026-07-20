"use client";

import { useEffect, useMemo, useState } from "react";

import { usePathname } from "next/navigation";

import { SpotlightLayer, useTargetRect } from "@/features/guided";

import { useDemo } from "./DemoProvider";
import { isViewConfirmStep, shouldAutoAdvance } from "./stepAdvance";
import {
  actionSelectorsForStep,
  highlightSelectorsForStep,
  useActionWaiter,
} from "./useActionWaiter";

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
  const viewStep = !!step && isViewConfirmStep(step);
  const autoAdvance = !!step && shouldAutoAdvance(step);

  // Re-resolve targets when modals / data surfaces mount.
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

    // View/inspect: prefer large data surfaces so the user sees what to look at.
    if (viewStep) {
      const surfaces = highlightSelectorsForStep(step.step_key);
      const fromTarget = step.target_selector
        ? [normalizeSelector(step.target_selector)]
        : [];
      return Array.from(new Set([...surfaces, ...fromTarget]));
    }

    const fromActions = actionSelectorsForStep(step.step_key);
    const fromTarget = step.target_selector
      ? [normalizeSelector(step.target_selector)]
      : [];
    const submit = fromActions.filter((s) => s.includes("submit-"));
    const openers = fromActions.filter((s) => !s.includes("submit-"));

    const ordered = modalOpen
      ? [...submit, ...openers, ...fromTarget]
      : [...openers, ...submit, ...fromTarget];

    return Array.from(new Set(ordered));
  }, [step, modalOpen, viewStep]);

  const selector = useMemo(
    () => (running ? firstMatchingSelector(candidates) : null),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [running, candidates, pathname, tick]
  );

  const rect = useTargetRect(
    selector,
    running,
    `${step?.step_key ?? ""}:${selector ?? ""}:${tick}:${modalOpen ? "m" : "p"}`,
    pathname
  );

  // Keep the highlighted control / data area in view.
  useEffect(() => {
    if (!running || !selector) return;
    if (!viewStep && step?.step_key !== "create_record") return;
    const el = document.querySelector(selector);
    if (!(el instanceof HTMLElement)) return;
    el.scrollIntoView({ block: "nearest", behavior: "smooth" });
  }, [running, selector, viewStep, step?.step_key]);

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
    stepKey: step?.step_key,
    actionSelectors: autoActionSelectors,
    onAction: (_sel, stepKey) => {
      // stepKey is frozen at click time — ignored if the session already advanced.
      void demo?.validate({ silent: true, stepKey });
    },
    pollMs: running && autoAdvance ? 2500 : 0,
    onPoll: (stepKey) => {
      if (!demo?.busy) void demo?.validate({ silent: true, stepKey });
    },
  });

  if (!running || !step) return null;

  // Form create: highlight the submit button only (no blockers so inputs/scroll work).
  const formCreateStep = step.step_key === "create_record";

  if (!selector && step.validator_key === "none") return null;

  const zIndex = modalOpen ? 80 : 50;
  const softHighlight = viewStep || formCreateStep;

  return (
    <SpotlightLayer
      rect={rect}
      mode="guide"
      pulse={!!rect}
      zIndex={zIndex}
      // View + form-create + open modals: never block scroll/clicks outside the ring.
      captureOutside={!modalOpen && !softHighlight}
      variant={softHighlight ? "highlight" : "focus"}
    />
  );
}
