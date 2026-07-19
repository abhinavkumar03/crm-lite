import api from "@/services/api";

import {
  DemoSession,
  DemoWorkflowInfo,
  ValidateStepResult,
} from "./types";

export async function getDemoCatalog(): Promise<DemoWorkflowInfo> {
  const res = await api.get("/demo/catalog");
  return res.data.data;
}

export async function getActiveDemoSession(): Promise<DemoSession | null> {
  const res = await api.get("/demo/session");
  return res.data.data ?? null;
}

export async function startDemo(): Promise<DemoSession> {
  const res = await api.post("/demo/start");
  return res.data.data;
}

export async function restartDemo(): Promise<DemoSession> {
  const res = await api.post("/demo/restart");
  return res.data.data;
}

export async function validateDemoStep(
  sessionId: string,
  stepKey: string,
  route?: string
): Promise<ValidateStepResult> {
  try {
    const res = await api.post(`/demo/sessions/${sessionId}/validate`, {
      step_key: stepKey,
      route: route ?? "",
    });
    return res.data.data;
  } catch (err: unknown) {
    const axiosErr = err as {
      response?: { data?: { message?: string; data?: ValidateStepResult } };
    };
    const data = axiosErr.response?.data?.data;
    if (data) return data;
    return {
      ok: false,
      message:
        axiosErr.response?.data?.message ?? "Validation failed — try again",
    };
  }
}

export async function skipDemoStep(
  sessionId: string,
  stepKey: string
): Promise<DemoSession> {
  const res = await api.post(`/demo/sessions/${sessionId}/skip`, {
    step_key: stepKey,
  });
  return res.data.data;
}

export async function completeDemo(sessionId: string): Promise<DemoSession> {
  const res = await api.post(`/demo/sessions/${sessionId}/complete`);
  return res.data.data;
}

export async function cleanupDemo(
  sessionId: string,
  keepData: boolean
): Promise<DemoSession> {
  const res = await api.post(`/demo/sessions/${sessionId}/cleanup`, {
    keep_data: keepData,
  });
  return res.data.data;
}

export async function logDemoEvent(
  sessionId: string,
  eventType: string,
  payload?: Record<string, unknown>
): Promise<void> {
  await api.post(`/demo/sessions/${sessionId}/events`, {
    event_type: eventType,
    payload: payload ?? {},
  });
}
