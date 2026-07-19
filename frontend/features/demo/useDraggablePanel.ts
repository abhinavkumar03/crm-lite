"use client";

import {
  useCallback,
  useEffect,
  useRef,
  useState,
  type PointerEvent as ReactPointerEvent,
} from "react";

type Pos = { x: number; y: number };

type Options = {
  storageKey: string;
  defaultPos: () => Pos;
  /** Approximate width used before the element mounts. */
  fallbackWidth?: number;
};

function clamp(n: number, min: number, max: number) {
  return Math.min(Math.max(n, min), max);
}

function loadPos(storageKey: string, fallback: () => Pos): Pos {
  try {
    const raw = sessionStorage.getItem(storageKey);
    if (!raw) return fallback();
    const parsed = JSON.parse(raw) as Pos;
    if (typeof parsed.x === "number" && typeof parsed.y === "number") {
      return parsed;
    }
  } catch {
    // ignore
  }
  return fallback();
}

/**
 * Pointer-based drag for floating demo panels. Position persists per tab via
 * sessionStorage so users can park overlays away from action targets.
 */
export function useDraggablePanel({
  storageKey,
  defaultPos,
  fallbackWidth = 280,
}: Options) {
  const [pos, setPos] = useState<Pos>(() =>
    typeof window === "undefined" ? { x: 16, y: 120 } : loadPos(storageKey, defaultPos)
  );
  const dragRef = useRef<{
    startX: number;
    startY: number;
    origX: number;
    origY: number;
  } | null>(null);
  const panelRef = useRef<HTMLElement | null>(null);

  const persist = useCallback(
    (next: Pos) => {
      try {
        sessionStorage.setItem(storageKey, JSON.stringify(next));
      } catch {
        // ignore
      }
    },
    [storageKey]
  );

  const onPointerDown = useCallback(
    (e: ReactPointerEvent<HTMLDivElement>) => {
      if (e.button !== 0) return;
      const target = e.target as HTMLElement;
      if (target.closest("button, a, input, textarea, select")) return;

      e.preventDefault();
      e.currentTarget.setPointerCapture(e.pointerId);
      dragRef.current = {
        startX: e.clientX,
        startY: e.clientY,
        origX: pos.x,
        origY: pos.y,
      };
    },
    [pos.x, pos.y]
  );

  const onPointerMove = useCallback(
    (e: ReactPointerEvent<HTMLDivElement>) => {
      if (!dragRef.current) return;
      const dx = e.clientX - dragRef.current.startX;
      const dy = e.clientY - dragRef.current.startY;
      const el = panelRef.current;
      const w = el?.offsetWidth ?? fallbackWidth;
      const h = el?.offsetHeight ?? 400;
      const next = {
        x: clamp(dragRef.current.origX + dx, 8, window.innerWidth - w - 8),
        y: clamp(
          dragRef.current.origY + dy,
          8,
          Math.max(8, window.innerHeight - Math.min(h, 120))
        ),
      };
      setPos(next);
    },
    [fallbackWidth]
  );

  const onPointerUp = useCallback(
    (e: ReactPointerEvent<HTMLDivElement>) => {
      if (!dragRef.current) return;
      dragRef.current = null;
      try {
        e.currentTarget.releasePointerCapture(e.pointerId);
      } catch {
        // ignore
      }
      setPos((current) => {
        persist(current);
        return current;
      });
    },
    [persist]
  );

  useEffect(() => {
    const onResize = () => {
      setPos((current) => {
        const el = panelRef.current;
        const w = el?.offsetWidth ?? fallbackWidth;
        const next = {
          x: clamp(current.x, 8, window.innerWidth - w - 8),
          y: clamp(current.y, 8, window.innerHeight - 80),
        };
        persist(next);
        return next;
      });
    };
    window.addEventListener("resize", onResize);
    return () => window.removeEventListener("resize", onResize);
  }, [persist, fallbackWidth]);

  return {
    panelRef,
    style: { left: pos.x, top: pos.y, right: "auto", bottom: "auto" } as const,
    dragHandleProps: {
      onPointerDown,
      onPointerMove,
      onPointerUp,
      onPointerCancel: onPointerUp,
    },
  };
}

export function instructionPanelDefaultPos(): Pos {
  if (typeof window === "undefined") return { x: 24, y: 120 };
  const width = Math.min(window.innerWidth - 32, 340);
  return {
    x: Math.max(16, window.innerWidth - width - 24),
    y: Math.max(16, window.innerHeight - 520),
  };
}

export function walkthroughPanelDefaultPos(): Pos {
  if (typeof window === "undefined") return { x: 16, y: 120 };
  return {
    x: 16,
    y: Math.max(80, Math.round(window.innerHeight / 2 - 200)),
  };
}
