import api from "@/services/api";

export type UploadedFile = {
  url: string;
  public_id: string;
  resource_type: string;
  bytes: number;
  format: string;
};

export async function uploadFile(
  file: File
): Promise<UploadedFile> {
  const formData = new FormData();

  formData.append("file", file);

  const res = await api.post(
    "/uploads",
    formData
  );

  return res.data.data;
}