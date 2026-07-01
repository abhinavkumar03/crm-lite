"use client";

import {
  ContactRound,
  History,
  MessageSquare,
  Paperclip,
  Phone,
  CheckCircle2,
  Users,
} from "lucide-react";

import { DashboardActivity } from "../types";

import { formatRelativeTime } from "@/features/activity/utils";

type Props = {
  activity: DashboardActivity;
};

export default function RecentActivityItem({
  activity,
}: Props) {
  function icon() {
    switch (activity.action) {
      case "NOTE_ADDED":
      case "NOTE_UPDATED":
      case "NOTE_DELETED":
        return (
          <MessageSquare
            size={18}
          />
        );

      case "ATTACHMENT_ADDED":
      case "ATTACHMENT_DELETED":
        return (
          <Paperclip
            size={18}
          />
        );

      case "CALL_LOGGED":
      case "CALL_UPDATED":
        return (
          <Phone
            size={18}
          />
        );

      case "TASK_CREATED":
      case "TASK_UPDATED":
      case "TASK_COMPLETED":
        return (
          <CheckCircle2
            size={18}
          />
        );

      case "LEAD_CREATED":
      case "LEAD_UPDATED":
        return (
          <Users
            size={18}
          />
        );

      case "CONTACT_CREATED":
      case "CONTACT_UPDATED":
        return (
          <ContactRound
            size={18}
          />
        );

      default:
        return (
          <History
            size={18}
          />
        );
    }
  }

  function badgeColor() {
    switch (
      activity.entity_type
    ) {
      case "LEAD":
        return "bg-emerald-100 text-emerald-700";

      case "CONTACT":
        return "bg-blue-100 text-blue-700";

      case "TASK":
        return "bg-violet-100 text-violet-700";

      default:
        return "bg-slate-100 text-slate-700";
    }
  }

  return (
    <div
      className="
      flex
      items-start
      gap-4
      p-5
      transition
      hover:bg-slate-50
      "
    >
      <div
        className="
        flex
        h-11
        w-11
        items-center
        justify-center
        rounded-full
        bg-slate-100
        text-slate-700
        "
      >
        {icon()}
      </div>

      <div className="min-w-0 flex-1">
        <h4 className="font-medium text-slate-900">
          {activity.description}
        </h4>

        <div className="mt-2 flex flex-wrap items-center gap-2">
          <span
            className={`
              rounded-full
              px-2.5
              py-1
              text-xs
              font-medium
              ${badgeColor()}
            `}
          >
            {activity.entity_type}
          </span>

          <span className="text-xs text-slate-500">
            {formatRelativeTime(
              activity.created_at
            )}
          </span>
        </div>
      </div>
    </div>
  );
}