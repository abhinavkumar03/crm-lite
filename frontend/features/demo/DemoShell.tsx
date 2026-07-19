"use client";

import { Suspense } from "react";

import OnboardingCoordination from "@/features/guided/OnboardingCoordination";
import { useTour } from "@/features/tour/TourProvider";

import CleanupDialog from "./CleanupDialog";
import DemoCatalogModal from "./DemoCatalogModal";
import DemoGuideOverlay from "./DemoGuideOverlay";
import DemoGuidedNav from "./DemoGuidedNav";
import DemoLauncher from "./DemoLauncher";
import InstructionPanel from "./InstructionPanel";
import ResumeDialog from "./ResumeDialog";
import TutorialSidebar from "./TutorialSidebar";
import { useDemo } from "./DemoProvider";

/** Mounts all demo engine UI surfaces inside the authenticated shell. */
export default function DemoShell() {
  const tour = useTour();
  const demo = useDemo();
  const tourActive = tour?.active ?? false;

  // Launcher stays available; instructional surfaces hide while Tour owns the screen.
  const showPanels = !tourActive;

  return (
    <>
      <OnboardingCoordination />
      <DemoLauncher />
      {showPanels && (
        <>
          <DemoCatalogModal />
          <ResumeDialog />
          <CleanupDialog />
          {demo?.mode === "running" && (
            <>
              <Suspense fallback={null}>
                <DemoGuidedNav />
              </Suspense>
              <DemoGuideOverlay />
            </>
          )}
          <InstructionPanel />
          <TutorialSidebar />
        </>
      )}
    </>
  );
}
