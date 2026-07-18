"use client";

import Link from "next/link";
import { Upload, Download, ArrowRight } from "lucide-react";

const cards = [
  {
    href: "/imports",
    icon: Upload,
    title: "Import",
    description:
      "Bring records in from CSV or Excel. Columns are auto-mapped, rows are validated, and large files process in the background.",
    accent: "bg-sky-50 text-sky-600",
  },
  {
    href: "/exports",
    icon: Download,
    title: "Export",
    description:
      "Export any module to CSV or Excel — instantly or as a background job — and reuse saved export templates.",
    accent: "bg-emerald-50 text-emerald-600",
  },
];

export default function DataSettingsPage() {
  return (
    <div className="space-y-5">
      <div>
        <h2 className="text-lg font-semibold text-slate-900">Data</h2>
        <p className="text-sm text-slate-500">
          Move data in and out of your CRM. Both engines operate on dynamic
          modules and reuse the validation engine and record runtime.
        </p>
      </div>

      <div className="grid gap-5 sm:grid-cols-2">
        {cards.map((card) => {
          const Icon = card.icon;
          return (
            <Link
              key={card.href}
              href={card.href}
              className="group flex flex-col gap-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm transition hover:border-emerald-300 hover:shadow-md"
            >
              <span
                className={`flex h-12 w-12 items-center justify-center rounded-2xl ${card.accent}`}
              >
                <Icon className="h-6 w-6" />
              </span>
              <div>
                <p className="text-lg font-semibold text-slate-800">
                  {card.title}
                </p>
                <p className="mt-1 text-sm text-slate-500">{card.description}</p>
              </div>
              <span className="mt-auto inline-flex items-center gap-1 text-sm font-semibold text-emerald-600">
                Open {card.title.toLowerCase()}
                <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-1" />
              </span>
            </Link>
          );
        })}
      </div>
    </div>
  );
}
