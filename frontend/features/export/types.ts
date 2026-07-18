// Types mirroring the backend export engine (Phase 13).

export type ExportStatus = "pending" | "processing" | "completed" | "failed";

export type ExportFormat = "csv" | "xlsx";

// A filter clause reuses the record runtime's shape (Phase 10).
export interface ExportFilter {
  field: string;
  operator: string;
  value: unknown;
}

// ExportSpec is the shared config for sync downloads, async jobs and templates.
export interface ExportSpec {
  format: ExportFormat;
  columns?: string[];
  filters?: ExportFilter[];
  search?: string;
  sort?: string;
  order?: string;
  expand?: boolean;
}

export interface ExportJob {
  id: string;
  module_id: string;
  filename: string;
  format: ExportFormat;
  status: ExportStatus;
  columns: string[];
  row_count: number;
  byte_size: number;
  error: string | null;
  created_by: string | null;
  started_at: string | null;
  finished_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface ExportListParams {
  page?: number;
  page_size?: number;
  status?: ExportStatus;
}

export interface ExportListResult {
  exports: ExportJob[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

export interface TemplateSort {
  field: string;
  order: string;
}

export interface ExportTemplate {
  id: string;
  module_id: string;
  name: string;
  format: ExportFormat;
  columns: string[];
  filters: ExportFilter[];
  sort: TemplateSort;
  created_by: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateTemplatePayload {
  name: string;
  format?: ExportFormat;
  columns?: string[];
  filters?: ExportFilter[];
  sort?: TemplateSort;
}
