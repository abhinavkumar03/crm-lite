"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  ReactNode,
} from "react";

import { usePathname, useRouter } from "next/navigation";

import { useAuth } from "@/context/AuthContext";

import { getTourProgress, restartTour, saveTourProgress } from "./api";
import { APP_TOUR_KEY, APP_TOUR_STEPS, TourStep } from "./steps";
import { UpdateProgressPayload } from "./types";

interface TourContextValue {
  active: boolean;
  steps: TourStep[];
  stepIndex: number;
  currentStep: TourStep | null;
  totalSteps: number;
  next: () => void;
  back: () => void;
  goTo: (index: number) => void;
  skip: () => void;
  finish: () => void;
  start: () => void;
  restart: () => void;
}

const TourContext = createContext<TourContextValue | null>(null);

export function TourProvider({ children }: { children: ReactNode }) {
  const auth = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  const steps = APP_TOUR_STEPS;
  const totalSteps = steps.length;

  const [active, setActive] = useState(false);
  const [stepIndex, setStepIndex] = useState(0);
  const loadedRef = useRef(false);

  // Fire-and-forget persistence: onboarding is non-critical, so a failed save
  // never blocks the UI.
  const persist = useCallback((payload: UpdateProgressPayload) => {
    saveTourProgress({ tour_key: APP_TOUR_KEY, ...payload }).catch(() => {});
  }, []);

  const stepKeysUpTo = useCallback(
    (index: number) => steps.slice(0, index + 1).map((s) => s.key),
    [steps]
  );

  // Load persisted progress once, after authentication. Auto-resume only for
  // users who have not completed or skipped the tour.
  useEffect(() => {
    if (!auth.token || loadedRef.current) return;
    loadedRef.current = true;

    (async () => {
      try {
        const progress = await getTourProgress(APP_TOUR_KEY);
        if (progress.status === "active") {
          const idx = Math.min(
            Math.max(progress.current_step, 0),
            totalSteps - 1
          );
          setStepIndex(idx);
          setActive(true);
        }
      } catch {
        // Ignore: the tour is a progressive enhancement.
      }
    })();
  }, [auth.token, totalSteps]);

  // Navigate to a step's route (if any) when it becomes active.
  useEffect(() => {
    if (!active) return;
    const step = steps[stepIndex];
    if (step?.route && step.route !== pathname) {
      router.push(step.route);
    }
  }, [active, stepIndex, steps, pathname, router]);

  const finish = useCallback(() => {
    setActive(false);
    persist({
      status: "completed",
      current_step: totalSteps - 1,
      completed_steps: steps.map((s) => s.key),
    });
  }, [persist, steps, totalSteps]);

  const next = useCallback(() => {
    if (stepIndex >= totalSteps - 1) {
      finish();
      return;
    }
    const nextIdx = stepIndex + 1;
    setStepIndex(nextIdx);
    persist({
      status: "active",
      current_step: nextIdx,
      completed_steps: stepKeysUpTo(nextIdx),
    });
  }, [stepIndex, totalSteps, finish, persist, stepKeysUpTo]);

  const back = useCallback(() => {
    if (stepIndex <= 0) return;
    const prevIdx = stepIndex - 1;
    setStepIndex(prevIdx);
    persist({ status: "active", current_step: prevIdx });
  }, [stepIndex, persist]);

  const goTo = useCallback(
    (index: number) => {
      const clamped = Math.min(Math.max(index, 0), totalSteps - 1);
      setStepIndex(clamped);
      persist({
        status: "active",
        current_step: clamped,
        completed_steps: stepKeysUpTo(clamped),
      });
    },
    [totalSteps, persist, stepKeysUpTo]
  );

  const skip = useCallback(() => {
    setActive(false);
    persist({ status: "skipped", current_step: stepIndex });
  }, [stepIndex, persist]);

  const start = useCallback(() => {
    setActive(true);
    persist({ status: "active", current_step: stepIndex });
  }, [stepIndex, persist]);

  const restart = useCallback(() => {
    setStepIndex(0);
    setActive(true);
    restartTour(APP_TOUR_KEY).catch(() => {});
  }, []);

  const value = useMemo<TourContextValue>(
    () => ({
      active,
      steps,
      stepIndex,
      currentStep: steps[stepIndex] ?? null,
      totalSteps,
      next,
      back,
      goTo,
      skip,
      finish,
      start,
      restart,
    }),
    [
      active,
      steps,
      stepIndex,
      totalSteps,
      next,
      back,
      goTo,
      skip,
      finish,
      start,
      restart,
    ]
  );

  return (
    <TourContext.Provider value={value}>{children}</TourContext.Provider>
  );
}

// useTour returns the tour controller, or null when rendered outside a
// TourProvider (e.g. on the public marketing pages), so shared components can
// opt in safely.
export function useTour(): TourContextValue | null {
  return useContext(TourContext);
}
