package catalog

// BuiltinTemplate is a starter workflow definition (org-agnostic).
type BuiltinTemplate struct {
	Key           string
	Name          string
	Description   string
	ModuleAPIName string
	Category      string // sales | nurture | tasks | lifecycle
	Definition    map[string]any
}

// BuiltinTemplates returns the rich gallery used by seed + EnsureBuiltins.
func BuiltinTemplates() []BuiltinTemplate {
	and := "and"
	return []BuiltinTemplate{
		{
			Key: "lead_qualification", Name: "Lead Qualification", ModuleAPIName: "lead", Category: "sales",
			Description: "When status becomes Qualified: log activity and send a welcome email.",
			Definition: map[string]any{
				"triggers": []any{map[string]any{"type": "field_updated", "config": map[string]any{"field_api_name": "status"}}},
				"conditions": map[string]any{
					"node_type": "group", "logic": and,
					"children": []any{
						map[string]any{"node_type": "predicate", "field_api_name": "status", "operator": "eq", "value": "Qualified"},
					},
				},
				"actions": []any{
					map[string]any{"type": "create_activity", "config": map[string]any{"description": "Lead qualified — follow-up started"}},
					map[string]any{"type": "send_email", "config": map[string]any{"subject": "Welcome!", "body": "Hello {{lead.name}}, welcome aboard."}},
				},
			},
		},
		{
			Key: "lost_lead", Name: "Lost Lead", ModuleAPIName: "lead", Category: "sales",
			Description: "When a lead is marked Lost: create a timeline activity for managers.",
			Definition: map[string]any{
				"triggers": []any{map[string]any{"type": "field_updated", "config": map[string]any{"field_api_name": "status"}}},
				"conditions": map[string]any{
					"node_type": "group", "logic": and,
					"children": []any{
						map[string]any{"node_type": "predicate", "field_api_name": "status", "operator": "eq", "value": "Lost"},
					},
				},
				"actions": []any{
					map[string]any{"type": "create_activity", "config": map[string]any{"description": "Lead marked as Lost"}},
					map[string]any{"type": "create_note", "config": map[string]any{"body": "Lost reason review needed."}},
				},
			},
		},
		{
			Key: "new_website_lead", Name: "New Website Lead", ModuleAPIName: "lead", Category: "sales",
			Description: "On website lead create: activity + WhatsApp acknowledgment.",
			Definition: map[string]any{
				"triggers": []any{map[string]any{"type": "record_created", "config": map[string]any{}}},
				"conditions": map[string]any{
					"node_type": "group", "logic": and,
					"children": []any{
						map[string]any{"node_type": "predicate", "field_api_name": "source", "operator": "eq", "value": "Website"},
					},
				},
				"actions": []any{
					map[string]any{"type": "create_activity", "config": map[string]any{"description": "Website lead captured"}},
					map[string]any{"type": "send_whatsapp", "config": map[string]any{"body": "Thanks {{lead.name}} — we received your enquiry."}},
				},
			},
		},
		{
			Key: "high_value_lead", Name: "High Value Lead", ModuleAPIName: "lead", Category: "sales",
			Description: "When a high-value lead is created, flag with a note and activity.",
			Definition: map[string]any{
				"triggers": []any{map[string]any{"type": "record_created", "config": map[string]any{}}},
				"actions": []any{
					map[string]any{"type": "create_activity", "config": map[string]any{"description": "High value lead — prioritize"}},
					map[string]any{"type": "create_note", "config": map[string]any{"title": "Priority", "body": "High value lead — assign senior AE."}},
				},
			},
		},
		{
			Key: "lead_assigned", Name: "Lead Assigned", ModuleAPIName: "lead", Category: "sales",
			Description: "When assigned_to changes: notify via activity and optional email.",
			Definition: map[string]any{
				"triggers": []any{map[string]any{"type": "field_updated", "config": map[string]any{"field_api_name": "assigned_to"}}},
				"actions": []any{
					map[string]any{"type": "create_activity", "config": map[string]any{"description": "Lead reassigned"}},
					map[string]any{"type": "send_email", "config": map[string]any{"subject": "Lead assigned to you", "body": "A lead was assigned. Review {{lead.name}}."}},
				},
			},
		},
		{
			Key: "new_contact", Name: "New Contact", ModuleAPIName: "contact", Category: "nurture",
			Description: "Send a welcome email and attach a note when a contact is created.",
			Definition: map[string]any{
				"triggers": []any{map[string]any{"type": "record_created", "config": map[string]any{}}},
				"actions": []any{
					map[string]any{"type": "send_email", "config": map[string]any{"subject": "Welcome", "body": "Hello {{contact.name}}, nice to meet you."}},
					map[string]any{"type": "create_note", "config": map[string]any{"body": "Welcome automation ran."}},
				},
			},
		},
		{
			Key: "contact_updated", Name: "Contact Updated", ModuleAPIName: "contact", Category: "nurture",
			Description: "Log an activity whenever a contact record is updated.",
			Definition: map[string]any{
				"triggers": []any{map[string]any{"type": "record_updated", "config": map[string]any{}}},
				"actions": []any{
					map[string]any{"type": "create_activity", "config": map[string]any{"description": "Contact profile updated"}},
				},
			},
		},
		{
			Key: "birthday", Name: "Birthday Greeting", ModuleAPIName: "contact", Category: "nurture",
			Description: "Date-based: send a birthday greeting email when birthday is today.",
			Definition: map[string]any{
				"triggers": []any{map[string]any{"type": "date_based", "config": map[string]any{"field_api_name": "birthday", "offset_days": 0}}},
				"actions": []any{
					map[string]any{"type": "send_email", "config": map[string]any{"subject": "Happy Birthday!", "body": "Happy Birthday {{contact.name}}! From everyone at {{workspace.name}}."}},
					map[string]any{"type": "create_activity", "config": map[string]any{"description": "Birthday greeting sent"}},
				},
			},
		},
		{
			Key: "anniversary", Name: "Anniversary", ModuleAPIName: "contact", Category: "nurture",
			Description: "Date-based: celebrate anniversary field with a note + email.",
			Definition: map[string]any{
				"triggers": []any{map[string]any{"type": "date_based", "config": map[string]any{"field_api_name": "anniversary", "offset_days": 0}}},
				"actions": []any{
					map[string]any{"type": "send_email", "config": map[string]any{"subject": "Happy Anniversary", "body": "Congratulations {{contact.name}} on your anniversary!"}},
					map[string]any{"type": "create_note", "config": map[string]any{"body": "Anniversary outreach sent."}},
				},
			},
		},
		{
			Key: "task_created", Name: "Task Created", ModuleAPIName: "task", Category: "tasks",
			Description: "When a task is created, write a timeline activity.",
			Definition: map[string]any{
				"triggers": []any{map[string]any{"type": "record_created", "config": map[string]any{}}},
				"actions": []any{
					map[string]any{"type": "create_activity", "config": map[string]any{"description": "Task created automation"}},
				},
			},
		},
		{
			Key: "task_completed", Name: "Task Completed", ModuleAPIName: "task", Category: "tasks",
			Description: "When task status becomes Completed, log activity.",
			Definition: map[string]any{
				"triggers": []any{map[string]any{"type": "field_updated", "config": map[string]any{"field_api_name": "status"}}},
				"conditions": map[string]any{
					"node_type": "group", "logic": and,
					"children": []any{
						map[string]any{"node_type": "predicate", "field_api_name": "status", "operator": "eq", "value": "Completed"},
					},
				},
				"actions": []any{
					map[string]any{"type": "create_activity", "config": map[string]any{"description": "Task completed"}},
				},
			},
		},
		{
			Key: "task_overdue", Name: "Task Overdue", ModuleAPIName: "task", Category: "tasks",
			Description: "Date-based: WhatsApp reminder when due_date is today.",
			Definition: map[string]any{
				"triggers": []any{map[string]any{"type": "date_based", "config": map[string]any{"field_api_name": "due_date", "offset_days": 0}}},
				"actions": []any{
					map[string]any{"type": "send_whatsapp", "config": map[string]any{"body": "Reminder: task {{task.title}} is due today."}},
					map[string]any{"type": "create_activity", "config": map[string]any{"description": "Overdue reminder sent"}},
				},
			},
		},
		{
			Key: "manual_follow_up", Name: "Manual Follow-up", ModuleAPIName: "lead", Category: "lifecycle",
			Description: "Run on demand from record detail to leave a follow-up note + activity.",
			Definition: map[string]any{
				"triggers": []any{map[string]any{"type": "manual", "config": map[string]any{}}},
				"actions": []any{
					map[string]any{"type": "create_activity", "config": map[string]any{"description": "Manual follow-up workflow"}},
					map[string]any{"type": "create_note", "config": map[string]any{"body": "Manual follow-up note from workflow"}},
				},
			},
		},
	}
}
