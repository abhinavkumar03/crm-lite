"use client";

import { useEffect, useState } from "react";

import { getDashboard } from "@/features/dashboard/api";

import { DashboardResponse } from "@/features/dashboard/types";

import MetricsGrid from "@/features/dashboard/components/MetricsGrid";

import LeadStatusCard from "@/features/dashboard/components/LeadStatusCard";

import TaskStatusCard from "@/features/dashboard/components/TaskStatusCard";

import RecentLeadsCard from "@/features/dashboard/components/RecentLeadsCard";

import UpcomingTasksCard from "@/features/dashboard/components/UpcomingTasksCard";

export default function DashboardPage() {

    const [dashboard, setDashboard] =
        useState<DashboardResponse | null>(null);

    useEffect(() => {

        loadDashboard();

    }, []);

    async function loadDashboard() {

        const data =
            await getDashboard();

        setDashboard(data);

    }

    if (!dashboard) {

        return <p>Loading...</p>;

    }

    return (

        <div className="space-y-6">

            <MetricsGrid dashboard={dashboard} />

            <div className="grid grid-cols-2 gap-6">

                <LeadStatusCard dashboard={dashboard} />

                <TaskStatusCard dashboard={dashboard} />

            </div>

            <div className="grid grid-cols-2 gap-6">

                <RecentLeadsCard
                    leads={dashboard.recent_leads}
                />

                <UpcomingTasksCard
                    tasks={dashboard.upcoming_tasks}
                />

            </div>

        </div>

    );

}