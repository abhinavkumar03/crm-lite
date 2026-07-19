import api from "@/services/api";

import { APP_TOUR_STEPS, TourPlacement, TourStep } from "./steps";

type EngineStep = {
  step_key: string;
  title: string;
  description: string;
  route?: string | null;
  target_selector?: string | null;
  placement?: string | null;
};

/**
 * Prefer metadata workflow `crm_orientation_tour` when migrations are applied;
 * otherwise fall back to the client catalogue so the tour never breaks.
 */
export async function loadOrientationSteps(): Promise<TourStep[]> {
  try {
    const res = await api.get("/demo/workflows/crm_orientation_tour");
    const steps = (res.data?.data?.steps ?? []) as EngineStep[];
    if (!steps.length) return APP_TOUR_STEPS;

    return steps.map((s) => {
      let target = s.target_selector ?? undefined;
      if (target && !target.startsWith("[") && target.includes("=")) {
        target = `[${target}]`;
      }
      const placement = (s.placement ?? "center") as TourPlacement;
      return {
        key: s.step_key,
        title: s.title,
        body: s.description,
        target,
        route: s.route ?? undefined,
        placement,
      } satisfies TourStep;
    });
  } catch {
    return APP_TOUR_STEPS;
  }
}
