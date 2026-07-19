"use client";

import { Play, RotateCcw, X } from "lucide-react";

import { useDemo } from "./DemoProvider";

export default function ResumeDialog() {
  const demo = useDemo();
  if (!demo || demo.mode !== "resume" || !demo.session) return null;

  const step = demo.currentStep;

  return (
    <div className="fixed inset-0 z-[70] flex items-center justify-center bg-slate-900/50 p-4 backdrop-blur-sm">
      <div className="relative w-full max-w-md rounded-3xl bg-white p-8 shadow-2xl">
        <button
          type="button"
          onClick={() => demo.closeUI()}
          className="absolute right-4 top-4 rounded-xl p-2 text-slate-400 hover:bg-slate-100"
          aria-label="Close"
        >
          <X size={18} />
        </button>

        <h2 className="text-xl font-bold text-slate-900">Resume Demo?</h2>
        <p className="mt-2 text-sm text-slate-600">
          You have an active sandbox walkthrough at{" "}
          <strong>{demo.session.progress_percent}%</strong>
          {step ? (
            <>
              {" "}
              — next up: <em>{step.title}</em>
            </>
          ) : null}
          .
        </p>

        <div className="mt-6 space-y-3">
          <button
            type="button"
            onClick={() => demo.continueSession()}
            className="
              flex w-full items-center justify-center gap-2
              rounded-2xl bg-emerald-600 px-4 py-3
              text-sm font-semibold text-white hover:bg-emerald-700
            "
          >
            <Play size={16} />
            Continue
          </button>
          <button
            type="button"
            disabled={demo.busy}
            onClick={() => demo.restart()}
            className="
              flex w-full items-center justify-center gap-2
              rounded-2xl border border-slate-200 px-4 py-3
              text-sm font-semibold text-slate-700 hover:bg-slate-50
              disabled:opacity-60
            "
          >
            <RotateCcw size={16} />
            Restart
          </button>
          <button
            type="button"
            onClick={() => demo.closeUI()}
            className="w-full rounded-2xl px-4 py-2 text-sm text-slate-500 hover:bg-slate-50"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  );
}
