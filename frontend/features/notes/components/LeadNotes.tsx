"use client";

import { useEffect, useState } from "react";

import Modal from "@/components/common/Modal";
import NoteForm from "./NoteForm";

import {
    Plus,
    StickyNote,
} from "lucide-react";

import {
    getLeadNotes,
    createLeadNote,
    updateLeadNote,
    deleteLeadNote
} from "../api";

import { Note } from "../types";

import NoteCard from "./NoteCard";

type Props = {
    leadId: string;
};

export default function LeadNotes({
    leadId,
}: Props) {
    const [notes, setNotes] =
        useState<Note[]>([]);

    const [loading, setLoading] =
        useState(true);

    const [open, setOpen] =
        useState(false);

    const [editingNote, setEditingNote] =
        useState<Note | null>(
            null
        );

    async function loadNotes() {
        try {
            const data =
                await getLeadNotes(
                    leadId
                );

            setNotes(data);
        } finally {
            setLoading(false);
        }
    }

    async function handleDelete(
        note: Note
    ) {
        const ok =
            window.confirm(
                "Delete this note?"
            );

        if (!ok) {
            return;
        }

        await deleteLeadNote(
            leadId,
            note.id
        );

        loadNotes();
    }

    useEffect(() => {
        loadNotes();
    }, []);

    if (loading) {
        return (
            <div
                className="
        rounded-3xl
        border
        border-slate-200
        bg-white
        p-10
        text-center
        "
            >
                Loading notes...
            </div>
        );
    }

    return (
        <div className="space-y-6">
            {/* Header */}

            <div
                className="
        flex
        flex-col
        gap-4
        rounded-3xl
        border
        border-slate-200
        bg-white
        p-6
        shadow-sm
        md:flex-row
        md:items-center
        md:justify-between
        "
            >
                <div>
                    <h2 className="text-xl font-semibold text-slate-900">
                        Notes
                    </h2>

                    <p className="mt-1 text-sm text-slate-500">
                        Internal notes related to this lead.
                    </p>
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
          transition
          hover:bg-emerald-600
          "
                >
                    <Plus size={18} />

                    Add Note
                </button>
            </div>

            {/* Empty */}

            {notes.length === 0 && (
                <div
                    className="
          rounded-3xl
          border
          border-dashed
          border-slate-300
          bg-white
          p-14
          text-center
          "
                >
                    <div
                        className="
            mx-auto
            flex
            h-20
            w-20
            items-center
            justify-center
            rounded-full
            bg-emerald-50
            "
                    >
                        <StickyNote
                            size={36}
                            className="text-emerald-500"
                        />
                    </div>

                    <h3 className="mt-6 text-xl font-semibold text-slate-900">
                        No Notes Yet
                    </h3>

                    <p className="mx-auto mt-3 max-w-md text-slate-500">
                        Capture important discussions, meeting outcomes,
                        pricing decisions and follow-up reminders for this lead.
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

                        Create First Note
                    </button>
                </div>
            )}

            {/* Notes */}

            {notes.map((note) => (
                <NoteCard
                    key={note.id}
                    note={note}
                    onEdit={setEditingNote}
                    onDelete={handleDelete}
                />
            ))}
            <Modal
                open={open}
                title="Create Note"
                onClose={() =>
                    setOpen(false)
                }
            >
                <NoteForm
                    submitText="Create Note"
                    onSubmit={async (
                        values
                    ) => {
                        await createLeadNote(
                            leadId,
                            values
                        );

                        setOpen(false);

                        loadNotes();
                    }}
                />
            </Modal>
            <Modal
                open={!!editingNote}
                title="Edit Note"
                onClose={() =>
                    setEditingNote(null)
                }
            >
                {editingNote && (
                    <NoteForm
                        initialValue={
                            editingNote.note
                        }
                        submitText="Update Note"
                        onSubmit={async (
                            values
                        ) => {
                            await updateLeadNote(
                                editingNote.id,
                                values
                            );

                            setEditingNote(
                                null
                            );

                            loadNotes();
                        }}
                    />
                )}
            </Modal>
        </div>
    );
}