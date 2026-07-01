"use client";

import Link from "next/link";

import {
  ArrowRight,
} from "lucide-react";

import { DashboardActivity } from "../types";

import RecentActivityItem from "./RecentActivityItem";

type Props = {
  activities: DashboardActivity[];
};

export default function RecentActivityCard({
  activities,
}: Props) {
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
      {/* Header */}

      <div className="border-b border-slate-100 p-6">
        <h3 className="text-lg font-semibold text-slate-900">
          Recent Activity
        </h3>

        <p className="mt-1 text-sm text-slate-500">
          Latest updates across leads, contacts and tasks.
        </p>
      </div>

      {/* Body */}

      {activities.length === 0 ? (
        <div className="p-10 text-center">
          <p className="font-medium text-slate-700">
            No recent activity
          </p>

          <p className="mt-2 text-sm text-slate-500">
            Activity will appear here automatically.
          </p>
        </div>
      ) : (
        <div className="divide-y divide-slate-100">
          {activities.map((activity) => (
            <RecentActivityItem
              key={activity.id}
              activity={activity}
            />
          ))}
        </div>
      )}

      {/* Footer */}

      <div className="border-t border-slate-100 p-4">
        <Link
          href="/activity"
          className="
          flex
          items-center
          justify-center
          gap-2
          rounded-xl
          py-2
          text-sm
          font-medium
          text-emerald-600
          transition
          hover:bg-emerald-50
          "
        >
          View Activity

          <ArrowRight size={16} />
        </Link>
      </div>
    </div>
  );
}