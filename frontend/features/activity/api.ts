import api from "@/services/api";

export async function getLeadActivities(
  leadId: string
) {
  const res = await api.get(
    `/activities/lead/${leadId}`
  );

  return res.data.data;
}