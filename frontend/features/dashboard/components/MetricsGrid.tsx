import {
  Users,
  UserRound,
  CheckSquare,
  Trophy,
} from "lucide-react";

import { DashboardResponse } from "../types";

import MetricCard from "./MetricCard";

type Props = {
  dashboard: DashboardResponse;
};

export default function MetricsGrid({
  dashboard,
}: Props) {
  const metrics = [
    {
      title: "Total Leads",
      value: dashboard.total_leads,
      icon: Users,
      color: "bg-emerald-500",
      trend: "+18%",
    },
    {
      title: "Contacts",
      value: dashboard.total_contacts,
      icon: UserRound,
      color: "bg-blue-500",
      trend: "+8%",
    },
    {
      title: "Tasks",
      value: dashboard.total_tasks,
      icon: CheckSquare,
      color: "bg-amber-500",
      trend: "+15%",
    },
    {
      title: "Won Leads",
      value: dashboard.won_leads,
      icon: Trophy,
      color: "bg-violet-500",
      trend: "+22%",
    },
  ];

  return (
    <section className="grid gap-6 sm:grid-cols-2 xl:grid-cols-4">
      {metrics?.map((metric) => (
        <MetricCard
          key={metric.title}
          title={metric.title}
          value={metric.value}
          icon={metric.icon}
          color={metric.color}
          trend={metric.trend}
        />
      ))}
    </section>
  );
}