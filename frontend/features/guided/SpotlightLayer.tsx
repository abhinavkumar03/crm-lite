"use client";

import type { SpotlightInteractionMode } from "./types";
import { SPOTLIGHT_PADDING } from "./placement";

type Props = {
  rect: DOMRect | null;
  /** block: capture clicks outside hole. guide: hole is click-through to the app. */
  mode: SpotlightInteractionMode;
  pulse?: boolean;
  zIndex?: number;
  /**
   * When false (e.g. over a scrollable modal), only draw the ring/dim —
   * no blocker panes — so wheel/scroll and form interaction keep working.
   */
  captureOutside?: boolean;
  /**
   * focus = mentor spotlight (dim + optional blockers).
   * highlight = view-only ring around a data area; never blocks scroll.
   */
  variant?: "focus" | "highlight";
};

/**
 * Dims the page with a spotlight hole. In guide mode the wrapper is
 * pointer-events-none so only the four blocker panes capture clicks —
 * the hole itself passes events through to the highlighted control.
 */
export default function SpotlightLayer({
  rect,
  mode,
  pulse = false,
  zIndex = 60,
  captureOutside = true,
  variant = "focus",
}: Props) {
  const pad = SPOTLIGHT_PADDING;
  const isHighlight = variant === "highlight";

  if (!rect) {
    // View highlight: wait silently until the data area mounts.
    if (isHighlight || mode === "guide") {
      if (isHighlight) return null;
      return (
        <div
          className="pointer-events-none fixed inset-0 bg-slate-900/40"
          style={{ zIndex }}
          aria-hidden
        />
      );
    }
    return (
      <div
        className="fixed inset-0 bg-slate-900/55"
        style={{ zIndex }}
        aria-hidden
      />
    );
  }

  const top = rect.top - pad;
  const left = rect.left - pad;
  const width = rect.width + pad * 2;
  const height = rect.height + pad * 2;

  const showBlockers =
    !isHighlight &&
    (mode === "block" || (mode === "guide" && captureOutside));

  return (
    <div
      className="pointer-events-none fixed inset-0"
      style={{ zIndex }}
      aria-hidden
    >
      <div
        className={`
          absolute rounded-xl ring-2 ring-emerald-400
          transition-all duration-300
          ${pulse ? "animate-pulse" : ""}
        `}
        style={{
          top,
          left,
          width,
          height,
          boxShadow: isHighlight
            ? "0 0 0 3px rgba(16, 185, 129, 0.35), 0 0 0 9999px rgba(15, 23, 42, 0.28)"
            : "0 0 0 9999px rgba(15, 23, 42, 0.55)",
        }}
      />

      {mode === "block" && !isHighlight ? (
        <div className="pointer-events-auto absolute inset-0" />
      ) : showBlockers ? (
        <>
          <div
            className="pointer-events-auto absolute left-0 right-0 top-0"
            style={{ height: Math.max(0, top) }}
          />
          <div
            className="pointer-events-auto absolute left-0"
            style={{ top, height, width: Math.max(0, left) }}
          />
          <div
            className="pointer-events-auto absolute"
            style={{
              top,
              height,
              left: left + width,
              right: 0,
            }}
          />
          <div
            className="pointer-events-auto absolute bottom-0 left-0 right-0"
            style={{ top: top + height }}
          />
        </>
      ) : null}
    </div>
  );
}
