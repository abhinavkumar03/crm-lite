"use client";

import {
  Calendar,
  Clock3,
  Pencil,
  Phone,
  Trash2,
} from "lucide-react";

import { CallLog } from "../types";

type Props = {
  call: CallLog;

  onEdit?: (
    call: CallLog
  ) => void;

  onDelete?: (
    call: CallLog
  ) => void;
};

export default function CallLogCard({
  call,
  onEdit,
  onDelete,
}: Props) {
  function getStatusColor() {
    switch (call.status) {
      case "COMPLETED":
        return "bg-emerald-100 text-emerald-700";

      case "MISSED":
        return "bg-red-100 text-red-700";

      case "NO_ANSWER":
        return "bg-amber-100 text-amber-700";

      case "VOICEMAIL":
        return "bg-sky-100 text-sky-700";

      default:
        return "bg-slate-100 text-slate-700";
    }
  }

  return (
    <div
      className="
      rounded-3xl
      border
      border-slate-200
      bg-white
      p-6
      shadow-sm
      transition
      hover:shadow-md
      "
    >
      <div
        className="
        flex
        flex-col
        gap-5
        md:flex-row
        md:items-start
        md:justify-between
        "
      >
        <div className="flex gap-4">
          <div
            className="
            flex
            h-12
            w-12
            items-center
            justify-center
            rounded-full
            bg-emerald-50
            "
          >
            <Phone
              size={22}
              className="text-emerald-600"
            />
          </div>

          <div>
            <h3 className="text-lg font-semibold text-slate-900">
              {call.direction} Call
            </h3>

            <div className="mt-3 flex flex-wrap items-center gap-3">
              <span
                className={`rounded-full px-3 py-1 text-xs font-semibold ${getStatusColor()}`}
              >
                {call.status.replaceAll(
                  "_",
                  " "
                )}
              </span>

              <div className="flex items-center gap-1 text-sm text-slate-500">
                <Clock3 size={15} />

                {call.duration_seconds} sec
              </div>

              {call.follow_up_at && (
                <div className="flex items-center gap-1 text-sm text-slate-500">
                  <Calendar size={15} />

                  Follow up{" "}
                  {new Date(
                    call.follow_up_at
                  ).toLocaleString()}
                </div>
              )}
            </div>
          </div>
        </div>

        <div
          className="
          flex
          flex-wrap
          justify-end
          gap-2
          "
        >
          <button
            onClick={() =>
              onEdit?.(call)
            }
            className="
            flex
            items-center
            gap-2
            rounded-xl
            border
            border-slate-200
            px-3
            py-2
            transition
            hover:bg-slate-100
            "
          >
            <Pencil size={18} />

            <span className="sm:hidden">
              Edit
            </span>
          </button>

          <button
            onClick={() =>
              onDelete?.(call)
            }
            className="
            flex
            items-center
            gap-2
            rounded-xl
            border
            border-red-200
            px-3
            py-2
            text-red-600
            transition
            hover:bg-red-50
            "
          >
            <Trash2 size={18} />

            <span className="sm:hidden">
              Delete
            </span>
          </button>
        </div>
      </div>

      {call.summary && (
        <div
          className="
          mt-6
          rounded-2xl
          bg-slate-50
          p-5
          "
        >
          <p className="whitespace-pre-wrap leading-7 text-slate-700">
            {call.summary}
          </p>
        </div>
      )}
    </div>
  );
}