"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Bell } from "lucide-react";

import { listNotifications } from "@/features/notifications/api";

export default function NotificationBell() {
  const [failed, setFailed] = useState(0);

  useEffect(() => {
    let active = true;

    (async () => {
      try {
        const result = await listNotifications({
          page_size: 1,
          status: "failed",
        });

        if (active) {
          setFailed(result.total);
        }
      } catch {
        // Bell is non-critical; ignore load errors.
      }
    })();

    return () => {
      active = false;
    };
  }, []);

  return (
    <Link
      href="/notifications"
      data-tour="notification-bell"
      aria-label="Open notification center"
    >
      <button
        className="
          relative
          rounded-2xl
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
          className={`
            absolute
            right-2
            top-2
            h-2.5
            w-2.5
            rounded-full
            ${failed > 0 ? "bg-red-500" : "bg-emerald-500"}
          `}
        />
      </button>
    </Link>
  );
}