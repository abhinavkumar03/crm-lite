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
	TotalModules         int64                    `json:"total_modules"`
	TotalRecords         int64                    `json:"total_records"`
	ModuleCounts         []ModuleCount            `json:"module_counts"`
	RecentRecords        []RecentRecord           `json:"recent_records"`
	EmailsSentToday      int64                    `json:"emails_sent_today"`
	WhatsAppSentToday    int64                    `json:"whatsapp_sent_today"`
	FailedNotifications  int64                    `json:"failed_notifications"`
	ScheduledNotifications int64                  `json:"scheduled_notifications"`
	ActiveWorkflows      int64                    `json:"active_workflows"`
	DisabledWorkflows    int64                    `json:"disabled_workflows"`
	WorkflowsExecutedToday int64                  `json:"workflows_executed_today"`
	WorkflowsFailedToday int64                    `json:"workflows_failed_today"`
	AvgWorkflowDurationMs *float64                `json:"avg_workflow_duration_ms,omitempty"`
}
