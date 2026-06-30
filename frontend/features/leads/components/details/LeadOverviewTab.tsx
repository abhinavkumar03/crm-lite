"use client";

import {
  CalendarDays,
  Clock3,
  Building2,
  Mail,
  Phone,
  StickyNote,
} from "lucide-react";

import { LeadDetails } from "../../types";

type Props = {
  lead: LeadDetails;
};

export default function LeadOverviewTab({
  lead,
}: Props) {
  return (
      <div className="grid gap-8 lg:grid-cols-3">

      <div className="space-y-6 lg:col-span-2">
        {/* Contact */}

        <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="mb-6 text-lg font-semibold text-slate-900">
            Contact Information
          </h2>

          <div className="space-y-5">
            <InfoRow
              icon={<Mail size={18} />}
              label="Email"
              value={lead.email}
            />

            <InfoRow
              icon={<Phone size={18} />}
              label="Phone"
              value={lead.phone}
            />

            <InfoRow
              icon={<Building2 size={18} />}
              label="Company"
              value={lead.company}
            />
          </div>
        </div>

        {/* Notes */}

        <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="mb-6 text-lg font-semibold">
            Lead Notes
          </h2>

          <div className="flex gap-3">
            <StickyNote
              size={18}
              className="mt-1 text-slate-400"
            />

            <p className="leading-7 text-slate-600">
              {lead.notes ||
                "No notes available."}
            </p>
          </div>
        </div>
      </div>

      {/* Right */}

      <div className="space-y-6">
        <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="mb-6 text-lg font-semibold">
            Timeline
          </h2>

          <div className="space-y-5">
            <InfoRow
              icon={
                <CalendarDays size={18} />
              }
              label="Created"
              value={new Date(
                lead.created_at
              ).toLocaleString()}
            />

            <InfoRow
              icon={<Clock3 size={18} />}
              label="Updated"
              value={new Date(
                lead.updated_at
              ).toLocaleString()}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

function InfoRow({
  icon,
  label,
  value,
}: {
  icon: React.ReactNode;
  label: string;
  value: string;
}) {
  return (
    <div className="flex items-start gap-4">
      <div className="rounded-xl bg-slate-100 p-3 text-slate-600">
        {icon}
      </div>

      <div>
        <p className="text-xs uppercase tracking-wider text-slate-400">
          {label}
        </p>

        <p className="mt-1 font-medium text-slate-800">
          {value}
        </p>
      </div>
    </div>
  );
}