import api from "@/services/api";

export async function getLeadAttachments(
  leadId: string
) {
  const res = await api.get(
    `/attachments/lead/${leadId}`
  );

  return res.data.data;
}

export async function uploadLeadAttachment(
  leadId: string,
  file: File
) {
  const formData =
    new FormData();

  formData.append(
    "file",
    file
  );

  const res = await api.post(
    `/attachments/lead/${leadId}`,
    formData,
    {
      headers: {
        "Content-Type":
          "multipart/form-data",
      },
    }
  );

  return res.data.data;
}

export async function deleteLeadAttachment(
  leadId: string,
  attachmentId: string
) {
  await api.delete(
    `/attachments/${attachmentId}`
  );
}