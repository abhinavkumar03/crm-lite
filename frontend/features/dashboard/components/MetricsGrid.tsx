import { Boxes, Database, Mail, MessageCircle, Workflow, Zap } from "lucide-react";

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
          : "Sent / delivered",
    },
    {
      title: "Active workflows",
      value: dashboard.active_workflows ?? 0,
      icon: Workflow,
      color: "bg-indigo-500",
      trend: `${dashboard.disabled_workflows ?? 0} disabled`,
    },
    {
      title: "Workflow runs today",
      value: dashboard.workflows_executed_today ?? 0,
      icon: Zap,
      color: "bg-rose-500",
      trend:
        (dashboard.workflows_failed_today ?? 0) > 0
          ? `${dashboard.workflows_failed_today} failed`
          : "Executions",
    },
  ];

  return (
    <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
      {metrics.map((m) => (
        <MetricCard key={m.title} {...m} />
      ))}
    </div>
  );
}
