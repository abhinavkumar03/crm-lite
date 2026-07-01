import api from "@/services/api";

import { CreateCallPayload } from "./types";

export async function getLeadCalls(
  leadId: string
) {
  const res = await api.get(
    `/calllogs/lead/${leadId}`
  );

  return res.data.data;
}

export async function createLeadCall(
  leadId: string,
  payload: CreateCallPayload
) {
  const res = await api.post(
    `/calllogs/lead/${leadId}`,
    payload
  );

  return res.data.data;
}

export async function updateLeadCall(
  callId: string,
  payload: CreateCallPayload
) {
  const res = await api.put(
    `/calllogs/${callId}`,
    payload
  );

  return res.data.data;
}

export async function deleteLeadCall(
  callId: string
) {
  await api.delete(
    `/calllogs/${callId}`
  );
}