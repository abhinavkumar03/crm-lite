import api from "@/services/api";

import { TourProgress, UpdateProgressPayload } from "./types";

// getTourProgress returns the current user's progress. A brand-new user gets a
// synthesized "active" default from the backend.
export async function getTourProgress(tourKey?: string): Promise<TourProgress> {
  const res = await api.get("/tour", {
    params: tourKey ? { key: tourKey } : undefined,
  });
  return res.data.data;
}

// saveTourProgress is the single write path used to advance, complete or skip.
export async function saveTourProgress(
  payload: UpdateProgressPayload
): Promise<TourProgress> {
  const res = await api.put("/tour", payload);
  return res.data.data;
}

// restartTour resets progress back to the first step.
export async function restartTour(tourKey?: string): Promise<TourProgress> {
  const res = await api.post("/tour/restart", tourKey ? { tour_key: tourKey } : {});
  return res.data.data;
}
