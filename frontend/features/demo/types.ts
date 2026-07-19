export type DemoStepStatus =
  | "locked"
  | "active"
  | "completed"
  | "skipped"
  | "failed";

export type DemoSessionStatus =
  | "active"
  | "completed"
  | "abandoned"
  | "cleaned";

export type DemoWorkflowInfo = {
  workflow_key: string;
  name: string;
  description: string;
  version: number;
  duration_min: number;
  features: string[];
};

export type DemoStep = {
  step_key: string;
  sort_order: number;
  title: string;
  description: string;
  why_it_matters: string;
  how_it_works: string;
  expected_result: string;
  route?: string | null;
  target_selector?: string | null;
  action_label?: string | null;
  validator_key: string;
  validator_params: Record<string, unknown>;
  is_skippable: boolean;
  status: DemoStepStatus;
  required_action?: string | null;
  success_event?: string | null;
  failure_message?: string | null;
  hint?: string | null;
  max_attempts?: number | null;
  allow_selectors?: string[] | null;
  placement?: "top" | "bottom" | "left" | "right" | "center" | null;
};

export type DemoSession = {
  id: string;
  workflow_key: string;
  workflow_version: number;
  sandbox_organization_id?: string | null;
  status: DemoSessionStatus;
  current_step_key?: string | null;
  started_at: string;
  completed_at?: string | null;
  stats: Record<string, unknown>;
  steps: DemoStep[];
  progress_percent: number;
};

export type ValidateStepResult = {
  ok: boolean;
  message: string;
  session?: DemoSession;
};
