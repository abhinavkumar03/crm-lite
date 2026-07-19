"use client";

import { useEffect, useState } from "react";

import {
  CheckCircle2,
  ChevronDown,
  ChevronRight,
  ChevronUp,
  GripVertical,
  Loader2,
  RotateCcw,
  SkipForward,
  X,
} from "lucide-react";

import { useDemo } from "./DemoProvider";
import { isCreateActionStep, isViewConfirmStep } from "./stepAdvance";
import { coachForStep } from "./stepCoach";
import {
  instructionPanelDefaultPos,
  useDraggablePanel,
} from "./useDraggablePanel";

const COLLAPSE_KEY = "crm-demo-instruction-collapsed";

export default function InstructionPanel() {
  const demo = useDemo();
  const { panelRef, style, dragHandleProps } = useDraggablePanel({
    storageKey: "crm-demo-instruction-pos",
    defaultPos: instructionPanelDefaultPos,
    fallbackWidth: 340,
  });
  const [detailsOpen, setDetailsOpen] = useState(false);
  const [collapsed, setCollapsed] = useState(false);

  const step = demo?.currentStep ?? null;
  const session = demo?.session ?? null;
  const running = demo?.mode === "running" && !!session && !!step;

  useEffect(() => {
    try {
      setCollapsed(sessionStorage.getItem(COLLAPSE_KEY) === "1");
    } catch {
      // ignore
    }
  }, []);

  function toggleCollapsed() {
    setCollapsed((prev) => {
      const next = !prev;
      try {
        sessionStorage.setItem(COLLAPSE_KEY, next ? "1" : "0");
      } catch {
        // ignore
      }
      return next;
    });
  }

  // Keyboard: Enter = validate/continue, Esc = minimize, S = skip (when allowed)
  useEffect(() => {
    if (!running || !demo || !step) return;
    const onKey = (e: KeyboardEvent) => {
      const tag = (e.target as HTMLElement | null)?.tagName;
      if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") return;
      if ((e.target as HTMLElement | null)?.isContentEditable) return;

      if (e.key === "Escape") {
        e.preventDefault();
        demo.closeUI();
        return;
      }
      if (e.key === "Enter" && !e.metaKey && !e.ctrlKey && !e.altKey) {
        e.preventDefault();
        if (step.step_key === "completion") void demo.finish();
        else void demo.validate();
        return;
      }
      if ((e.key === "s" || e.key === "S") && step.is_skippable) {
        e.preventDefault();
        void demo.skip();
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [running, demo, step]);

  if (!running || !demo || !session || !step) return null;

  const index = session.steps.findIndex((s) => s.step_key === step.step_key);
  const isCompletion = step.step_key === "completion";
  const createStep = isCreateActionStep(step);
  const viewStep = isViewConfirmStep(step);
  const hasDetails = !!(step.why_it_matters || step.how_it_works);
  const coach = coachForStep(step);

  return (
    <aside
      ref={(node) => {
        panelRef.current = node;
      }}
      data-demo="instruction-panel"
      style={style}
      className="
        pointer-events-auto
        fixed z-[85]
        flex w-[min(100vw-2rem,340px)] max-h-[min(70vh,560px)] flex-col
        overflow-hidden rounded-3xl
        border border-slate-200 bg-white
        shadow-2xl shadow-slate-900/15
      "
    >
      <div
        {...dragHandleProps}
        className="
          flex cursor-grab items-start justify-between gap-3
          border-b border-slate-100 bg-slate-50 px-4 py-3
          active:cursor-grabbing touch-none select-none
        "
        title="Drag to move"
      >
        <div className="flex min-w-0 items-start gap-2">
          <GripVertical
            size={16}
            className="mt-1 shrink-0 text-slate-300"
            aria-hidden
          />
          <div className="min-w-0">
            <p className="text-xs font-semibold uppercase tracking-widest text-emerald-600">
              Step {index + 1} of {session.steps.length}
            </p>
            <h3 className="mt-1 text-base font-bold text-slate-900">
              {step.title}
            </h3>
          </div>
        </div>
        <div className="flex shrink-0 items-center gap-0.5">
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              toggleCollapsed();
            }}
            onPointerDown={(e) => e.stopPropagation()}
            className="rounded-xl p-1.5 text-slate-400 hover:bg-white hover:text-slate-700"
            aria-label={collapsed ? "Expand instructions" : "Collapse instructions"}
            title={collapsed ? "Expand" : "Collapse"}
          >
            {collapsed ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
          </button>
          <button
            type="button"
            onClick={() => demo.closeUI()}
            onPointerDown={(e) => e.stopPropagation()}
            className="rounded-xl p-1.5 text-slate-400 hover:bg-white hover:text-slate-700"
            aria-label="Minimize"
            title="Minimize (Esc)"
          >
            <X size={16} />
          </button>
        </div>
      </div>

      {collapsed ? (
        <div className="space-y-2 px-4 py-3">
          <p className="line-clamp-2 text-xs text-slate-500">
            {coach.coach ?? step.description}
          </p>
          <div className="flex gap-2">
            {(step.route ||
              step.step_key === "add_note" ||
              step.step_key === "timeline" ||
              step.step_key === "record_workspace") && (
              <button
                type="button"
                onClick={() => demo.goToStepRoute()}
                className="flex-1 rounded-xl border border-slate-200 px-2 py-1.5 text-xs font-medium text-slate-700 hover:bg-slate-50"
              >
                Go
              </button>
            )}
            <button
              type="button"
              disabled={demo.busy}
              onClick={() =>
                isCompletion ? demo.finish() : demo.validate()
              }
              className="flex-1 rounded-xl bg-emerald-600 px-2 py-1.5 text-xs font-semibold text-white hover:bg-emerald-700 disabled:opacity-60"
            >
              {demo.busy ? "…" : isCompletion ? "Finish" : "Check"}
            </button>
          </div>
        </div>
      ) : (
        <>
          <div className="max-h-[40vh] space-y-3 overflow-y-auto px-5 py-4 text-sm">
            <p className="text-slate-700">{step.description}</p>

            {coach.coach && (
              <div className="rounded-2xl border border-emerald-100 bg-emerald-50/80 px-3 py-2.5 text-xs leading-relaxed text-emerald-900">
                <p className="font-semibold text-emerald-800">What to do</p>
                <p className="mt-1">{coach.coach}</p>
              </div>
            )}

            {hasDetails && (
              <button
                type="button"
                onClick={() => setDetailsOpen((o) => !o)}
                className="flex w-full items-center justify-between rounded-xl bg-slate-50 px-3 py-2 text-xs font-semibold text-slate-600 hover:bg-slate-100"
              >
                Why & how
                <ChevronDown
                  size={14}
                  className={`transition ${detailsOpen ? "rotate-180" : ""}`}
                />
              </button>
            )}

            {detailsOpen && (
              <div className="space-y-3">
                {step.why_it_matters && (
                  <div>
                    <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                      Why this exists
                    </p>
                    <p className="mt-1 text-slate-600">{step.why_it_matters}</p>
                  </div>
                )}
                {step.how_it_works && (
                  <div>
                    <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                      How it works
                    </p>
                    <p className="mt-1 text-slate-600">{step.how_it_works}</p>
                  </div>
                )}
              </div>
            )}

            {step.expected_result && (
              <div className="rounded-2xl bg-emerald-50 px-3 py-2 text-emerald-800">
                <p className="text-xs font-semibold uppercase tracking-wide text-emerald-700">
                  Expected result
                </p>
                <p className="mt-1">{step.expected_result}</p>
              </div>
            )}

            {createStep && (
              <div
                className={`
                  flex items-center gap-2 rounded-2xl px-3 py-2 text-xs font-medium
                  ${
                    demo.stepPhase === "failed"
                      ? "bg-amber-50 text-amber-800"
                      : "bg-slate-100 text-slate-600"
                  }
                `}
              >
                {demo.stepPhase === "validating" || demo.busy ? (
                  <Loader2 size={14} className="animate-spin" />
                ) : null}
                {demo.stepPhase === "failed"
                  ? demo.lastMessage ??
                    "Not done yet — complete the highlighted action"
                  : "Waiting for your action — this step advances automatically when done"}
              </div>
            )}

            {viewStep && !isCompletion && (
              <div className="rounded-2xl bg-sky-50 px-3 py-2 text-xs font-medium text-sky-900">
                Look around this screen, then press{" "}
                <strong>Continue</strong> when you&apos;re ready.
              </div>
            )}

            {demo.lastMessage && !createStep && (
              <p className="text-xs text-slate-500">{demo.lastMessage}</p>
            )}

            {coach.hint && (
              <p
                className={`text-xs ${
                  demo.stepPhase === "failed"
                    ? "text-amber-700"
                    : "text-slate-500"
                }`}
              >
                {demo.stepPhase === "failed" && step.failure_message
                  ? step.failure_message
                  : coach.hint}
              </p>
            )}

            <div className="h-1.5 overflow-hidden rounded-full bg-slate-100">
              <div
                className="h-full rounded-full bg-emerald-500 transition-all"
                style={{ width: `${session.progress_percent}%` }}
              />
            </div>
          </div>

          <div className="space-y-2 border-t border-slate-100 px-5 py-4">
            <button
              type="button"
              onClick={() => demo.goToStepRoute()}
              className="
                flex w-full items-center justify-center gap-2
                rounded-2xl border border-slate-200 px-3 py-2.5
                text-sm font-medium text-slate-700 hover:bg-slate-50
              "
            >
              {coach.goLabel}
              <ChevronRight size={16} />
            </button>

            {isCompletion ? (
              <button
                type="button"
                disabled={demo.busy}
                onClick={() => demo.finish()}
                className="
                  flex w-full items-center justify-center gap-2
                  rounded-2xl bg-emerald-600 px-3 py-2.5
                  text-sm font-semibold text-white hover:bg-emerald-700
                  disabled:opacity-60
                "
              >
                <CheckCircle2 size={16} />
                Finish walkthrough
              </button>
            ) : createStep ? (
              <button
                type="button"
                disabled={demo.busy}
                onClick={() => demo.validate()}
                className="
                  flex w-full items-center justify-center gap-2
                  rounded-2xl border border-emerald-200 bg-emerald-50 px-3 py-2.5
                  text-sm font-semibold text-emerald-800 hover:bg-emerald-100
                  disabled:opacity-60
                "
              >
                <CheckCircle2 size={16} />
                {demo.busy ? "Checking…" : "I finished — check now"}
              </button>
            ) : (
              <button
                type="button"
                disabled={demo.busy}
                onClick={() => demo.validate()}
                className="
                  flex w-full items-center justify-center gap-2
                  rounded-2xl bg-emerald-600 px-3 py-2.5
                  text-sm font-semibold text-white hover:bg-emerald-700
                  disabled:opacity-60
                "
              >
                <CheckCircle2 size={16} />
                {demo.busy ? "Checking…" : "I've seen this — continue"}
              </button>
            )}

            <div className="flex gap-2">
              {step.is_skippable && !isCompletion && (
                <button
                  type="button"
                  disabled={demo.busy}
                  onClick={() => demo.skip()}
                  className="
                    flex flex-1 items-center justify-center gap-1.5
                    rounded-2xl px-3 py-2 text-xs font-medium
                    text-slate-500 hover:bg-slate-50
                  "
                >
                  <SkipForward size={14} />
                  Skip
                </button>
              )}
              <button
                type="button"
                disabled={demo.busy}
                onClick={() => demo.restart()}
                className="
                  flex flex-1 items-center justify-center gap-1.5
                  rounded-2xl px-3 py-2 text-xs font-medium
                  text-slate-500 hover:bg-slate-50
                "
              >
                <RotateCcw size={14} />
                Restart
              </button>
            </div>

            <p className="text-center text-[10px] text-slate-400">
              Drag header · collapse · Enter check · Esc hide
            </p>
          </div>
        </>
      )}
    </aside>
  );
}
