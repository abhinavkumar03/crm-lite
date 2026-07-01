export interface Lead {
    id: string;
    name: string;
    email: string;
    company: string;
    status: string;
}

export interface Task {
    id: string;
    title: string;
    status: string;
    due_date?: string;
}

export interface DashboardResponse {
    total_leads: number;

    new_leads: number;

    contacted_leads: number;

    qualified_leads: number;

    won_leads: number;

    lost_leads: number;

    total_contacts: number;

    total_tasks: number;

    pending_tasks: number;

    in_progress_tasks: number;

    completed_tasks: number;

    recent_leads: Lead[];

    upcoming_tasks: Task[];

    recent_activities: DashboardActivity[];
}

export interface DashboardActivity {
  id: string;

  entity_type:
    | "LEAD"
    | "CONTACT"
    | "TASK";

  entity_id: string;

  action: string;

  description: string;

  metadata: string | null;

  created_at: string;
}