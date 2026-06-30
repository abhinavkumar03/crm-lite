"use client";

import {
    Pencil,
    Trash2,
    Clock3,
} from "lucide-react";

import { Note } from "../types";

type Props = {
    note: Note;

    onEdit: (
        note: Note
    ) => void;

    onDelete: (
        note: Note
    ) => void;
};

export default function NoteCard({
    note,
    onEdit,
    onDelete,
}: Props) {
    const initials =
        note.user.name
            .split(" ")
            .map((n) => n[0])
            .join("")
            .substring(0, 2)
            .toUpperCase();

    return (
        <div
            className="
      rounded-3xl
      border
      border-slate-200
      bg-white
      p-6
      shadow-sm
      "
        >
            {/* Header */}

            <div
                className="
flex
flex-col
gap-5
md:flex-row
md:items-start
md:justify-between
"
            >
                <div className="flex gap-4">
                    <div
                        className="
            flex
            h-12
            w-12
            items-center
            justify-center
            rounded-full
            bg-gradient-to-br
            from-emerald-500
            to-teal-500
            font-semibold
            text-white
            "
                    >
                        {initials}
                    </div>

                    <div>
                        <h4 className="font-semibold text-slate-900">
                            {note.user.name}
                        </h4>

                        <div
                            className="
              mt-1
              flex
              items-center
              gap-2
              text-sm
              text-slate-500
              "
                        >
                            <Clock3 size={15} />

                            {new Date(
                                note.created_at
                            ).toLocaleString()}
                        </div>
                    </div>
                </div>

                <div className="
flex
justify-end
gap-2
">
                    <button
                        onClick={() =>
                            onEdit(note)
                        }
                        className="
            rounded-xl
            border
            border-slate-200
            p-2
            transition
            hover:bg-slate-100
            "
                    >
                        <Pencil size={18} />
                    </button>

                    <button
                        onClick={() =>
                            onDelete(note)
                        }
                        className="
            rounded-xl
            border
            border-red-200
            p-2
            text-red-600
            transition
            hover:bg-red-50
            "
                    >
                        <Trash2 size={18} />
                    </button>
                </div>
            </div>

            {/* Body */}

            <div className="mt-6">
                <p className="whitespace-pre-wrap leading-7 text-slate-700">
                    {note.note}
                </p>
            </div>
        </div>
    );
}