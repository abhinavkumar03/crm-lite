export interface Activity {
  id: string;

  action: string;

  description: string;

  performed_by: string;

  metadata: string | null;

  created_at: string;
}