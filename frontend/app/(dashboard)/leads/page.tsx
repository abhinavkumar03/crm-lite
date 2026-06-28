"use client";

import { useEffect, useState } from "react";

import { getLeads } from "@/features/leads/api";

import { Lead } from "@/features/leads/types";

import LeadTable from "@/features/leads/components/LeadTable";

export default function LeadsPage() {

    const [leads, setLeads] =
        useState<Lead[]>([]);

    const [page, setPage] =
        useState(1);

    const [search, setSearch] =
        useState("");

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

    return (

        <div className="space-y-6">

            <div className="flex items-center justify-between">

                <h1 className="text-3xl font-bold">

                    Leads

                </h1>

                <button
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
            />

        </div>

    );

}