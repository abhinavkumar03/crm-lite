"use client";

import { useState } from "react";
import { CreateContactPayload } from "../types";


interface Props {

    initialValues?: CreateContactPayload;

    submitText: string;

    onSubmit: (
        values: CreateContactPayload
    ) => Promise<void>;

}

export default function ContactForm({
    initialValues,
    onSubmit,
    submitText,
}: Props) {

    const [loading, setLoading] =
        useState(false);

    const [form, setForm] =
    useState<CreateContactPayload>(

        initialValues ?? {

            first_name: "",

            last_name: "",

            email: "",

            phone: "",

            company: "",

        }

    );

    function update(
        key: keyof CreateContactPayload,
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

            await onSubmit(form);

        } finally {

            setLoading(false);

        }

    }

    return (

        <form
            onSubmit={submit}
            className="space-y-4"
        >

            <div className="grid grid-cols-2 gap-4">

                <input
                    className="w-full rounded border p-3"
                    placeholder="First Name"
                    value={form.first_name}
                    onChange={(e) =>
                        update("first_name", e.target.value)
                    }
                    required
                />

                <input
                    className="w-full rounded border p-3"
                    placeholder="Last Name"
                    value={form.last_name}
                    onChange={(e) =>
                        update("last_name", e.target.value)
                    }
                    required
                />

            </div>

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