package dto

import (
	leadDto "github.com/abhinavkumar03/crm-lite/backend/internal/lead/dto"
	taskDto "github.com/abhinavkumar03/crm-lite/backend/internal/task/dto"
)

type DashboardResponse struct {
	TotalLeads int64 `json:"total_leads"`

	NewLeads int64 `json:"new_leads"`

	ContactedLeads int64 `json:"contacted_leads"`

	QualifiedLeads int64 `json:"qualified_leads"`

	WonLeads int64 `json:"won_leads"`

	LostLeads int64 `json:"lost_leads"`

	TotalContacts int64 `json:"total_contacts"`

	TotalTasks int64 `json:"total_tasks"`

	PendingTasks int64 `json:"pending_tasks"`

	InProgressTasks int64 `json:"in_progress_tasks"`

	CompletedTasks int64 `json:"completed_tasks"`

	RecentLeads []leadDto.LeadResponse `json:"recent_leads"`

	UpcomingTasks []taskDto.TaskResponse `json:"upcoming_tasks"`

	RecentActivities []RecentActivityResponse `json:"recent_activities"`
}
