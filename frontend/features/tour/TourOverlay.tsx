"use client";

import { useEffect, useState } from "react";

import { usePathname } from "next/navigation";

import { ArrowLeft, ArrowRight, X } from "lucide-react";

import {
  GUIDED_CARD_WIDTH,
  SpotlightLayer,
  tooltipPosition,
  useTargetRect,
} from "@/features/guided";

import { useTour } from "./TourProvider";

export default function TourOverlay() {
  const tour = useTour();
  const pathname = usePathname();
  const [vp, setVp] = useState<{ w: number; h: number } | null>(null);

  const step = tour?.currentStep ?? null;
  const active = tour?.active ?? false;
  const stepIndex = tour?.stepIndex ?? 0;

  const rect = useTargetRect(step?.target, active, stepIndex, pathname);

  useEffect(() => {
    const update = () =>
      setVp({ w: window.innerWidth, h: window.innerHeight });
    const id = requestAnimationFrame(update);
    window.addEventListener("resize", update);
    return () => {
      cancelAnimationFrame(id);
      window.removeEventListener("resize", update);
    };
  }, []);

  useEffect(() => {
    if (!active || !tour) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") tour.skip();
      else if (e.key === "ArrowRight") tour.next();
      else if (e.key === "ArrowLeft") tour.back();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [active, tour]);

  if (!tour || !active || !step || !vp) return null;

  const isLast = stepIndex >= tour.totalSteps - 1;
  const pos = tooltipPosition(rect, step.placement, vp);

  return (
    <div className="fixed inset-0 z-[60]">
      <SpotlightLayer rect={rect} mode="block" />

      <div
        role="dialog"
        aria-modal="true"
        aria-label={step.title}
        className="absolute rounded-2xl border border-slate-200 bg-white p-5 shadow-2xl transition-all duration-300"
        style={{ top: pos.top, left: pos.left, width: GUIDED_CARD_WIDTH, zIndex: 61 }}
      >
        <div className="mb-3 flex items-center justify-between">
          <span className="text-xs font-semibold uppercase tracking-widest text-emerald-600">
            Step {stepIndex + 1} of {tour.totalSteps}
          </span>

          <button
            onClick={tour.skip}
            aria-label="Close tour"
            className="rounded-lg p-1 text-slate-400 transition hover:bg-slate-100 hover:text-slate-600"
          >
            <X size={16} />
          </button>
        </div>

        <h3 className="text-lg font-bold text-slate-900">{step.title}</h3>

        <p className="mt-2 text-sm leading-6 text-slate-500">{step.body}</p>

        <div className="mt-4 flex flex-wrap gap-1.5">
          {tour.steps.map((s, i) => (
            <button
              key={s.key}
              onClick={() => tour.goTo(i)}
              aria-label={`Go to step ${i + 1}`}
              className={`h-1.5 rounded-full transition-all ${
                i === stepIndex
                  ? "w-5 bg-emerald-500"
                  : "w-1.5 bg-slate-200 hover:bg-slate-300"
              }`}
            />
          ))}
        </div>

        <div className="mt-5 flex items-center justify-between">
          <button
            onClick={tour.skip}
            className="text-sm font-medium text-slate-400 transition hover:text-slate-600"
          >
            Skip tour
          </button>

          <div className="flex items-center gap-2">
            <button
              onClick={tour.back}
              disabled={stepIndex === 0}
              className="
              flex items-center gap-1 rounded-xl border border-slate-200
              px-3 py-2 text-sm font-medium text-slate-600
              transition hover:bg-slate-50
              disabled:cursor-not-allowed disabled:opacity-40
              "
            >
              <ArrowLeft size={16} />
              Back
            </button>

            <button
              onClick={tour.next}
              className="
              flex items-center gap-1 rounded-xl bg-emerald-500
              px-4 py-2 text-sm font-semibold text-white shadow
              transition hover:bg-emerald-600
              "
            >
              {isLast ? "Finish" : "Next"}
              {!isLast && <ArrowRight size={16} />}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
