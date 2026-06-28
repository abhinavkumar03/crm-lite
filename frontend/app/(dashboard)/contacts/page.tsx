"use client";

import { useEffect, useState } from "react";

import {
    createContact,
    deleteContact,
    getContacts,
    updateContact,
} from "@/features/contacts/api";

import ContactTable from "@/features/contacts/components/ContactTable";
import Modal from "@/components/common/Modal";
import { Contact } from "@/features/contacts/types";
import ContactForm from "@/features/contacts/components/ContactForm";

export default function LeadsPage() {

    const [contacts, setContacts] = useState<Contact[]>([]);

    const [page, setPage] = useState(1);

    const [search, setSearch] = useState("");

    const [open, setOpen] = useState(false);

    const [editingContact, setEditingContact] =
        useState<Contact | null>(null);

    useEffect(() => {

        loadContacts();

    }, [page]);

    async function loadContacts() {

        const res = await getContacts(
            page,
            search
        );

        setContacts(res.data);

    }

    async function handleDelete(contact: Contact) {

        const ok = window.confirm(
            "Delete this contact?"
        );

        if (!ok) {
            return;
        }

        await deleteContact(contact.id);

        await loadContacts();
    }

    return (

        <div className="space-y-6">

            <div className="flex items-center justify-between">

                <h1 className="text-3xl font-bold">

                    Contacts

                </h1>

                <button
                    onClick={() => setOpen(true)}
                    className="rounded bg-blue-600 px-4 py-2 text-white"
                >
                    Create Contact
                </button>

            </div>

            <div className="flex gap-2">
                <input
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    placeholder="Search..."
                    className="flex-1 rounded border p-3"
                />

                <button
                    onClick={() => {
                        setPage(1);
                        loadContacts();
                    }}
                    className="rounded bg-blue-600 px-4 py-2 text-white"
                >
                    Search
                </button>
            </div>

            <ContactTable
                contacts={contacts}
                page={page}
                setPage={setPage}
                onEdit={setEditingContact}
                onDelete={handleDelete}
            />

            <Modal
                open={open}
                title="Create Contact"
                onClose={() => setOpen(false)}
            >
                <ContactForm
                    submitText="Create Contact"
                    onSubmit={async (values) => {
                        await createContact(values);

                        setOpen(false);

                        loadContacts();
                    }}
                />
            </Modal>

            <Modal
                open={!!editingContact}
                title="Edit Contact"
                onClose={() => setEditingContact(null)}
            >
                {editingContact && (
                    <ContactForm
                        initialValues={{
                            first_name: editingContact.first_name,
                            last_name: editingContact.last_name,
                            email: editingContact.email,
                            phone: editingContact.phone,
                            company: editingContact.company,
                        }}
                        submitText="Update Contact"
                        onSubmit={async (values) => {
                            await updateContact(editingContact.id, values);

                            setEditingContact(null);

                            loadContacts();
                        }}
                    />
                )}
            </Modal>

        </div>

    );

}