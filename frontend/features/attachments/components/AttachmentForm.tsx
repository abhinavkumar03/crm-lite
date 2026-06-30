"use client";

import { useState } from "react";

type Props = {
  submitText: string;

  onSubmit: (
    file: File
  ) => Promise<void>;

  onClose: () => void;
};

export default function AttachmentForm({
  submitText,
  onSubmit,
  onClose,
}: Props) {
  const [file, setFile] =
    useState<File | null>(null);

  const [loading, setLoading] =
    useState(false);

  async function submit(
    e: React.FormEvent
  ) {
    e.preventDefault();

    if (!file) {
      return;
    }

    try {
      setLoading(true);

      await onSubmit(file);
    } finally {
      setLoading(false);
    }
  }

  return (
    <form
      onSubmit={submit}
      className="space-y-6"
    >
      <div>
        <label className="mb-2 block text-sm font-medium text-slate-700">
          Select File
        </label>

        <input
          type="file"
          onChange={(e) =>
            setFile(
              e.target.files?.[0] ??
                null
            )
          }
          className="
          block
          w-full
          rounded-2xl
          border
          border-slate-200
          px-4
          py-3
          file:mr-4
          file:rounded-lg
          file:border-0
          file:bg-emerald-500
          file:px-4
          file:py-2
          file:text-white
          "
          required
        />
      </div>

      {file && (
        <div className="rounded-2xl bg-slate-50 p-4">
          <p className="font-medium">
            {file.name}
          </p>

          <p className="mt-1 text-sm text-slate-500">
            {(
              file.size /
              1024
            ).toFixed(1)}{" "}
            KB
          </p>
        </div>
      )}

      <div className="flex justify-end gap-3">
        <button
          type="button"
          onClick={onClose}
          className="
          rounded-2xl
          border
          border-slate-200
          px-5
          py-3
          transition
          hover:bg-slate-100
          "
        >
          Cancel
        </button>

        <button
          disabled={
            loading || !file
          }
          className="
          rounded-2xl
          bg-emerald-500
          px-6
          py-3
          font-medium
          text-white
          transition
          hover:bg-emerald-600
          disabled:cursor-not-allowed
          disabled:opacity-50
          "
        >
          {loading
            ? "Uploading..."
            : submitText}
        </button>
      </div>
    </form>
  );
}