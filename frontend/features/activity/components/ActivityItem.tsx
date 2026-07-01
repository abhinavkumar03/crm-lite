"use client";

import {
  Calendar,
  FileText,
  History,
  MessageSquare,
  Paperclip,
  Phone,
} from "lucide-react";

import { Activity } from "../types";

import { formatRelativeTime } from "../utils";

import MetadataPreview from "./MetadataPreview";

type Props = {
    activity: Activity;

    isLast: boolean;
};

export default function ActivityItem({
    activity,
    isLast,
}: Props) {
  function getIcon() {
    switch (activity.action) {
      case "NOTE_ADDED":
      case "NOTE_UPDATED":
      case "NOTE_DELETED":
        return (
          <MessageSquare
            size={18}
            className="text-blue-600"
          />
        );

      case "ATTACHMENT_ADDED":
      case "ATTACHMENT_DELETED":
        return (
          <Paperclip
            size={18}
            className="text-emerald-600"
          />
        );

      case "CALL_LOGGED":
      case "CALL_UPDATED":
      case "CALL_DELETED":
        return (
          <Phone
            size={18}
            className="text-violet-600"
          />
        );

      case "LEAD_CREATED":
      case "LEAD_UPDATED":
      case "LEAD_STATUS_CHANGED":
        return (
          <History
            size={18}
            className="text-orange-600"
          />
        );

      default:
        return (
          <FileText
            size={18}
            className="text-slate-600"
          />
        );
    }
  }

  function getColor() {
    switch (activity.action) {
      case "NOTE_ADDED":
      case "NOTE_UPDATED":
      case "NOTE_DELETED":
        return "bg-blue-50";

      case "ATTACHMENT_ADDED":
      case "ATTACHMENT_DELETED":
        return "bg-emerald-50";

      case "CALL_LOGGED":
      case "CALL_UPDATED":
      case "CALL_DELETED":
        return "bg-violet-50";

      case "LEAD_CREATED":
      case "LEAD_UPDATED":
      case "LEAD_STATUS_CHANGED":
        return "bg-orange-50";

      default:
        return "bg-slate-100";
    }
  }

  return (
      <div className="relative flex items-start gap-4">
          {/* Timeline */}

      <div className="flex flex-col items-center">
        <div
          className={`
            flex
            h-12
            w-12
            items-center
            justify-center
            rounded-full
            ${getColor()}
          `}
        >
          {getIcon()}
        </div>

        {!isLast && (
            <div
                className="
                mt-2
                h-full
                w-px
                bg-slate-200
                "
            />
        )}      </div>

      {/* Card */}

      <div
        className="
        flex-1
        rounded-2xl
        sm:rounded-3xl
        border
        border-slate-200
        bg-white
        p-5
        sm:p-6
        shadow-sm
        "
      >
        <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <div>
            <h3 className="font-semibold text-slate-900">
              {activity.description}
            </h3>

            <p className="mt-1 text-xs uppercase tracking-wide text-slate-400">
              {activity.action.replaceAll(
                "_",
                " "
              )}
            </p>
          </div>

          <div className="flex items-center gap-2 text-sm text-slate-500">
            <Calendar size={15} />

            {formatRelativeTime(
                activity.created_at
                )}
          </div>
        </div>

        {activity.metadata && (
          <div
            className="
            mt-5
            rounded-2xl
            bg-slate-50
            p-4
            "
          >
            <MetadataPreview
              metadata={
                activity.metadata
              }
            />
          </div>
        )}
      </div>
    </div>
  );
}