"use client";

import { useEffect, useRef } from "react";

import { usePathname, useRouter, useSearchParams } from "next/navigation";

import { getModules } from "@/features/metadata/api";

import { useDemo } from "./DemoProvider";
import {
  resolveTutorialWorkspace,
  workspacePath,
  workspaceTabForStep,
} from "./resolveTutorialWorkspace";

/**
 * Auto-navigates view/create mentored steps to the right screen so the user
 * can see them — advancing still requires Continue (view) or the create action.
 */
export default function DemoGuidedNav() {
  const demo = useDemo();
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const lastKey = useRef<string | null>(null);

  useEffect(() => {
    if (!demo || demo.mode !== "running" || !demo.currentStep) return;

    const stepKey = demo.currentStep.step_key;

    // Workspace tabs (overview / notes / timeline)
    const tab = workspaceTabForStep(stepKey);
    if (tab) {
      const onRecord =
        /^\/m\/[^/]+\/[^/]+$/.test(pathname) &&
        searchParams.get("tab") === tab;

      if (lastKey.current === stepKey && onRecord) return;

      let cancelled = false;
      (async () => {
        try {
          const target = await resolveTutorialWorkspace(tab);
          if (cancelled || !target) return;
          lastKey.current = stepKey;
          router.replace(workspacePath(target));
        } catch {
          // Instruction panel still explains the manual path.
        }
      })();

      return () => {
        cancelled = true;
      };
    }

    // Product demo showcase table
    if (stepKey === "product_demo_module") {
      if (lastKey.current === stepKey && pathname.startsWith("/m/")) {
        return;
      }
      let cancelled = false;
      (async () => {
        try {
          const modules = await getModules();
          const mod = modules.find((m) => m.api_name === "product_demo");
          if (cancelled) return;
          lastKey.current = stepKey;
          if (mod) router.replace(`/m/${mod.api_name}`);
          else router.replace("/dashboard");
        } catch {
          // ignore
        }
      })();
      return () => {
        cancelled = true;
      };
    }

    lastKey.current = null;
  }, [
    demo?.mode,
    demo?.currentStep?.step_key,
    pathname,
    searchParams,
    router,
    demo,
  ]);

  return null;
}
