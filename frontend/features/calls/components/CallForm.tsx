"use client";

import { useState } from "react";

import { CreateCallPayload } from "../types";

type Props = {
  initialValues?: CreateCallPayload;

  submitText: string;

  onSubmit: (
    values: CreateCallPayload
  ) => Promise<void>;

  onClose: () => void;
};

export default function CallForm({
  initialValues,
  submitText,
  onSubmit,
  onClose,
}: Props) {
  const [loading, setLoading] =
    useState(false);

  const [form, setForm] =
    useState<CreateCallPayload>(
      initialValues ?? {
        direction: "OUTGOING",
        status: "COMPLETED",
        duration_seconds: 0,
        summary: "",
        follow_up_at: "",
        }
    );

  function update(
    key: keyof CreateCallPayload,
    value: string | number
  ) {
    setForm((prev) => ({
      ...prev,
      [key]: value,
    }));
  }

  async function submit(
    e: React.FormEvent
  ) {
    e.preventDefault();

    try {
      setLoading(true);

      await onSubmit({
            ...form,
            duration_seconds: Number(
                form.duration_seconds
            ),
            follow_up_at: form.follow_up_at
                ? new Date(
                    form.follow_up_at
                ).toISOString()
                : undefined,
        });
    } finally {
      setLoading(false);
    }
  }

  return (
    <form
      onSubmit={submit}
      className="space-y-5"
    >
      {/* Direction */}

      <div>
        <label className="mb-2 block text-sm font-medium text-slate-700">
            Direction
        </label>

        <select
            value={form.direction}
            onChange={(e) =>
            update(
                "direction",
                e.target.value
            )
            }
            className="w-full rounded-2xl border border-slate-200 px-4 py-3"
        >
            <option value="OUTGOING">
            Outgoing
            </option>

            <option value="INCOMING">
            Incoming
            </option>
        </select>
        </div>

      {/* Outcome */}

            {/* Status */}

<div>
  <label className="mb-2 block text-sm font-medium text-slate-700">
    Status
  </label>

  <select
    value={form.status}
    onChange={(e) =>
      update(
        "status",
        e.target.value
      )
    }
    className="
    w-full
    rounded-2xl
    border
    border-slate-200
    px-4
    py-3
    outline-none
    focus:border-emerald-500
    "
  >
    <option value="COMPLETED">
      Completed
    </option>

    <option value="MISSED">
      Missed
    </option>

    <option value="NO_ANSWER">
      No Answer
    </option>

    <option value="VOICEMAIL">
      Voicemail
    </option>
  </select>
</div>
      {/* Duration */}

      <div>
        <label className="mb-2 block text-sm font-medium text-slate-700">
          Duration (minutes)
        </label>

        <input
          type="number"
          min={0}
          value={form.duration_seconds}
          onChange={(e) =>
            update(
              "duration_seconds",
              Number(
                e.target.value
              )
            )
          }
          className="
          w-full
          rounded-2xl
          border
          border-slate-200
          px-4
          py-3
          outline-none
          focus:border-emerald-500
          "
        />
      </div>

      {/* Date */}

      <div>
        <label className="mb-2 block text-sm font-medium text-slate-700">
          Follow Up At
        </label>

        <input
          type="datetime-local"
          value={form.follow_up_at}
          onChange={(e) =>
            update(
              "follow_up_at",
              e.target.value
            )
          }
          required
          className="
          w-full
          rounded-2xl
          border
          border-slate-200
          px-4
          py-3
          outline-none
          focus:border-emerald-500
          "
        />
      </div>

      {/* Notes */}

      <div>
        <label className="mb-2 block text-sm font-medium text-slate-700">
          Summary
        </label>

        <textarea
          rows={5}
          value={form.summary}
          onChange={(e) =>
            update(
              "summary",
              e.target.value
            )
          }
          placeholder="Conversation summary..."
          className="
          w-full
          rounded-2xl
          border
          border-slate-200
          px-4
          py-3
          outline-none
          focus:border-emerald-500
          "
        />
      </div>

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
          hover:bg-slate-100
          "
        >
          Cancel
        </button>

        <button
          disabled={loading}
          className="
          rounded-2xl
          bg-emerald-500
          px-6
          py-3
          font-medium
          text-white
          hover:bg-emerald-600
          disabled:opacity-50
          "
        >
          {loading
            ? "Saving..."
            : submitText}
        </button>
      </div>
    </form>
  );
}