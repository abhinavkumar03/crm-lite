"use client";

import { Check, Clock, Sparkles, X } from "lucide-react";

import { useDemo } from "./DemoProvider";

export default function DemoCatalogModal() {
  const demo = useDemo();
  if (!demo || demo.mode !== "catalog") return null;

  const catalog = demo.catalog;

  return (
    <div className="fixed inset-0 z-[70] flex items-center justify-center bg-slate-900/50 p-4 backdrop-blur-sm">
      <div className="relative w-full max-w-lg overflow-hidden rounded-3xl bg-white shadow-2xl">
        <button
          type="button"
          onClick={() => demo.closeUI()}
          className="absolute right-4 top-4 rounded-xl p-2 text-slate-400 hover:bg-slate-100 hover:text-slate-700"
          aria-label="Close"
        >
          <X size={18} />
        </button>

        <div className="bg-gradient-to-br from-emerald-600 to-teal-600 px-8 py-8 text-white">
          <div className="mb-3 flex items-center gap-2 text-sm font-medium text-emerald-100">
            <Sparkles size={16} />
            Interactive sandbox tutorial
          </div>
          <h2 className="text-2xl font-bold tracking-tight">
            {catalog?.name ?? "Interactive CRM Walkthrough"}
          </h2>
          <p className="mt-2 text-sm text-emerald-50/90">
            {catalog?.description ??
              "Learn every major capability by doing real actions in an isolated workspace."}
          </p>
        </div>

        <div className="space-y-6 px-8 py-6">
          <div className="flex items-center gap-3 text-sm text-slate-600">
            <Clock size={16} className="text-emerald-600" />
            <span>
              Estimated duration:{" "}
              <strong className="text-slate-900">
                {catalog?.duration_min ?? 15} minutes
              </strong>
            </span>
          </div>

          <div>
            <p className="mb-3 text-xs font-semibold uppercase tracking-widest text-slate-400">
              Features covered
            </p>
            <ul className="grid grid-cols-2 gap-2">
              {(catalog?.features ?? []).map((feature) => (
                <li
                  key={feature}
                  className="flex items-center gap-2 text-sm text-slate-700"
                >
                  <Check size={14} className="shrink-0 text-emerald-600" />
                  {feature}
                </li>
              ))}
            </ul>
          </div>

          <p className="rounded-2xl bg-slate-50 px-4 py-3 text-xs leading-relaxed text-slate-500">
            Starting creates a temporary sandbox organization with seeded demo
            data. Your real CRM records are never touched. You can resume,
            restart, or delete the sandbox anytime.
          </p>

          <div className="flex gap-3">
            <button
              type="button"
              disabled={demo.busy}
              onClick={() => demo.start()}
              className="
                flex-1 rounded-2xl bg-emerald-600 px-4 py-3
                text-sm font-semibold text-white
                transition hover:bg-emerald-700
                disabled:opacity-60
              "
            >
              {demo.busy ? "Provisioning sandbox…" : "Start Demo"}
            </button>
            <button
              type="button"
              onClick={() => demo.closeUI()}
              className="rounded-2xl px-4 py-3 text-sm font-medium text-slate-600 hover:bg-slate-100"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
