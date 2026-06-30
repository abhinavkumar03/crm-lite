export interface Note {
  id: string;

  note: string;

  created_at: string;

  updated_at: string;

  user: {
    id: string;

    name: string;
  };
}

export interface CreateNotePayload {
  note: string;
}

export interface UpdateNotePayload {
  note: string;
}