"use client";

import { useEffect, useState } from "react";

import { usePathname } from "next/navigation";

import { ArrowLeft, ArrowRight, X } from "lucide-react";

import { useTour } from "./TourProvider";
import { TourPlacement } from "./steps";

const CARD_WIDTH = 340;
const CARD_HEIGHT = 210;
const MARGIN = 16;
const SPOTLIGHT_PADDING = 8;

// useTargetRect measures the highlighted element, retrying across a few frames
// so it works even while a route transition is still painting.
function useTargetRect(
  selector: string | undefined,
  active: boolean,
  stepIndex: number,
  pathname: string
): DOMRect | null {
  const [rect, setRect] = useState<DOMRect | null>(null);

  useEffect(() => {
    if (!active || !selector) {
      const id = requestAnimationFrame(() => setRect(null));
      return () => cancelAnimationFrame(id);
    }

    let raf = 0;
    let attempts = 0;

    const measure = () => {
      const el = document.querySelector(selector);
      if (el) {
        setRect(el.getBoundingClientRect());
        return true;
      }
      return false;
    };

    const tick = () => {
      if (measure()) return;
      attempts += 1;
      if (attempts < 45) {
        raf = requestAnimationFrame(tick);
      } else {
        setRect(null);
      }
    };
    raf = requestAnimationFrame(tick);

    const onViewportChange = () => measure();
    window.addEventListener("resize", onViewportChange);
    window.addEventListener("scroll", onViewportChange, true);

    return () => {
      cancelAnimationFrame(raf);
      window.removeEventListener("resize", onViewportChange);
      window.removeEventListener("scroll", onViewportChange, true);
    };
  }, [selector, active, stepIndex, pathname]);

  return rect;
}

function tooltipPosition(
  rect: DOMRect | null,
  placement: TourPlacement | undefined,
  vp: { w: number; h: number }
): { top: number; left: number } {
  if (!rect || placement === "center") {
    return {
      top: vp.h / 2 - CARD_HEIGHT / 2,
      left: vp.w / 2 - CARD_WIDTH / 2,
    };
  }

  let top: number;
  let left: number;

  switch (placement) {
    case "right":
      top = rect.top;
      left = rect.right + MARGIN;
      break;
    case "left":
      top = rect.top;
      left = rect.left - CARD_WIDTH - MARGIN;
      break;
    case "top":
      top = rect.top - CARD_HEIGHT - MARGIN;
      left = rect.left;
      break;
    case "bottom":
    default:
      top = rect.bottom + MARGIN;
      left = rect.left;
      break;
  }

  left = Math.max(MARGIN, Math.min(left, vp.w - CARD_WIDTH - MARGIN));
  top = Math.max(MARGIN, Math.min(top, vp.h - CARD_HEIGHT - MARGIN));
  return { top, left };
}

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

  // Keyboard shortcuts: Esc skips, arrows navigate.
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
      {/* Click-blocking + dimming layer. When a target is measured we use the
          box-shadow "hole" trick to spotlight it; otherwise we dim the screen. */}
      {rect ? (
        <div
          aria-hidden
          className="pointer-events-none absolute rounded-xl ring-2 ring-emerald-400 transition-all duration-300"
          style={{
            top: rect.top - SPOTLIGHT_PADDING,
            left: rect.left - SPOTLIGHT_PADDING,
            width: rect.width + SPOTLIGHT_PADDING * 2,
            height: rect.height + SPOTLIGHT_PADDING * 2,
            boxShadow: "0 0 0 9999px rgba(15, 23, 42, 0.55)",
          }}
        />
      ) : (
        <div
          aria-hidden
          className="absolute inset-0 bg-slate-900/55"
        />
      )}

      {/* Transparent capture layer keeps the app behind non-interactive. */}
      <div className="absolute inset-0" />

      {/* Step card */}
      <div
        role="dialog"
        aria-modal="true"
        aria-label={step.title}
        className="absolute rounded-2xl border border-slate-200 bg-white p-5 shadow-2xl transition-all duration-300"
        style={{ top: pos.top, left: pos.left, width: CARD_WIDTH }}
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

        {/* Progress dots */}
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
              flex
              items-center
              gap-1
              rounded-xl
              border
              border-slate-200
              px-3
              py-2
              text-sm
              font-medium
              text-slate-600
              transition
              hover:bg-slate-50
              disabled:cursor-not-allowed
              disabled:opacity-40
              "
            >
              <ArrowLeft size={16} />
              Back
            </button>

            <button
              onClick={tour.next}
              className="
              flex
              items-center
              gap-1
              rounded-xl
              bg-emerald-500
              px-4
              py-2
              text-sm
              font-semibold
              text-white
              shadow
              transition
              hover:bg-emerald-600
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
