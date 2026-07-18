"use client";

import Link from "next/link";
import { ArrowRight } from "lucide-react";

import { RecentRecord } from "../types";

type Props = {
  records: RecentRecord[];
};

export default function RecentRecordsCard({ records }: Props) {
  return (
    <section className="rounded-3xl border border-slate-200 bg-white shadow-sm">
      <div className="border-b border-slate-100 p-6">
        <p className="text-sm font-medium text-slate-500">Activity</p>
        <h2 className="mt-1 text-2xl font-bold text-slate-900">
          Recent records
        </h2>
      </div>

      <div className="divide-y divide-slate-100">
        {records.length === 0 ? (
          <p className="p-6 text-sm text-slate-500">
            No records yet. Create one from Forms.
          </p>
        ) : (
          records.map((r) => (
            <Link
              key={r.id}
              href={`/tables?module=${r.module_id}`}
              className="flex items-center justify-between px-6 py-4 transition hover:bg-slate-50"
            >
              <div>
                <p className="font-semibold text-slate-900">{r.title}</p>
                <p className="text-sm text-slate-500">
                  {r.module_label} ·{" "}
                  {new Date(r.created_at).toLocaleDateString()}
                </p>
              </div>
              <ArrowRight size={16} className="text-slate-400" />
            </Link>
          ))
        )}
      </div>
    </section>
  );
}
