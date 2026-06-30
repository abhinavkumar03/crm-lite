import api from "@/services/api";
import { uploadFile } from "@/features/uploads/api";

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
    const uploaded = await uploadFile(file);

    const payload = {
        file_name: file.name,
        file_url: uploaded.url,
        public_id: uploaded.public_id,
        resource_type: uploaded.resource_type,
        file_size: uploaded.bytes,
    };

    const res = await api.post(
        `/attachments/lead/${leadId}`,
        payload
    );

    return res.data.data;
}

export async function deleteLeadAttachment(
    attachmentId: string
) {
    await api.delete(
        `/attachments/${attachmentId}`
    );
}