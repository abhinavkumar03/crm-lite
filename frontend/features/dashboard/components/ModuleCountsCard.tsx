"use client";

import Link from "next/link";
import { ArrowRight, Boxes } from "lucide-react";

import { ModuleCount } from "../types";

type Props = {
  modules: ModuleCount[];
};

export default function ModuleCountsCard({ modules }: Props) {
  return (
    <section className="rounded-3xl border border-slate-200 bg-white shadow-sm">
      <div className="border-b border-slate-100 p-6">
        <p className="text-sm font-medium text-slate-500">Catalog</p>
        <h2 className="mt-1 text-2xl font-bold text-slate-900">
          Modules
        </h2>
      </div>

      <div className="divide-y divide-slate-100">
        {modules.length === 0 ? (
          <p className="p-6 text-sm text-slate-500">
            No dynamic modules yet. Create one in Settings.
          </p>
        ) : (
          modules.map((m) => (
            <Link
              key={m.module_id}
              href={`/m/${m.api_name}`}
              className="flex items-center justify-between px-6 py-4 transition hover:bg-slate-50"
            >
              <div className="flex items-center gap-3">
                <div
                  className="flex h-10 w-10 items-center justify-center rounded-2xl"
                  style={{
                    backgroundColor: m.color || "#10b981",
                  }}
                >
                  <Boxes size={18} className="text-white" />
                </div>
                <div>
                  <p className="font-semibold text-slate-900">
                    {m.plural_label}
                  </p>
                  <p className="text-sm text-slate-500">{m.api_name}</p>
                </div>
              </div>
              <div className="flex items-center gap-3 text-slate-500">
                <span className="text-sm font-medium">
                  {m.record_count} records
                </span>
                <ArrowRight size={16} />
              </div>
            </Link>
          ))
        )}
      </div>
    </section>
  );
}
