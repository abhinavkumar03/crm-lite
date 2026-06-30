"use client";

import Link from "next/link";

import {
  Download,
  Trash2,
  Eye,
  Calendar,
  File,
  FileArchive,
  FileSpreadsheet,
  FileText,
  ImageIcon,
} from "lucide-react";

import { Attachment } from "../types";

type Props = {
  attachment: Attachment;

  onDelete: (
    attachment: Attachment
  ) => void;
};

export default function AttachmentCard({
  attachment,
  onDelete,
}: Props) {
  const extension =
    attachment.file_name
      .split(".")
      .pop()
      ?.toLowerCase() ?? "";

  const imageExtensions = [
    "jpg",
    "jpeg",
    "png",
    "gif",
    "webp",
    "svg",
  ];

  function renderPreview() {
    if (
      imageExtensions.includes(
        extension
      )
    ) {
      return (
        <img
          src={attachment.file_url}
          alt={attachment.file_name}
          className="
          h-44
          w-full
          rounded-2xl
          object-cover
          "
        />
      );
    }

    if (extension === "pdf") {
      return (
        <div className="flex h-44 items-center justify-center rounded-2xl bg-red-50">
          <FileText
            size={64}
            className="text-red-500"
          />
        </div>
      );
    }

    if (
      extension === "xls" ||
      extension === "xlsx"
    ) {
      return (
        <div className="flex h-44 items-center justify-center rounded-2xl bg-green-50">
          <FileSpreadsheet
            size={64}
            className="text-green-600"
          />
        </div>
      );
    }

    if (
      extension === "zip" ||
      extension === "rar"
    ) {
      return (
        <div className="flex h-44 items-center justify-center rounded-2xl bg-yellow-50">
          <FileArchive
            size={64}
            className="text-yellow-600"
          />
        </div>
      );
    }

    return (
      <div className="flex h-44 items-center justify-center rounded-2xl bg-slate-100">
        <File
          size={64}
          className="text-slate-500"
        />
      </div>
    );
  }

  function formatSize(
    bytes: number
  ) {
    if (bytes < 1024)
      return `${bytes} B`;

    if (
      bytes <
      1024 * 1024
    ) {
      return `${(
        bytes / 1024
      ).toFixed(1)} KB`;
    }

    return `${(
      bytes /
      1024 /
      1024
    ).toFixed(2)} MB`;
  }

  return (
    <div
      className="
      overflow-hidden
      rounded-3xl
      border
      border-slate-200
      bg-white
      shadow-sm
      transition
      hover:-translate-y-1
      hover:shadow-lg
      "
    >
      {renderPreview()}

      <div className="space-y-4 p-5">
        <div>
          <h3 className="truncate font-semibold text-slate-900">
            {attachment.file_name}
          </h3>

          <p className="mt-1 text-sm text-slate-500">
            {formatSize(
              attachment.file_size
            )}
          </p>
        </div>

        <div className="flex items-center gap-2 text-sm text-slate-500">
          <Calendar size={15} />

          {new Date(
            attachment.created_at
          ).toLocaleDateString()}
        </div>

        <div className="flex flex-col gap-2 sm:flex-row">
              <Link
            href={attachment.file_url}
            target="_blank"
            className="
            flex
            w-full
            items-center
            justify-center
            gap-2
            rounded-xl
            border
            border-slate-200
            py-3
            transition
            hover:bg-slate-50
            sm:flex-1
            "
          >
            <Eye size={18} />

            Preview
          </Link>

          <Link
            href={attachment.file_url}
            download
            className="
            flex
            w-full
            items-center
            justify-center
            gap-2
            rounded-xl
            border
            border-slate-200
            py-3
            transition
            hover:bg-slate-50
            sm:flex-1
            "
          >
            <Download
              size={18}
            />

            Download
          </Link>

          <button
            onClick={() =>
                onDelete(attachment)
            }
            className="
            flex
            w-full
            items-center
            justify-center
            gap-2
            rounded-xl
            border
            border-red-200
            py-3
            text-red-600
            transition
            hover:bg-red-50
            sm:w-auto
            sm:px-3
            "
          >
            <Trash2 size={18} />

            <span className="sm:hidden">
                Delete
            </span>
          </button>
        </div>
      </div>
    </div>
  );
}