"use client";

import {
  CalendarDays,
  Pencil,
  Trash2,
  User,
  Users,
  ClipboardList,
} from "lucide-react";

import { Task } from "../types";

import StatusBadge from "@/components/common/table/StatusBadge";

type Props = {
  task: Task;
  onEdit: (task: Task) => void;
  onDelete: (task: Task) => void;
};

export default function TaskCard({
  task,
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

      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0">
          <div className="flex items-center gap-2">
            <ClipboardList
              size={20}
              className="text-emerald-500"
            />

            <h3 className="truncate text-lg font-semibold text-slate-900">
              {task.title}
            </h3>
          </div>

          <div className="mt-3">
            <StatusBadge status={task.status} />
          </div>
        </div>
      </div>

      {/* Description */}

      {task.description && (
        <p className="mt-5 line-clamp-3 text-sm leading-6 text-slate-600">
          {task.description}
        </p>
      )}

      {/* Details */}

      <div className="mt-6 space-y-4">
        <div className="flex items-center gap-3">
          <CalendarDays
            size={17}
            className="text-slate-400"
          />

          <span className="text-sm text-slate-600">
            {task.due_date
              ? new Date(
                  task.due_date
                ).toLocaleString()
              : "No due date"}
          </span>
        </div>

        <div className="flex items-center gap-3">
          <Users
            size={17}
            className="text-slate-400"
          />

          <span className="text-sm text-slate-600">
            Lead ID: {task.lead_id || "-"}
          </span>
        </div>

        <div className="flex items-center gap-3">
          <User
            size={17}
            className="text-slate-400"
          />

          <span className="text-sm text-slate-600">
            Contact ID: {task.contact_id || "-"}
          </span>
        </div>
      </div>

      {/* Footer */}

      <div className="mt-6 flex gap-3 border-t border-slate-100 pt-5">
        <button
          onClick={() => onEdit(task)}
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
          onClick={() => onDelete(task)}
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