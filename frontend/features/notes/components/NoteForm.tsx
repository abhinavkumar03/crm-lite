"use client";

import { useState } from "react";

import { CreateNotePayload } from "../types";

type Props = {
    initialValue?: string;

    submitText: string;

    onSubmit: (
        values: CreateNotePayload
    ) => Promise<void>;
};

export default function NoteForm({
    initialValue,
    submitText,
    onSubmit,
}: Props) {
    const [loading, setLoading] =
        useState(false);

    const [content, setContent] =
        useState(
            initialValue ?? ""
        );
    async function submit(
        e: React.FormEvent
    ) {
        e.preventDefault();

        if (!content.trim()) return;

        try {
            setLoading(true);

            await onSubmit({
                note: content,
            });

            setContent("");
        } finally {
            setLoading(false);
        }
    }

    return (
        <form
            onSubmit={submit}
            className="space-y-5"
        >
            <div>
                <label className="mb-2 block text-sm font-medium text-slate-700">
                    Note
                </label>

                <textarea
                    rows={6}
                    value={content}
                    onChange={(e) =>
                        setContent(
                            e.target.value
                        )
                    }
                    placeholder="Write an internal note..."
                    className="
          w-full
          rounded-2xl
          border
          border-slate-200
          bg-white
          px-4
          py-3
          outline-none
          transition
          focus:border-emerald-500
          "
                    required
                />
            </div>

            <button
                disabled={loading}
                className="
        w-full
        rounded-2xl
        bg-emerald-500
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
                    ? "Saving..."
                    : submitText}
            </button>
        </form>
    );
}