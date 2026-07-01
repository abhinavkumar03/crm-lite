import {
  ArrowRight,
  Building2,
  UserRound,
} from "lucide-react";

import { Lead } from "../types";

type Props = {
  leads: Lead[];
};

export default function RecentLeadsCard({
  leads,
}: Props) {
  return (
    <section className="rounded-3xl border border-slate-200 bg-white shadow-sm">
      {/* Header */}

      <div className="flex items-center justify-between border-b border-slate-100 p-6">
        <div>
          <p className="text-sm font-medium text-slate-500">
            CRM
          </p>

          <h2 className="mt-1 text-2xl font-bold text-slate-900">
            Recent Leads
          </h2>
        </div>

        <button className="text-sm font-semibold text-emerald-600 transition hover:text-emerald-700">
          View All
        </button>
      </div>

      {/* Empty State */}

      {(leads ?? []).length === 0 && (
        <div className="flex flex-col items-center justify-center py-16">
          <UserRound
            size={48}
            className="text-slate-300"
          />

          <p className="mt-4 font-medium text-slate-700">
            No leads found
          </p>

          <p className="mt-1 text-sm text-slate-500">
            Create your first lead to get started.
          </p>
        </div>
      )}

      {/* Leads */}

      <div className="divide-y divide-slate-100">
        {leads?.map((lead, index) => (
          <div
            key={lead.id}
            className="group flex items-center justify-between p-5 transition-all duration-300 hover:bg-slate-50"
          >
            <div className="flex items-center gap-4">
              {/* Avatar */}

              <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-emerald-500 to-teal-500 text-lg font-bold text-white shadow-sm">
                {lead.name
                  ?.split(" ")
                  ?.map((word) => word[0])
                  .join("")
                  .slice(0, 2)
                  .toUpperCase()}
              </div>

              {/* Details */}

              <div>
                <h3 className="font-semibold text-slate-900">
                  {lead.name}
                </h3>

                <div className="mt-1 flex items-center gap-2 text-sm text-slate-500">
                  <Building2 size={14} />

                  <span>
                    {lead.company || "No Company"}
                  </span>
                </div>
              </div>
            </div>

            {/* Right */}

            <div className="flex items-center gap-4">
              <span
                className={`rounded-full px-3 py-1 text-xs font-semibold ${
                  index % 3 === 0
                    ? "bg-emerald-100 text-emerald-700"
                    : index % 3 === 1
                    ? "bg-blue-100 text-blue-700"
                    : "bg-amber-100 text-amber-700"
                }`}
              >
                {index % 3 === 0
                  ? "Qualified"
                  : index % 3 === 1
                  ? "Contacted"
                  : "New"}
              </span>

              <ArrowRight
                size={18}
                className="text-slate-400 transition group-hover:translate-x-1"
              />
            </div>
          </div>
        ))}
      </div>

      {/* Footer */}

      {(leads ?? []).length > 0 && (
        <div className="border-slate-100 p-5">
          <button className="flex w-full items-center justify-center gap-2 rounded-xl border border-slate-200 bg-white py-3 font-medium text-slate-700 transition hover:bg-slate-100">
            View Complete Lead List

            <ArrowRight size={18} />
          </button>
        </div>
      )}
    </section>
  );
}