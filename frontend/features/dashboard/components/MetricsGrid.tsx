import { Boxes, Database, Mail, MessageCircle } from "lucide-react";

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
    {
      title: "Emails today",
      value: dashboard.emails_sent_today ?? 0,
      icon: Mail,
      color: "bg-amber-500",
      trend: "Sent / delivered",
    },
    {
      title: "WhatsApp today",
      value: dashboard.whatsapp_sent_today ?? 0,
      icon: MessageCircle,
      color: "bg-teal-500",
      trend:
        (dashboard.failed_notifications ?? 0) > 0
          ? `${dashboard.failed_notifications} failed`
          : (dashboard.scheduled_notifications ?? 0) > 0
            ? `${dashboard.scheduled_notifications} scheduled`
            : "Sent / delivered",
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
    </section>
  );
}
