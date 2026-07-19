"use client";

import { useEffect, useState } from "react";

/**
 * Measures a highlighted element, retrying across frames so it works while
 * a route transition is still painting.
 */
export function useTargetRect(
  selector: string | undefined | null,
  active: boolean,
  remountKey: string | number,
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
  }, [selector, active, remountKey, pathname]);

  return rect;
}
