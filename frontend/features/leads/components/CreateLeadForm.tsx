"use client";

import { useState } from "react";

import {
    createLead,
    CreateLeadPayload,
} from "../api";

interface Props {
    onSuccess: () => void;
}

export default function CreateLeadForm({
    onSuccess,
}: Props) {

    const [loading, setLoading] =
        useState(false);

    const [form, setForm] =
        useState<CreateLeadPayload>({
            name: "",
            email: "",
            phone: "",
            company: "",
            status: "NEW",
            notes: "",
        });

    function update(
        key: keyof CreateLeadPayload,
        value: string
    ) {

        setForm(prev => ({
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

            await createLead(form);

            onSuccess();

        } finally {

            setLoading(false);

        }

    }

    return (

        <form
            onSubmit={submit}
            className="space-y-4"
        >

            <input
                className="w-full rounded border p-3"
                placeholder="Name"
                value={form.name}
                onChange={(e)=>
                    update("name", e.target.value)
                }
                required
            />

            <input
                className="w-full rounded border p-3"
                placeholder="Email"
                value={form.email}
                onChange={(e)=>
                    update("email", e.target.value)
                }
                required
            />

            <input
                className="w-full rounded border p-3"
                placeholder="Phone"
                value={form.phone}
                onChange={(e)=>
                    update("phone", e.target.value)
                }
            />

            <input
                className="w-full rounded border p-3"
                placeholder="Company"
                value={form.company}
                onChange={(e)=>
                    update("company", e.target.value)
                }
            />

            <textarea
                className="w-full rounded border p-3"
                placeholder="Notes"
                value={form.notes}
                onChange={(e)=>
                    update("notes", e.target.value)
                }
            />

            <select
                className="w-full rounded border p-3"
                value={form.status}
                onChange={(e)=>
                    update("status", e.target.value)
                }
            >
                <option value="NEW">New</option>
                <option value="CONTACTED">Contacted</option>
                <option value="QUALIFIED">Qualified</option>
                <option value="WON">Won</option>
                <option value="LOST">Lost</option>
            </select>

            <button
                disabled={loading}
                className="w-full rounded bg-blue-600 py-3 text-white"
            >
                {loading
                    ? "Creating..."
                    : "Create Lead"}
            </button>

        </form>

    );

}