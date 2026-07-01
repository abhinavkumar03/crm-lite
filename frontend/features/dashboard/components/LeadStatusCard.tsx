import {
  CircleDot,
  Phone,
  BadgeCheck,
  Trophy,
  XCircle,
} from "lucide-react";

import { DashboardResponse } from "../types";

type Props = {
  dashboard: DashboardResponse;
};

export default function LeadStatusCard({
  dashboard,
}: Props) {
  const items = [
    {
      label: "New",
      value: dashboard.new_leads,
      icon: CircleDot,
      color: "bg-sky-500",
    },
    {
      label: "Contacted",
      value: dashboard.contacted_leads,
      icon: Phone,
      color: "bg-amber-500",
    },
    {
      label: "Qualified",
      value: dashboard.qualified_leads,
      icon: BadgeCheck,
      color: "bg-indigo-500",
    },
    {
      label: "Won",
      value: dashboard.won_leads,
      icon: Trophy,
      color: "bg-emerald-500",
    },
    {
      label: "Lost",
      value: dashboard.lost_leads,
      icon: XCircle,
      color: "bg-red-500",
    },
  ];

  const total = items.reduce(
    (sum, item) => sum + item.value,
    0
  );

  return (
    <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
      {/* Header */}

      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-slate-500">
            Lead Analytics
          </p>

          <h2 className="mt-1 text-2xl font-bold text-slate-900">
            Sales Pipeline
          </h2>
        </div>

        {/* <div className="rounded-2xl bg-emerald-50 px-4 py-2 text-right">
          <p className="text-xs text-slate-500">
            Total Leads
          </p>

          <p className="text-2xl font-bold text-emerald-600">
            {total}
          </p>
        </div> */}
      </div>

      {/* Status */}

      <div className="mt-8 space-y-6">
        {items?.map((item) => {
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
                      size={18}
                      className="text-white"
                    />
                  </div>

                  <div>
                    <p className="font-medium text-slate-800">
                      {item.label}
                    </p>

                    <p className="text-xs text-slate-500">
                      {percent}% of leads
                    </p>
                  </div>
                </div>

                <div className="text-right">
                  <h3 className="text-lg font-bold text-slate-900">
                    {item.value}
                  </h3>
                </div>
              </div>

              <div className="h-2 overflow-hidden rounded-full bg-slate-100">
                <div
                  className={`h-full rounded-full ${item.color}`}
                  style={{
                    width: `${percent}%`,
                  }}
                />
              </div>
            </div>
          );
        })}
      </div>
    </section>
  );
}