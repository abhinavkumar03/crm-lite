import { Boxes, Database } from "lucide-react";

import { DashboardResponse } from "../types";
import MetricCard from "./MetricCard";

type Props = {
  dashboard: DashboardResponse;
};

export default function MetricsGrid({ dashboard }: Props) {
  const metrics = [
    {
      title: "Modules",
      value: dashboard.total_modules,
      icon: Boxes,
      color: "bg-emerald-500",
      trend: "Dynamic",
    },
    {
      title: "Records",
      value: dashboard.total_records,
      icon: Database,
      color: "bg-blue-500",
      trend: "All modules",
    },
  ];

  return (
    <section className="grid gap-6 sm:grid-cols-2 xl:grid-cols-4">
      {metrics.map((metric) => (
        <MetricCard
          key={metric.title}
          title={metric.title}
          value={metric.value}
          icon={metric.icon}
          color={metric.color}
          trend={metric.trend}
        />
      ))}
      {(dashboard.module_counts ?? []).slice(0, 2).map((m) => (
        <MetricCard
          key={m.module_id}
          title={m.plural_label}
          value={m.record_count}
          icon={Boxes}
          color="bg-violet-500"
          trend={m.api_name}
        />
      ))}
    </section>
  );
}
