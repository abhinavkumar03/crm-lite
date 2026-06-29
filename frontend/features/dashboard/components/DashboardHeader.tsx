"use client";

import { CalendarDays, Plus, RefreshCw } from "lucide-react";

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
      {/* Left */}

      <div>
        <div className="inline-flex items-center rounded-full bg-emerald-50 px-4 py-1 text-sm font-medium text-emerald-700">
          CRM Dashboard
        </div>

        <h1 className="mt-4 text-4xl font-bold tracking-tight text-slate-900">
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

      {/* Right */}

      <div className="flex flex-wrap gap-3">
        <button
          onClick={onRefresh}
          className="inline-flex items-center gap-2 rounded-xl border border-slate-300 bg-white px-5 py-3 font-medium transition-all duration-300 hover:bg-slate-100"
        >
          <RefreshCw size={18} />
          Refresh
        </button>

        <button className="inline-flex items-center gap-2 rounded-xl bg-emerald-500 px-5 py-3 font-medium text-white shadow transition-all duration-300 hover:-translate-y-0.5 hover:bg-emerald-600">
          <Plus size={18} />
          New Lead
        </button>
      </div>
    </section>
  );
}