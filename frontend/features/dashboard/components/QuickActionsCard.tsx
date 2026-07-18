"use client";

import Link from "next/link";
import {
  LayoutTemplate,
  Table2,
  Settings,
  ArrowRight,
} from "lucide-react";

const actions = [
  {
    title: "New record",
    description: "Open a metadata-driven form",
    href: "/forms",
    icon: LayoutTemplate,
    color: "bg-emerald-500",
  },
  {
    title: "Browse tables",
    description: "List and edit module records",
    href: "/tables",
    icon: Table2,
    color: "bg-blue-500",
  },
  {
    title: "Manage modules",
    description: "Create fields and validation",
    href: "/settings/modules",
    icon: Settings,
    color: "bg-violet-500",
  },
];

export default function QuickActionsCard() {
  return (
    <section className="rounded-3xl border border-slate-200 bg-white shadow-sm">
      <div className="border-b border-slate-100 p-6">
        <p className="text-sm font-medium text-slate-500">Productivity</p>
        <h2 className="mt-1 text-2xl font-bold text-slate-900">
          Quick Actions
        </h2>
      </div>

      <div className="grid gap-4 p-6">
        {actions.map((action) => {
          const Icon = action.icon;
          return (
            <Link
              key={action.href}
              href={action.href}
              className="
                group
                flex
                items-center
                justify-between
                rounded-2xl
                border
                border-slate-200
                p-4
                text-left
                transition-all
                duration-300
                hover:-translate-y-1
                hover:border-emerald-300
                hover:bg-emerald-50
                hover:shadow-lg
              "
            >
              <div className="flex items-center gap-4">
                <div
                  className={`flex h-12 w-12 items-center justify-center rounded-2xl ${action.color}`}
                >
                  <Icon size={20} className="text-white" />
                </div>
                <div>
                  <h3 className="font-semibold text-slate-900">
                    {action.title}
                  </h3>
                  <p className="text-sm text-slate-500">
                    {action.description}
                  </p>
                </div>
              </div>
              <ArrowRight
                size={18}
                className="text-slate-400 transition-transform duration-300 group-hover:translate-x-1"
              />
            </Link>
          );
        })}
      </div>
    </section>
  );
}
