import api from "@/services/api";
import {
  CreateNotePayload,
  UpdateNotePayload,
} from "./types";

/*
GET
/api/v1/notes/lead/:leadId
*/
export async function getLeadNotes(
  leadId: string
) {
  const res = await api.get(
    `/notes/lead/${leadId}`
  );

  return res.data.data;
}

/*
POST
/api/v1/notes/lead/:leadId
*/
export async function createLeadNote(
  leadId: string,
  payload: CreateNotePayload
) {
  const res = await api.post(
    `/notes/lead/${leadId}`,
    payload
  );

  return res.data.data;
}

/*
PUT
/api/v1/notes/:noteId
*/
export async function updateLeadNote(
  noteId: string,
  payload: UpdateNotePayload
) {
  const res = await api.put(
    `/notes/${noteId}`,
    payload
  );

  return res.data.data;
}

/*
DELETE
/api/v1/notes/:noteId
*/
export async function deleteLeadNote(
  noteId: string
) {
  await api.delete(
    `/notes/${noteId}`
  );
}