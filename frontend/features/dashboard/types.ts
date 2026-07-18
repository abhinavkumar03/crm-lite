export interface ModuleCount {
  module_id: string;
  api_name: string;
  plural_label: string;
  icon?: string;
  color?: string;
  record_count: number;
}

export interface RecentRecord {
  id: string;
  module_id: string;
  module_label: string;
  api_name: string;
  title: string;
  created_at: string;
}

export interface DashboardResponse {
  total_modules: number;
  total_records: number;
  module_counts: ModuleCount[];
  recent_records: RecentRecord[];
}
