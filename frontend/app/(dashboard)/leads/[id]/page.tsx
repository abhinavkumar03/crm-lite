"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";

import { getLeadDetails } from "@/features/leads/api";
import { LeadDetails } from "@/features/leads/types";


import LeadNotes from "@/features/notes/components/LeadNotes";
import LeadDetailsHeader from "@/features/leads/components/details/LeadDetailsHeader";
import LeadDetailsTabs from "@/features/leads/components/details/LeadDetailsTabs";
import LeadOverviewTab from "@/features/leads/components/details/LeadOverviewTab";
import LeadAttachments from "@/features/attachments/components/LeadAttachments";
import LeadCallLogs from "@/features/calls/components/LeadCallLogs";
import LeadActivityTimeline from "@/features/activity/components/LeadActivityTimeline";

export default function LeadDetailsPage() {
    const params = useParams();

    const [lead, setLead] =
        useState<LeadDetails | null>(null);

    const [loading, setLoading] =
        useState(true);

    const [activeTab, setActiveTab] =
        useState("overview");

    async function loadLead() {
        try {
            const data =
                await getLeadDetails(
                    params.id as string
                );

            setLead(data);
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        loadLead();
    }, []);

    if (loading) {
        return (
            <div className="py-20 text-center">
                Loading Lead...
            </div>
        );
    }

    if (!lead) {
        return (
            <div className="py-20 text-center">
                Lead not found.
            </div>
        );
    }

    return (
        <div className="mx-auto max-w-7xl space-y-6">
            <div className="sticky z-20 space-y-4">
                <LeadDetailsHeader lead={lead} />
                <LeadDetailsTabs active={activeTab} onChange={setActiveTab}/>
            </div>

            <div
            className="
            rounded-3xl
            "
            >
            {activeTab === "overview" && (
                <LeadOverviewTab
                lead={lead}
                />
            )}

            {activeTab === "notes" && (
                <LeadNotes
                leadId={lead.id}
                />
            )}

            {activeTab === "attachments" && (
            <LeadAttachments
                leadId={lead.id}
            />
            )}

            {activeTab === "calls" && (
            <LeadCallLogs
                leadId={lead.id}
            />
            )}


            {activeTab === "activity" && (
                <LeadActivityTimeline
                    leadId={lead.id}
                />
            )}
            </div>
        </div>
    );
}