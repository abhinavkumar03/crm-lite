"use client";

import {
  Building2,
  Mail,
  Phone,
  Pencil,
  Trash2,
  User,
} from "lucide-react";

import { Contact } from "../types";

type Props = {
  contact: Contact;
  onEdit: (contact: Contact) => void;
  onDelete: (contact: Contact) => void;
};

export default function ContactCard({
  contact,
  onEdit,
  onDelete,
}: Props) {
  const fullName =
    `${contact.first_name} ${contact.last_name}`.trim();

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

      <div className="flex items-start gap-4">
        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100 text-emerald-600">
          <User size={22} />
        </div>

        <div className="min-w-0 flex-1">
          <h3 className="truncate text-lg font-semibold text-slate-900">
            {fullName}
          </h3>

          <div className="mt-2 flex items-center gap-2 text-sm text-slate-500">
            <Building2
              size={15}
              className="flex-shrink-0"
            />

            <span className="truncate">
              {contact.company || "No company"}
            </span>
          </div>
        </div>
      </div>

      {/* Contact */}

      <div className="mt-6 space-y-4">
        <div className="flex items-center gap-3">
          <Mail
            size={17}
            className="text-slate-400"
          />

          <span className="truncate text-sm text-slate-600">
            {contact.email}
          </span>
        </div>

        <div className="flex items-center gap-3">
          <Phone
            size={17}
            className="text-slate-400"
          />

          <span className="text-sm text-slate-600">
            {contact.phone || "-"}
          </span>
        </div>
      </div>

      {/* Actions */}

      <div className="mt-6 flex gap-3">
        <button
          onClick={() => onEdit(contact)}
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
          onClick={() => onDelete(contact)}
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