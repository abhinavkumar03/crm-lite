"use client";

import { useEffect, useState } from "react";

import { getDashboard } from "@/features/dashboard/api";

import { DashboardResponse } from "@/features/dashboard/types";

import DashboardHeader from "@/features/dashboard/components/DashboardHeader";
import MetricsGrid from "@/features/dashboard/components/MetricsGrid";
import LeadStatusCard from "@/features/dashboard/components/LeadStatusCard";
import TaskStatusCard from "@/features/dashboard/components/TaskStatusCard";
import RecentLeadsCard from "@/features/dashboard/components/RecentLeadsCard";
import UpcomingTasksCard from "@/features/dashboard/components/UpcomingTasksCard";
import QuickActionsCard from "@/features/dashboard/components/QuickActionsCard";

export default function DashboardPage() {
    const [dashboard, setDashboard] =
        useState<DashboardResponse | null>(null);

    const [loading, setLoading] =
        useState(true);

    useEffect(() => {
        loadDashboard();
    }, []);

    async function loadDashboard() {
        try {
            setLoading(true);

            const data = await getDashboard();

            setDashboard(data);
        } finally {
            setLoading(false);
        }
    }

    if (loading) {
        return (
            <div className="space-y-6">
                <div className="h-44 animate-pulse rounded-3xl bg-slate-200" />

                <div className="grid gap-6 sm:grid-cols-2 xl:grid-cols-4">
                    {Array.from({ length: 4 })?.map((_, index) => (
                        <div
                            key={index}
                            className="h-36 animate-pulse rounded-2xl bg-slate-200"
                        />
                    ))}
                </div>

                <div className="grid gap-6 lg:grid-cols-2">
                    <div className="h-72 animate-pulse rounded-2xl bg-slate-200" />

                    <div className="h-72 animate-pulse rounded-2xl bg-slate-200" />
                </div>
            </div>
        );
    }

    if (!dashboard) {
        return (
            <div className="rounded-2xl border border-red-200 bg-red-50 p-8 text-center text-red-600">
                Failed to load dashboard.
            </div>
        );
    }

    return (
        <div className="space-y-8">
            <DashboardHeader
                onRefresh={loadDashboard}
            />

            <MetricsGrid
                dashboard={dashboard}
            />

            <div className="grid gap-6 xl:grid-cols-2">
                <LeadStatusCard
                    dashboard={dashboard}
                />

                <TaskStatusCard
                    dashboard={dashboard}
                />
            </div>

            <div className="grid gap-6 xl:grid-cols-3">
                <div className="xl:col-span-2">
                    <RecentLeadsCard
                        leads={dashboard.recent_leads}
                    />
                </div>

                <UpcomingTasksCard
                    tasks={dashboard.upcoming_tasks}
                />
            </div>

            <div className="grid gap-6 xl:grid-cols-3">
                <div className="xl:col-span-2">
                    <div className="rounded-3xl border border-dashed border-slate-300 bg-white p-12 text-center">
                        <h3 className="text-xl font-bold">
                            Recent Activity
                        </h3>

                        <p className="mt-2 text-slate-500">
                            Coming in the next phase.
                        </p>
                    </div>
                </div>

                <QuickActionsCard />
            </div>
        </div>
    );
}