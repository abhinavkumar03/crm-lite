"use client";

import { useEffect, useRef } from "react";

import { useDemo } from "@/features/demo/DemoProvider";
import { useTour } from "@/features/tour/TourProvider";

/**
 * Ensures Tour and Interactive Demo never fight for the screen.
 * Demo UI takes priority; Tour is paused (not skipped) while Demo is open.
 */
export default function OnboardingCoordination() {
  const tour = useTour();
  const demo = useDemo();
  const pausedRef = useRef(false);

  useEffect(() => {
    if (!tour || !demo) return;

    const demoOpen = ["catalog", "resume", "running", "cleanup"].includes(
      demo.mode
    );

    if (demoOpen && tour.active) {
      tour.pause();
      pausedRef.current = true;
    }
  }, [demo, tour]);

  return null;
}
