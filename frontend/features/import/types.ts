// Types mirroring the backend import engine (Phase 12).

export type ImportStatus = "pending" | "processing" | "completed" | "failed";

export interface RowError {
  row: number;
  field?: string;
  message: string;
}

// AnalyzeResult drives the mapping UI: detected columns, a preview sample and
// an auto-suggested column (header) -> field api_name mapping.
export interface AnalyzeResult {
  headers: string[];
  sample_rows: Record<string, string>[];
  suggested_mapping: Record<string, string>;
  row_count: number;
}

export interface ImportJob {
  id: string;
  module_id: string;
  filename: string;
  status: ImportStatus;
  mapping: Record<string, string>;
  total_rows: number;
  processed_rows: number;
  success_rows: number;
  error_rows: number;
  errors: RowError[];
  created_by: string | null;
  started_at: string | null;
  finished_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface ImportListParams {
  page?: number;
  page_size?: number;
  status?: ImportStatus;
}

export interface ImportListResult {
  imports: ImportJob[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}
