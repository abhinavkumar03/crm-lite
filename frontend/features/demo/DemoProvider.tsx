"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";

import { usePathname, useRouter } from "next/navigation";
import { toast } from "sonner";

import { useAuth } from "@/context/AuthContext";

import {
  cleanupDemo,
  completeDemo,
  getActiveDemoSession,
  getDemoCatalog,
  restartDemo,
  skipDemoStep,
  startDemo,
  validateDemoStep,
} from "./api";
import type { DemoSession, DemoStep, DemoWorkflowInfo } from "./types";
import { isViewConfirmStep } from "./stepAdvance";

type DemoUIMode =
  | "idle"
  | "catalog"
  | "resume"
  | "running"
  | "cleanup"
  | "busy";

type DemoContextValue = {
  catalog: DemoWorkflowInfo | null;
  session: DemoSession | null;
  currentStep: DemoStep | null;
  mode: DemoUIMode;
  busy: boolean;
  lastMessage: string | null;
  /** waiting | validating | failed — UI hint for mentor mode */
  stepPhase: "idle" | "waiting" | "validating" | "failed";
  openCatalog: () => void;
  openLauncher: () => void;
  closeUI: () => void;
  start: () => Promise<void>;
  continueSession: () => void;
  restart: () => Promise<void>;
  validate: (opts?: { silent?: boolean; stepKey?: string }) => Promise<void>;
  skip: () => Promise<void>;
  finish: () => Promise<void>;
  cleanup: (keepData: boolean) => Promise<void>;
  goToStepRoute: () => void;
};

const DemoContext = createContext<DemoContextValue | null>(null);

export function DemoProvider({ children }: { children: ReactNode }) {
  const auth = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  const [catalog, setCatalog] = useState<DemoWorkflowInfo | null>(null);
  const [session, setSession] = useState<DemoSession | null>(null);
  const [mode, setMode] = useState<DemoUIMode>("idle");
  const [busy, setBusy] = useState(false);
  const [lastMessage, setLastMessage] = useState<string | null>(null);
  const [stepPhase, setStepPhase] = useState<
    "idle" | "waiting" | "validating" | "failed"
  >("waiting");
  const loadedRef = useRef(false);
  const currentStepKeyRef = useRef<string | null>(null);
  const validateInFlightRef = useRef(false);
  /** Blocks accidental double-Enter from completing the next view step instantly. */
  const confirmCooldownUntilRef = useRef(0);

  const currentStep = useMemo(() => {
    if (!session?.current_step_key) return null;
    return (
      session.steps.find((s) => s.step_key === session.current_step_key) ??
      session.steps.find((s) => s.status === "active") ??
      null
    );
  }, [session]);

  currentStepKeyRef.current = currentStep?.step_key ?? null;

  useEffect(() => {
    if (!auth.token || loadedRef.current) return;
    loadedRef.current = true;

    (async () => {
      try {
        const [cat, active] = await Promise.all([
          getDemoCatalog().catch(() => null),
          getActiveDemoSession().catch(() => null),
        ]);
        if (cat) setCatalog(cat);
        if (active && (active.status === "active" || active.status === "completed")) {
          setSession(active);
          if (active.status === "completed") {
            setMode("cleanup");
          } else {
            setMode("resume");
          }
        }
      } catch {
        // Demo is progressive enhancement.
      }
    })();
  }, [auth.token]);

  const reloadAfterTenantSwitch = useCallback(() => {
    // Sandbox org switch invalidates module lists / dashboard cache in memory.
    window.location.assign("/dashboard");
  }, []);

  const openCatalog = useCallback(() => {
    setMode(session?.status === "active" ? "resume" : "catalog");
  }, [session]);

  const openLauncher = useCallback(() => {
    if (session?.status === "completed") {
      setMode("cleanup");
      return;
    }
    if (session?.status === "active") {
      // Re-open the instruction panel (resume prompt only on first load).
      setMode("running");
      return;
    }
    setMode("catalog");
  }, [session]);

  const closeUI = useCallback(() => {
    // Dismiss overlays; launcher remains so the user can reopen.
    setMode("idle");
  }, []);

  const continueSession = useCallback(() => {
    setMode("running");
    const step =
      session?.steps.find((s) => s.step_key === session.current_step_key) ??
      session?.steps.find((s) => s.status === "active");
    if (step?.route && step.route !== pathname) {
      router.push(step.route);
    }
  }, [session, pathname, router]);

  const start = useCallback(async () => {
    setBusy(true);
    setLastMessage(null);
    try {
      const sess = await startDemo();
      setSession(sess);
      toast.success("Sandbox organization ready");
      reloadAfterTenantSwitch();
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { message?: string } } })?.response?.data
          ?.message ?? "Could not start demo";
      toast.error(msg);
      setBusy(false);
    }
  }, [reloadAfterTenantSwitch]);

  const restart = useCallback(async () => {
    setBusy(true);
    setLastMessage(null);
    try {
      await restartDemo();
      toast.success("Demo restarted with a fresh sandbox");
      reloadAfterTenantSwitch();
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { message?: string } } })?.response?.data
          ?.message ?? "Could not restart demo";
      toast.error(msg);
      setBusy(false);
    }
  }, [reloadAfterTenantSwitch]);

  const validate = useCallback(
    async (opts?: { silent?: boolean; stepKey?: string }) => {
      if (!session || !currentStep) return;
      if (validateInFlightRef.current || busy) return;

      // Ignore stale timers from a previous create step (e.g. add_note → timeline).
      if (opts?.stepKey && opts.stepKey !== currentStep.step_key) return;

      const silent = opts?.silent ?? false;

      // View / navigate steps must never silent-auto advance.
      if (silent && isViewConfirmStep(currentStep)) return;

      // After advancing, ignore rapid Enter/click that would clear the next view step.
      if (
        !silent &&
        isViewConfirmStep(currentStep) &&
        Date.now() < confirmCooldownUntilRef.current
      ) {
        return;
      }

      const stepKey = currentStep.step_key;
      validateInFlightRef.current = true;
      setBusy(true);
      setStepPhase("validating");
      if (!silent) setLastMessage(null);
      try {
        if (currentStepKeyRef.current !== stepKey) return;

        const result = await validateDemoStep(session.id, stepKey, pathname);
        if (currentStepKeyRef.current !== stepKey) return;

        setLastMessage(result.message);
        if (result.ok && result.session) {
          setSession(result.session);
          setStepPhase("waiting");
          // Give the user time to see the next view step before another confirm.
          confirmCooldownUntilRef.current = Date.now() + 800;
          if (!silent) toast.success(result.message);
          if (result.session.status === "completed") {
            setMode("cleanup");
          } else {
            const next = result.session.steps.find(
              (s) => s.step_key === result.session?.current_step_key
            );
            if (next?.route) router.push(next.route);
            setMode("running");
          }
        } else {
          setStepPhase("failed");
          if (!silent) toast.error(result.message);
        }
      } finally {
        validateInFlightRef.current = false;
        setBusy(false);
      }
    },
    [session, currentStep, pathname, router, busy]
  );

  const skip = useCallback(async () => {
    if (!session || !currentStep) return;
    if (!currentStep.is_skippable) {
      toast.error("This step requires a real action — it cannot be skipped");
      return;
    }
    setBusy(true);
    try {
      const sess = await skipDemoStep(session.id, currentStep.step_key);
      setSession(sess);
      setLastMessage("Step skipped");
      if (sess.status === "completed") {
        setMode("cleanup");
      } else {
        const next = sess.steps.find(
          (s) => s.step_key === sess.current_step_key
        );
        if (next?.route) router.push(next.route);
        setMode("running");
      }
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { message?: string } } })?.response?.data
          ?.message ?? "Could not skip step";
      toast.error(msg);
    } finally {
      setBusy(false);
    }
  }, [session, currentStep, router]);

  const finish = useCallback(async () => {
    if (!session) return;
    setBusy(true);
    try {
      if (currentStep?.validator_key === "none") {
        const result = await validateDemoStep(
          session.id,
          currentStep.step_key,
          pathname
        );
        if (result.session) setSession(result.session);
      }
      const sess = await completeDemo(session.id);
      setSession(sess);
      setMode("cleanup");
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { message?: string } } })?.response?.data
          ?.message ?? "Could not complete demo";
      toast.error(msg);
    } finally {
      setBusy(false);
    }
  }, [session, currentStep, pathname]);

  const cleanup = useCallback(
    async (keepData: boolean) => {
      if (!session) return;
      setBusy(true);
      try {
        await cleanupDemo(session.id, keepData);
        toast.success(
          keepData
            ? "Sandbox kept — you remain in the demo organization"
            : "Demo data deleted — restored previous organization"
        );
        setSession(null);
        setMode("idle");
        reloadAfterTenantSwitch();
      } catch (err: unknown) {
        const msg =
          (err as { response?: { data?: { message?: string } } })?.response
            ?.data?.message ?? "Cleanup failed";
        toast.error(msg);
        setBusy(false);
      }
    },
    [session, reloadAfterTenantSwitch]
  );

  const goToStepRoute = useCallback(() => {
    if (!currentStep) return;
    const key = currentStep.step_key;
    if (
      key === "add_note" ||
      key === "timeline" ||
      key === "record_workspace"
    ) {
      void import("./resolveTutorialWorkspace").then(
        async ({ resolveTutorialWorkspace, workspacePath, workspaceTabForStep }) => {
          const tab = workspaceTabForStep(key) ?? "overview";
          const target = await resolveTutorialWorkspace(tab);
          if (target) router.push(workspacePath(target));
          else if (currentStep.route) router.push(currentStep.route);
        }
      );
      return;
    }
    if (key === "product_demo_module") {
      void import("@/features/metadata/api").then(async ({ getModules }) => {
        const modules = await getModules();
        const mod = modules.find((m) => m.api_name === "product_demo");
        router.push(mod ? `/m/${mod.api_name}` : "/dashboard");
      });
      return;
    }
    if (currentStep.route) router.push(currentStep.route);
  }, [currentStep, router]);

  // Guided navigation: keep the user on the step's target page while running.
  // Workspace / product-demo deep links are handled by DemoGuidedNav.
  useEffect(() => {
    if (mode !== "running" || !currentStep?.route) return;
    if (
      currentStep.step_key === "add_note" ||
      currentStep.step_key === "timeline" ||
      currentStep.step_key === "record_workspace" ||
      currentStep.step_key === "product_demo_module"
    ) {
      return;
    }
    if (currentStep.route !== pathname && !pathname.startsWith(currentStep.route)) {
      router.push(currentStep.route);
    }
  }, [mode, currentStep?.route, currentStep?.step_key, pathname, router]);

  const value = useMemo<DemoContextValue>(
    () => ({
      catalog,
      session,
      currentStep,
      mode,
      busy,
      lastMessage,
      stepPhase,
      openCatalog,
      openLauncher,
      closeUI,
      start,
      continueSession,
      restart,
      validate,
      skip,
      finish,
      cleanup,
      goToStepRoute,
    }),
    [
      catalog,
      session,
      currentStep,
      mode,
      busy,
      lastMessage,
      stepPhase,
      openCatalog,
      openLauncher,
      closeUI,
      start,
      continueSession,
      restart,
      validate,
      skip,
      finish,
      cleanup,
      goToStepRoute,
    ]
  );

  return <DemoContext.Provider value={value}>{children}</DemoContext.Provider>;
}

export function useDemo(): DemoContextValue | null {
  return useContext(DemoContext);
}
