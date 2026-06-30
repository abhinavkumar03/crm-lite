"use client";

import {
  useEffect,
  useState,
} from "react";

import {
  Paperclip,
  Plus,
} from "lucide-react";

import { Attachment } from "../types";

import AttachmentForm from "./AttachmentForm";
import Modal from "@/components/common/Modal";

import {
  getLeadAttachments,
  uploadLeadAttachment,
  deleteLeadAttachment
} from "../api";

import AttachmentCard from "./AttachmentCard";

type Props = {
  leadId: string;
};

export default function LeadAttachments({
  leadId,
}: Props) {
  const [
    attachments,
    setAttachments,
  ] = useState<
    Attachment[]
  >([]);

  const [
    loading,
    setLoading,
  ] = useState(true);

  const [open, setOpen] =
  useState(false);

  async function load() {
    try {
      const data =
        await getLeadAttachments(
          leadId
        );

      setAttachments(data);
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete(
  attachment: Attachment
) {
  const ok =
    window.confirm(
      `Delete "${attachment.file_name}"?`
    );

  if (!ok) {
    return;
  }

  await deleteLeadAttachment(
    attachment.id
  );

  await load();
}

  useEffect(() => {
    load();
  }, []);

  if (loading) {
    return (
      <div className="rounded-3xl border border-slate-200 bg-white p-10 text-center">
        Loading attachments...
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}

      <div className="flex flex-col gap-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm md:flex-row md:items-center md:justify-between">
        <div>
          <h2 className="text-xl font-semibold">
            Attachments
          </h2>

          <p className="mt-1 text-sm text-slate-500">
            Upload proposals,
            quotations,
            contracts and
            documents.
          </p>
          <button
            onClick={() =>
                setOpen(true)
            }
            className="
            mt-8
            inline-flex
            items-center
            gap-2
            rounded-2xl
            bg-emerald-500
            px-5
            py-3
            font-medium
            text-white
            transition
            hover:bg-emerald-600
            "
            >
            <Plus size={18} />

            Upload First File
            </button>
        </div>

        <button
         onClick={() =>
            setOpen(true)
        }
          className="
          inline-flex
          items-center
          gap-2
          rounded-2xl
          bg-emerald-500
          px-5
          py-3
          font-medium
          text-white
          "
        >
          <Plus size={18} />

          Upload
        </button>
      </div>

      {/* Empty */}

      {attachments.length ===
        0 && (
        <div className="rounded-3xl border border-dashed border-slate-300 bg-white p-14 text-center">
          <div className="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-emerald-50">
            <Paperclip
              size={36}
              className="text-emerald-500"
            />
          </div>

          <h3 className="mt-6 text-xl font-semibold">
            No Attachments
          </h3>

          <p className="mx-auto mt-3 max-w-md text-slate-500">
            Upload proposals, quotations, signed contracts,
            presentations and supporting documents related
            to this lead.
            </p>
        </div>
      )}

      {/* Grid */}

      <div
        className="
        grid
        gap-6
        md:grid-cols-2
        xl:grid-cols-3
        "
      >
        {attachments.map(
          (attachment) => (
            <AttachmentCard
              key={
                attachment.id
              }
              attachment={
                attachment
              }
              onDelete={handleDelete}
            />
          )
        )}
      </div>
      <Modal
  open={open}
  title="Upload Attachment"
  onClose={() =>
    setOpen(false)
  }
>
  <AttachmentForm
    submitText="Upload File"
    onClose={() =>
      setOpen(false)
    }
    onSubmit={async (
      file
    ) => {
      await uploadLeadAttachment(
        leadId,
        file
      );

      setOpen(false);

      load();
    }}
  />
</Modal>
    </div>
  );
}