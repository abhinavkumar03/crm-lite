"use client";

import { useEffect, useState } from "react";

import { getDashboard } from "@/features/dashboard/api";
import { DashboardResponse } from "@/features/dashboard/types";
import DashboardHeader from "@/features/dashboard/components/DashboardHeader";
import MetricsGrid from "@/features/dashboard/components/MetricsGrid";
import ModuleCountsCard from "@/features/dashboard/components/ModuleCountsCard";
import RecentRecordsCard from "@/features/dashboard/components/RecentRecordsCard";
import QuickActionsCard from "@/features/dashboard/components/QuickActionsCard";

export default function DashboardPage() {
  const [dashboard, setDashboard] = useState<DashboardResponse | null>(
    null
  );
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboard();
  }, []);

  async function loadDashboard(refresh = false) {
    try {
      setLoading(true);
      const data = await getDashboard(refresh);
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
          {Array.from({ length: 4 }).map((_, index) => (
            <div
              key={index}
              className="h-36 animate-pulse rounded-2xl bg-slate-200"
            />
          ))}
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
    <div data-tour="dashboard" className="space-y-8">
      <DashboardHeader onRefresh={() => loadDashboard(true)} />

      <MetricsGrid dashboard={dashboard} />

      <div className="grid gap-6 xl:grid-cols-3">
        <div className="xl:col-span-2">
          <ModuleCountsCard modules={dashboard.module_counts ?? []} />
        </div>
        <QuickActionsCard />
      </div>

      <RecentRecordsCard records={dashboard.recent_records ?? []} />
    </div>
  );
}
