"use client";

import Link from "next/link";
import { CalendarDays, CircleHelp, RefreshCw } from "lucide-react";

type Props = {
  onRefresh?: () => void;
};

export default function DashboardHeader({ onRefresh }: Props) {
  const today = new Date().toLocaleDateString("en-US", {
    weekday: "long",
    month: "long",
    day: "numeric",
    year: "numeric",
  });

  return (
    <section className="flex flex-col gap-6 rounded-3xl border border-slate-200 bg-white p-8 shadow-sm lg:flex-row lg:items-center lg:justify-between">
      <div>
        <h1 className="text-4xl font-bold tracking-tight text-slate-900">
          Welcome back 👋
        </h1>

        <p className="mt-2 max-w-2xl text-slate-500">
          Manage leads, contacts, tasks and monitor your CRM performance from
          one place.
        </p>

        <div className="mt-5 flex items-center gap-2 text-sm text-slate-500">
          <CalendarDays size={18} />
          <span>{today}</span>
        </div>
      </div>

      <div className="flex flex-wrap gap-3">
        <Link
          href="/help"
          className="inline-flex items-center gap-2 rounded-xl border border-emerald-200 bg-emerald-50 px-5 py-3 font-medium text-emerald-700 transition-all duration-300 hover:border-emerald-300 hover:bg-emerald-100"
        >
          <CircleHelp size={18} />
          How it works
        </Link>

        <button
          type="button"
          onClick={onRefresh}
          className="inline-flex items-center gap-2 rounded-xl border border-slate-300 bg-white px-5 py-3 font-medium transition-all duration-300 hover:bg-slate-100"
        >
          <RefreshCw size={18} />
          Refresh
        </button>
      </div>
    </section>
  );
}