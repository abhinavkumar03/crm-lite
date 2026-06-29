"use client";

import {
  Building2,
  Mail,
  Phone,
  Pencil,
  Trash2,
} from "lucide-react";

import { Lead } from "../types";
import StatusBadge from "@/components/common/table/StatusBadge";

type Props = {
  lead: Lead;
  onEdit: (lead: Lead) => void;
  onDelete: (lead: Lead) => void;
};

export default function LeadCard({
  lead,
  onEdit,
  onDelete,
}: Props) {
  return (
    <div
      className="
      rounded-3xl
      border
      border-slate-200
      bg-white
      p-5
      shadow-sm
      transition-all
      duration-300
      hover:-translate-y-0.5
      hover:shadow-md
      "
    >
      {/* Header */}

      <div className="flex items-start justify-between">
        <div>
          <h3 className="text-lg font-semibold text-slate-900">
            {lead.name}
          </h3>

          <div className="mt-2 flex items-center gap-2 text-sm text-slate-500">
            <Building2 size={15} />

            {lead.company || "No company"}
          </div>
        </div>

        <StatusBadge status={lead.status} />
      </div>

      {/* Information */}

      <div className="mt-6 space-y-4">
        <div className="flex items-center gap-3">
          <Mail
            size={17}
            className="text-slate-400"
          />

          <span className="truncate text-sm text-slate-600">
            {lead.email}
          </span>
        </div>

        <div className="flex items-center gap-3">
          <Phone
            size={17}
            className="text-slate-400"
          />

          <span className="text-sm text-slate-600">
            {lead.phone || "-"}
          </span>
        </div>
      </div>

      {/* Actions */}

      <div className="mt-6 flex gap-3 border-t border-slate-100 pt-5">
        <button
          onClick={() => onEdit(lead)}
          className="
          flex-1
          rounded-xl
          border
          border-slate-200
          py-2.5
          text-sm
          font-medium
          text-slate-700
          transition
          hover:bg-slate-50
          "
        >
          <span className="flex items-center justify-center gap-2">
            <Pencil size={16} />
            Edit
          </span>
        </button>

        <button
          onClick={() => onDelete(lead)}
          className="
          flex-1
          rounded-xl
          bg-red-500
          py-2.5
          text-sm
          font-medium
          text-white
          transition
          hover:bg-red-600
          "
        >
          <span className="flex items-center justify-center gap-2">
            <Trash2 size={16} />
            Delete
          </span>
        </button>
      </div>
    </div>
  );
}