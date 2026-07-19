import type { GuidedPlacement } from "./types";

export const GUIDED_CARD_WIDTH = 340;
export const GUIDED_CARD_HEIGHT = 210;
export const GUIDED_MARGIN = 16;
export const SPOTLIGHT_PADDING = 8;

export function tooltipPosition(
  rect: DOMRect | null,
  placement: GuidedPlacement | undefined,
  vp: { w: number; h: number },
  cardWidth = GUIDED_CARD_WIDTH,
  cardHeight = GUIDED_CARD_HEIGHT
): { top: number; left: number } {
  if (!rect || placement === "center") {
    return {
      top: vp.h / 2 - cardHeight / 2,
      left: vp.w / 2 - cardWidth / 2,
    };
  }

  let top: number;
  let left: number;

  switch (placement) {
    case "right":
      top = rect.top;
      left = rect.right + GUIDED_MARGIN;
      break;
    case "left":
      top = rect.top;
      left = rect.left - cardWidth - GUIDED_MARGIN;
      break;
    case "top":
      top = rect.top - cardHeight - GUIDED_MARGIN;
      left = rect.left;
      break;
    case "bottom":
    default:
      top = rect.bottom + GUIDED_MARGIN;
      left = rect.left;
      break;
  }

  left = Math.max(GUIDED_MARGIN, Math.min(left, vp.w - cardWidth - GUIDED_MARGIN));
  top = Math.max(GUIDED_MARGIN, Math.min(top, vp.h - cardHeight - GUIDED_MARGIN));
  return { top, left };
}
