package dto

type ModuleCount struct {
	ModuleID     string `json:"module_id"`
	APIName      string `json:"api_name"`
	PluralLabel  string `json:"plural_label"`
	Icon         string `json:"icon,omitempty"`
	Color        string `json:"color,omitempty"`
	RecordCount  int64  `json:"record_count"`
}

type RecentRecord struct {
	ID          string `json:"id"`
	ModuleID    string `json:"module_id"`
	ModuleLabel string `json:"module_label"`
	APIName     string `json:"api_name"`
	Title       string `json:"title"`
	CreatedAt   string `json:"created_at"`
}

type DashboardResponse struct {
	TotalModules int64           `json:"total_modules"`
	TotalRecords int64           `json:"total_records"`
	ModuleCounts []ModuleCount   `json:"module_counts"`
	RecentRecords []RecentRecord `json:"recent_records"`
}
