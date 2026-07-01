export interface CallLog {
  id: string;

  entity_type:
    | "LEAD"
    | "CONTACT"
    | "TASK";

  entity_id: string;

  direction: "INCOMING" | "OUTGOING";

  status:
    | "COMPLETED"
    | "MISSED"
    | "NO_ANSWER"
    | "VOICEMAIL";

  duration_seconds: number;

  summary: string;

  follow_up_at?: string;

  created_by: string;

  created_at: string;
}

export interface CreateCallPayload {
  direction: "INCOMING" | "OUTGOING";

  status:
    | "COMPLETED"
    | "MISSED"
    | "NO_ANSWER"
    | "VOICEMAIL";

  duration_seconds: number;

  summary: string;

  follow_up_at?: string;
}