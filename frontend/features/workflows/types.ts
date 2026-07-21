export type WorkflowStatus = "draft" | "active" | "disabled" | "archived";

export type TriggerInput = {
  type: string;
  config?: Record<string, unknown>;
};

export type ConditionInput = {
  node_type: "group" | "predicate";
  logic?: "and" | "or";
  field_api_name?: string;
  operator?: string;
  value?: unknown;
  children?: ConditionInput[];
};

export type ActionInput = {
  type: string;
  config?: Record<string, unknown>;
  max_retries?: number;
  continue_on_error?: boolean;
};

export type WorkflowSummary = {
  id: string;
  module_id?: string | null;
  module_api_name?: string | null;
  name: string;
  description: string;
  status: WorkflowStatus;
  priority: number;
  version: number;
  trigger_types: string[];
  action_count: number;
  created_at: string;
  updated_at: string;
};

export type WorkflowDetail = {
  id: string;
  module_id?: string | null;
  module_api_name?: string | null;
  name: string;
  description: string;
  status: WorkflowStatus;
  on_action_error: "continue" | "stop";
  priority: number;
  published_version_id?: string | null;
  draft_version_id?: string | null;
  version: number;
  triggers: Array<{ id: string; type: string; config: Record<string, unknown> }>;
  conditions?: ConditionInput & { id?: string };
  actions: Array<{
    id: string;
    sort_order: number;
    type: string;
    config: Record<string, unknown>;
    max_retries: number;
    continue_on_error?: boolean;
  }>;
  created_at: string;
  updated_at: string;
};

export type CreateWorkflowPayload = {
  name: string;
  description?: string;
  module_id?: string | null;
  on_action_error?: "continue" | "stop";
  priority?: number;
  triggers?: TriggerInput[];
  conditions?: ConditionInput | null;
  actions?: ActionInput[];
};

export type UpdateWorkflowPayload = Partial<CreateWorkflowPayload>;

export type ListWorkflowsResult = {
  items: WorkflowSummary[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
};

export type ExecutionSummary = {
  id: string;
  workflow_id: string;
  workflow_name?: string;
  module_id?: string | null;
  record_id?: string | null;
  trigger_type: string;
  status: string;
  source: string;
  depth: number;
  error_summary?: string | null;
  started_at?: string | null;
  finished_at?: string | null;
  duration_ms?: number | null;
  created_at: string;
};

export type ExecutionDetail = ExecutionSummary & {
  steps: Array<{
    id: string;
    sort_order: number;
    action_type: string;
    status: string;
    input: Record<string, unknown>;
    output: Record<string, unknown>;
    error?: string | null;
    started_at?: string | null;
    finished_at?: string | null;
  }>;
};

export type ListExecutionsResult = {
  items: ExecutionSummary[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
};

export type WorkflowTemplate = {
  id: string;
  key: string;
  name: string;
  description: string;
  module_api_name?: string | null;
  category?: string;
  definition: Record<string, unknown>;
  is_builtin: boolean;
};

export type VersionSummary = {
  id: string;
  version: number;
  state: string;
  changelog: string;
  published_at?: string | null;
  published_by?: string | null;
  created_at: string;
};

export type BuilderMetadata = {
  modules: Array<{
    id: string;
    api_name: string;
    label: string;
    fields: Array<{
      api_name: string;
      label: string;
      type: string;
      options?: string[];
      required: boolean;
    }>;
  }>;
  operators: Array<{ key: string; label: string; value_arity: string }>;
  actions: Array<{ type: string; label: string; description: string; mvp: boolean }>;
  triggers: Array<{ type: string; label: string; description: string; mvp: boolean }>;
  variables: Array<{ key: string; label: string; description: string }>;
  users: Array<{ id: string; name: string; email: string }>;
  templates: WorkflowTemplate[];
};

export type WorkflowMetrics = {
  active_workflows: number;
  disabled_workflows: number;
  draft_workflows: number;
  executed_today: number;
  failed_today: number;
  avg_duration_ms?: number | null;
};
