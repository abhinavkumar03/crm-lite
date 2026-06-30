"use client";

import Link from "next/link";

import {
  ArrowLeft,
  Building2,
  Mail,
  Phone,
  Pencil,
  Trash2,
  MoreHorizontal,
} from "lucide-react";

import StatusBadge from "@/components/common/table/StatusBadge";

import { LeadDetails } from "../../types";

type Props = {
  lead: LeadDetails;
};

export default function LeadDetailsHeader({
  lead,
}: Props) {
  const initials =
    lead.name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .substring(0, 2)
      .toUpperCase();

  return (
    <div
      className="
      rounded-3xl
      border
      border-slate-200
      bg-white
      shadow-sm
      "
    >
      {/* Top */}

      <div
        className="
        flex
        flex-col
        gap-6
        border-b
        border-slate-200
        p-6
        lg:flex-row
        lg:items-center
        lg:justify-between
        "
      >
        <div className="flex items-center gap-5">
          <div
            className="
            flex
            h-20
            w-20
            items-center
            justify-center
            rounded-3xl
            bg-gradient-to-br
            from-emerald-500
            to-teal-500
            text-2xl
            font-bold
            text-white
            shadow-lg
            "
          >
            {initials}
          </div>

          <div>
            <div className="flex flex-wrap items-center gap-3">
              <h1 className="text-3xl font-bold text-slate-900">
                {lead.name}
              </h1>

              <StatusBadge
                status={lead.status}
              />
            </div>

            <div className="mt-3 flex items-center gap-2 text-slate-600">
              <Building2 size={18} />

              <span>{lead.company}</span>
            </div>
          </div>
        </div>

        {/* Actions */}

        <div className="flex flex-wrap gap-3">
          <button
            className="
            flex
            items-center
            gap-2
            rounded-xl
            border
            border-slate-200
            px-5
            py-3
            font-medium
            transition
            hover:bg-slate-50
            "
          >
            <Pencil size={18} />

            Edit
          </button>

          <button
            className="
            flex
            items-center
            gap-2
            rounded-xl
            border
            border-red-200
            px-5
            py-3
            font-medium
            text-red-600
            transition
            hover:bg-red-50
            "
          >
            <Trash2 size={18} />

            Delete
          </button>

          <button
            className="
            rounded-xl
            border
            border-slate-200
            p-3
            transition
            hover:bg-slate-50
            "
          >
            <MoreHorizontal size={20} />
          </button>
        </div>
      </div>

      {/* Bottom */}

      <div
        className="
        grid
        gap-6
        p-6
        md:grid-cols-2
        "
      >
        <div className="flex items-center gap-3">
          <div className="rounded-xl bg-slate-100 p-3">
            <Mail
              size={18}
              className="text-slate-600"
            />
          </div>

          <div>
            <p className="text-xs uppercase tracking-wider text-slate-400">
              Email
            </p>

            <p className="font-medium text-slate-800">
              {lead.email}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <div className="rounded-xl bg-slate-100 p-3">
            <Phone
              size={18}
              className="text-slate-600"
            />
          </div>

          <div>
            <p className="text-xs uppercase tracking-wider text-slate-400">
              Phone
            </p>

            <p className="font-medium text-slate-800">
              {lead.phone}
            </p>
          </div>
        </div>
      </div>

      {/* Back */}

      <div className="border-t border-slate-200 px-6 py-4">
        <Link
          href="/leads"
          className="
          inline-flex
          items-center
          gap-2
          text-sm
          font-medium
          text-emerald-600
          transition
          hover:text-emerald-700
          "
        >
          <ArrowLeft size={16} />

          Back to Leads
        </Link>
      </div>
    </div>
  );
}