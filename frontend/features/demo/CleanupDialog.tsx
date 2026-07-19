"use client";

import { Archive, PartyPopper, Trash2 } from "lucide-react";

import { useDemo } from "./DemoProvider";

export default function CleanupDialog() {
  const demo = useDemo();
  if (!demo || demo.mode !== "cleanup" || !demo.session) return null;

  const completed = demo.session.steps.filter(
    (s) => s.status === "completed" || s.status === "skipped"
  ).length;

  return (
    <div className="fixed inset-0 z-[70] flex items-center justify-center bg-slate-900/50 p-4 backdrop-blur-sm">
      <div className="w-full max-w-md rounded-3xl bg-white p-8 shadow-2xl">
        <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-2xl bg-emerald-50 text-emerald-600">
          <PartyPopper size={22} />
        </div>
        <h2 className="text-xl font-bold text-slate-900">
          Congratulations!
        </h2>
        <p className="mt-2 text-sm text-slate-600">
          You completed the CRM walkthrough. {completed} of{" "}
          {demo.session.steps.length} steps finished.
        </p>

        <ul className="mt-4 space-y-1 rounded-2xl bg-slate-50 px-4 py-3 text-sm text-slate-600">
          <li>Progress: {demo.session.progress_percent}%</li>
          <li>Sandbox org kept until you choose below</li>
          <li>Certificate download — coming soon</li>
        </ul>

        <p className="mt-4 text-sm font-medium text-slate-800">
          Keep demo data or delete the sandbox?
        </p>

        <div className="mt-4 space-y-3">
          <button
            type="button"
            disabled={demo.busy}
            onClick={() => demo.cleanup(true)}
            className="
              flex w-full items-center justify-center gap-2
              rounded-2xl bg-emerald-600 px-4 py-3
              text-sm font-semibold text-white hover:bg-emerald-700
              disabled:opacity-60
            "
          >
            <Archive size={16} />
            Keep Demo Data
          </button>
          <button
            type="button"
            disabled={demo.busy}
            onClick={() => demo.cleanup(false)}
            className="
              flex w-full items-center justify-center gap-2
              rounded-2xl border border-red-200 px-4 py-3
              text-sm font-semibold text-red-600 hover:bg-red-50
              disabled:opacity-60
            "
          >
            <Trash2 size={16} />
            Delete Demo Data
          </button>
        </div>
      </div>
    </div>
  );
}
