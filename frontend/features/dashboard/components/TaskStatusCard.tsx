import {
  Clock3,
  PlayCircle,
  CheckCircle2,
  Activity,
} from "lucide-react";

import { DashboardResponse } from "../types";

type Props = {
  dashboard: DashboardResponse;
};

export default function TaskStatusCard({
  dashboard,
}: Props) {
  const items = [
    {
      label: "Pending",
      value: dashboard.pending_tasks,
      icon: Clock3,
      color: "bg-amber-500",
    },
    {
      label: "In Progress",
      value: dashboard.in_progress_tasks,
      icon: PlayCircle,
      color: "bg-blue-500",
    },
    {
      label: "Completed",
      value: dashboard.completed_tasks,
      icon: CheckCircle2,
      color: "bg-emerald-500",
    },
  ];

  const total = items.reduce(
    (sum, item) => sum + item.value,
    0
  );

  const completed =
    dashboard.completed_tasks;

  const completion =
    total === 0
      ? 0
      : Math.round((completed / total) * 100);

  return (
    <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
      {/* Header */}

      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-slate-500">
            Task Analytics
          </p>

          <h2 className="mt-1 text-2xl font-bold text-slate-900">
            Task Overview
          </h2>
        </div>

        <div className="rounded-2xl bg-blue-50 px-4 py-2 text-right">
          <p className="text-xs text-slate-500">
            Completion
          </p>

          <p className="text-2xl font-bold text-blue-600">
            {completion}%
          </p>
        </div>
      </div>

      {/* Overall Progress */}

      <div className="mt-8">
        <div className="mb-2 flex items-center justify-between">
          <span className="text-sm font-medium text-slate-600">
            Overall Progress
          </span>

          <span className="text-sm font-semibold text-slate-900">
            {completed}/{total}
          </span>
        </div>

        <div className="h-3 overflow-hidden rounded-full bg-slate-100">
          <div
            className="h-full rounded-full bg-gradient-to-r from-blue-500 to-emerald-500 transition-all duration-700"
            style={{
              width: `${completion}%`,
            }}
          />
        </div>
      </div>

      {/* Task Status */}

      <div className="mt-8 space-y-6">
        {items.map((item) => {
          const Icon = item.icon;

          const percent =
            total === 0
              ? 0
              : Math.round((item.value / total) * 100);

          return (
            <div key={item.label}>
              <div className="mb-2 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div
                    className={`flex h-10 w-10 items-center justify-center rounded-xl ${item.color}`}
                  >
                    <Icon
                      className="text-white"
                      size={18}
                    />
                  </div>

                  <div>
                    <p className="font-medium text-slate-800">
                      {item.label}
                    </p>

                    <p className="text-xs text-slate-500">
                      {percent}% of tasks
                    </p>
                  </div>
                </div>

                <span className="text-lg font-bold text-slate-900">
                  {item.value}
                </span>
              </div>

              <div className="h-2 overflow-hidden rounded-full bg-slate-100">
                <div
                  className={`${item.color} h-full rounded-full transition-all duration-500`}
                  style={{
                    width: `${percent}%`,
                  }}
                />
              </div>
            </div>
          );
        })}
      </div>

      {/* Footer */}

      <div className="mt-8 rounded-2xl bg-slate-50 p-4">
        <div className="flex items-center gap-3">
          <Activity
            size={20}
            className="text-emerald-500"
          />

          <div>
            <p className="text-sm font-medium text-slate-900">
              Productivity
            </p>

            <p className="text-xs text-slate-500">
              {completion >= 75
                ? "Excellent progress this week."
                : completion >= 50
                ? "You're making steady progress."
                : "Several tasks still need attention."}
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}