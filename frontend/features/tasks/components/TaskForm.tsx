"use client";

import { useState } from "react";

import { CreateTaskPayload } from "../types";
import { Lead } from "@/features/leads/types";
import { Contact } from "@/features/contacts/types";

interface Props {

    initialValues?: CreateTaskPayload;

    submitText: string;

    leads: Lead[];

    contacts: Contact[];

    onSubmit: (
        values: CreateTaskPayload
    ) => Promise<void>;

}

export default function TaskForm({
    initialValues,
    submitText,
    leads,
    contacts,
    onSubmit,
}: Props) {

    const [loading, setLoading] =
        useState(false);

    const [form, setForm] =
        useState<CreateTaskPayload>(

            initialValues ?? {

                title: "",

                description: "",

                status: "PENDING",

                due_date: "",

                lead_id: "",

                contact_id: "",

            }

        );

    function update(
        key: keyof CreateTaskPayload,
        value: string,
    ) {

        setForm(prev => ({
            ...prev,
            [key]: value,
        }));

    }

    async function submit(
        e: React.FormEvent,
    ) {

        e.preventDefault();

        try {

            setLoading(true);

            const payload = {

                ...form,

                due_date: form.due_date
                    ? new Date(form.due_date).toISOString()
                    : undefined,

                lead_id:
                    form.lead_id || undefined,

                contact_id:
                    form.contact_id || undefined,

            };

            await onSubmit(payload);

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
                placeholder="Title"
                value={form.title}
                onChange={(e) =>
                    update("title", e.target.value)
                }
                required
            />

            <textarea
                className="w-full rounded border p-3"
                placeholder="Description"
                value={form.description}
                onChange={(e) =>
                    update("description", e.target.value)
                }
            />

            <select
                className="w-full rounded border p-3"
                value={form.status}
                onChange={(e) =>
                    update("status", e.target.value)
                }
            >

                <option value="PENDING">

                    Pending

                </option>

                <option value="IN_PROGRESS">

                    In Progress

                </option>

                <option value="COMPLETED">

                    Completed

                </option>

            </select>

            <input
                type="datetime-local"
                className="w-full rounded border p-3"
                value={form.due_date ?? ""}
                onChange={(e) =>
                    update("due_date", e.target.value)
                }
            />

            <select
                className="w-full rounded border p-3"
                value={form.lead_id ?? ""}
                onChange={(e) =>
                    update("lead_id", e.target.value)
                }
            >

                <option value="">

                    None

                </option>

                {leads.map((lead) => (

                    <option
                        key={lead.id}
                        value={lead.id}
                    >

                        {lead.name}

                    </option>

                ))}

            </select>

            <select
                className="w-full rounded border p-3"
                value={form.contact_id ?? ""}
                onChange={(e) =>
                    update("contact_id", e.target.value)
                }
            >

                <option value="">

                    None

                </option>

                {contacts.map((contact) => (

                    <option
                        key={contact.id}
                        value={contact.id}
                    >

                        {contact.first_name} {contact.last_name}

                    </option>

                ))}

            </select>

            <button
                disabled={loading}
                className="w-full rounded bg-blue-600 py-3 text-white"
            >

                {loading
                    ? "Saving..."
                    : submitText}

            </button>

        </form>

    );

}