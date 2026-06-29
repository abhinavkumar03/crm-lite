"use client";

import { Bell } from "lucide-react";

export default function NotificationBell() {
  return (
    <button
      className="
      relative
      rounded-2xl
      border
      border-slate-200
      bg-white
      p-3
      transition
      hover:bg-slate-100
      "
    >
      <Bell
        size={20}
        className="text-slate-700"
      />

      <span
        className="
        absolute
        right-2
        top-2
        h-2.5
        w-2.5
        rounded-full
        bg-emerald-500
        "
      />
    </button>
  );
}