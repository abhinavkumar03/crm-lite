// Types mirroring the backend guided-tour engine (Phase 14).

export type TourStatus = "active" | "completed" | "skipped";

// TourProgress is the per-user progress persisted server-side. The step
// catalogue itself lives on the client (see steps.ts).
export interface TourProgress {
  tour_key: string;
  status: TourStatus;
  current_step: number;
  completed_steps: string[];
  started_at: string;
  completed_at: string | null;
  updated_at: string;
}

// UpdateProgressPayload is the single write path used to advance/complete/skip.
export interface UpdateProgressPayload {
  tour_key?: string;
  status?: TourStatus;
  current_step?: number;
  completed_steps?: string[];
}
