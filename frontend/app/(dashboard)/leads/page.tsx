"use client";

import { useEffect, useState } from "react";

import { createLead, deleteLead, getLeads, updateLead } from "@/features/leads/api";

import { Lead } from "@/features/leads/types";

import LeadTable from "@/features/leads/components/LeadTable";
import LeadForm from "@/features/leads/components/LeadForm";
import Modal from "@/components/common/Modal";

export default function LeadsPage() {

    const [leads, setLeads] =
        useState<Lead[]>([]);

    const [page, setPage] =
        useState(1);

    const [search, setSearch] =
        useState("");

    const [editingLead, setEditingLead] = useState<Lead | null>(null);

    const [open, setOpen] = useState(false);

    useEffect(() => {

        loadLeads();

    }, [page]);

    async function loadLeads() {

        const res =
            await getLeads(
                page,
                search,
            );

        setLeads(res.data);

    }

    async function handleDelete(lead: Lead) {

        const ok = window.confirm(
            "Delete this lead?"
        );

        if (!ok) {
            return;
        }

        await deleteLead(lead.id);

        await loadLeads();
    }

    return (

        <div className="space-y-6">

            <div className="flex items-center justify-between">

                <h1 className="text-3xl font-bold">

                    Leads

                </h1>

                <button
                    onClick={() => setOpen(true)}
                    className="rounded bg-blue-600 px-4 py-2 text-white"
                >
                    Create Lead
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
                        loadLeads();
                    }}
                    className="rounded bg-blue-600 px-4 py-2 text-white"
                >
                    Search
                </button>
            </div>

            <LeadTable
                leads={leads}
                page={page}
                setPage={setPage}
                onEdit={setEditingLead}
                onDelete={handleDelete}
            />

            <Modal
                open={!!editingLead}
                title="Edit Lead"
                onClose={() => setEditingLead(null)}
            >
                {editingLead && (
                    <LeadForm
                        initialValues={{
                            name: editingLead.name,
                            email: editingLead.email,
                            phone: editingLead.phone,
                            company: editingLead.company,
                            status: editingLead.status,
                            notes: editingLead.notes,
                        }}
                        submitText="Update Lead"
                        onSubmit={async (values) => {
                            await updateLead(editingLead.id, values);

                            setEditingLead(null);

                            loadLeads();
                        }}
                    />
                )}
            </Modal>

        </div>

    );

}