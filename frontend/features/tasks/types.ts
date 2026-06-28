export interface Task {
    id: string;
    title: string;
    description: string;
    status: "PENDING" | "IN_PROGRESS" | "COMPLETED";
    due_date?: string;
    lead_id?: string;
    contact_id?: string;
}

export interface CreateTaskPayload {
    title: string;
    description: string;
    status: "PENDING" | "IN_PROGRESS" | "COMPLETED";
    due_date?: string;
    lead_id?: string;
    contact_id?: string;
}