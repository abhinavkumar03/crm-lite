"use client";

import { useEffect, useState } from "react";

import {
  Check,
  ChevronDown,
  ChevronUp,
  Circle,
  GripVertical,
  Lock,
  Minus,
} from "lucide-react";

import { useDemo } from "./DemoProvider";
import type { DemoStepStatus } from "./types";
import {
  useDraggablePanel,
  walkthroughPanelDefaultPos,
} from "./useDraggablePanel";

const COLLAPSE_KEY = "crm-demo-walkthrough-collapsed";

function StatusIcon({ status }: { status: DemoStepStatus }) {
  if (status === "completed") {
    return <Check size={12} className="text-emerald-600" />;
  }
  if (status === "skipped") {
    return <Minus size={12} className="text-slate-400" />;
  }
  if (status === "active" || status === "failed") {
    return <Circle size={12} className="text-emerald-500" />;
  }
  return <Lock size={12} className="text-slate-300" />;
}

export default function TutorialSidebar() {
  const demo = useDemo();
  const { panelRef, style, dragHandleProps } = useDraggablePanel({
    storageKey: "crm-demo-walkthrough-pos",
    defaultPos: walkthroughPanelDefaultPos,
    fallbackWidth: 224,
  });

  const [collapsed, setCollapsed] = useState(false);

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

  if (!demo || demo.mode !== "running" || !demo.session) return null;

  const activeStep = demo.currentStep;
  const activeIndex = activeStep
    ? demo.session.steps.findIndex((s) => s.step_key === activeStep.step_key)
    : -1;

  return (
    <div
      ref={(node) => {
        panelRef.current = node;
      }}
      data-demo="tutorial-sidebar"
      style={style}
      className="
        pointer-events-auto
        fixed z-[85]
        w-56 overflow-hidden rounded-2xl
        border border-slate-200 bg-white/95
        shadow-lg backdrop-blur
      "
    >
      <div
        {...dragHandleProps}
        className="
          flex cursor-grab items-center gap-2
          border-b border-slate-100 bg-slate-50 px-3 py-2.5
          active:cursor-grabbing touch-none select-none
        "
        title="Drag to move"
      >
        <GripVertical size={14} className="shrink-0 text-slate-300" aria-hidden />
        <p className="min-w-0 flex-1 truncate text-[10px] font-semibold uppercase tracking-widest text-slate-400">
          Walkthrough
        </p>
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            toggleCollapsed();
          }}
          onPointerDown={(e) => e.stopPropagation()}
          className="rounded-lg p-1 text-slate-400 hover:bg-white hover:text-slate-700"
          aria-label={collapsed ? "Expand walkthrough" : "Collapse walkthrough"}
          title={collapsed ? "Expand" : "Collapse"}
        >
          {collapsed ? <ChevronUp size={14} /> : <ChevronDown size={14} />}
        </button>
      </div>

      {collapsed ? (
        <div className="px-3 py-2.5 text-xs text-slate-600">
          <p className="font-semibold text-emerald-800">
            {activeIndex >= 0 ? `Step ${activeIndex + 1}` : "In progress"}
          </p>
          <p className="mt-0.5 line-clamp-2 text-slate-500">
            {activeStep?.title ?? "Walkthrough"}
          </p>
        </div>
      ) : (
        <ol className="max-h-[60vh] space-y-1 overflow-y-auto p-3">
          {demo.session.steps.map((step, i) => {
            const active = step.step_key === demo.session?.current_step_key;
            return (
              <li
                key={step.step_key}
                className={`
                  flex items-start gap-2 rounded-xl px-2 py-1.5 text-xs
                  ${active ? "bg-emerald-50 text-emerald-900" : "text-slate-600"}
                `}
              >
                <span className="mt-0.5 shrink-0">
                  <StatusIcon status={step.status} />
                </span>
                <span className={active ? "font-semibold" : ""}>
                  {i + 1}. {step.title}
                </span>
              </li>
            );
          })}
        </ol>
      )}
    </div>
  );
}
