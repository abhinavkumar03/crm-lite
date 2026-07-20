"use client";

import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { Building2, Check, ChevronDown, Plus } from "lucide-react";

import { useOrgOptional } from "@/features/organization/OrgContext";

type Props = {
  /** Compact = sidebar brand strip; menu = denser list for user menu embedding. */
  variant?: "sidebar" | "menu";
  className?: string;
};

export default function WorkspaceSwitcher({
  variant = "sidebar",
  className = "",
}: Props) {
  const org = useOrgOptional();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function onDoc(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    function onEsc(e: KeyboardEvent) {
      if (e.key === "Escape") setOpen(false);
    }
    document.addEventListener("mousedown", onDoc);
    document.addEventListener("keydown", onEsc);
    return () => {
      document.removeEventListener("mousedown", onDoc);
      document.removeEventListener("keydown", onEsc);
    };
  }, []);

  if (!org || org.loading) {
    return (
      <div className={`animate-pulse rounded-xl bg-slate-100 ${className}`}>
        <div className="h-10" />
      </div>
    );
  }

  const active = org.activeOrg;

  return (
    <div ref={ref} className={`relative ${className}`}>
      <button
        type="button"
        data-tour="workspace-switcher"
        onClick={() => setOpen((v) => !v)}
        disabled={org.switching}
        className={`
          flex w-full items-center gap-3 rounded-2xl border border-slate-200
          bg-white px-3 py-2.5 text-left transition
          hover:border-slate-300 hover:bg-slate-50
          disabled:opacity-60
          ${variant === "sidebar" ? "" : "shadow-sm"}
        `}
      >
        <OrgAvatar name={active?.name ?? "Workspace"} logoUrl={active?.logo_url} />
        <div className="min-w-0 flex-1">
          <p className="text-[10px] font-semibold uppercase tracking-wide text-slate-400">
            Current workspace
          </p>
          <p className="truncate text-sm font-semibold text-slate-800">
            {active?.name ?? "Select workspace"}
          </p>
        </div>
        <ChevronDown
          size={16}
          className={`shrink-0 text-slate-400 transition ${open ? "rotate-180" : ""}`}
        />
      </button>

      {open && (
        <div
          className={`
            absolute z-50 mt-2 w-full min-w-[240px] overflow-hidden rounded-2xl
            border border-slate-200 bg-white shadow-xl
            ${variant === "sidebar" ? "left-0" : "right-0"}
          `}
        >
          <div className="max-h-64 overflow-y-auto px-2 py-2">
            <p className="px-2 py-1 text-xs font-semibold uppercase tracking-wide text-slate-400">
              Workspaces
            </p>
            {org.orgs.map((item) => {
              const isActive = item.id === active?.id;
              return (
                <button
                  key={item.id}
                  type="button"
                  disabled={org.switching || isActive}
                  onClick={() => {
                    setOpen(false);
                    void org.switchWorkspace(item.id);
                  }}
                  className="flex w-full items-center gap-3 rounded-xl px-2 py-2.5 text-left text-slate-700 transition hover:bg-slate-100 disabled:opacity-60"
                >
                  <OrgAvatar name={item.name} logoUrl={item.logo_url} size="sm" />
                  <span className="min-w-0 flex-1 truncate text-sm">{item.name}</span>
                  {isActive ? <Check size={16} className="text-emerald-500" /> : null}
                </button>
              );
            })}
          </div>
          <div className="border-t border-slate-100 px-2 py-2">
            <Link
              href="/onboarding/organization"
              onClick={() => setOpen(false)}
              className="flex w-full items-center gap-3 rounded-xl px-2 py-2.5 text-sm font-semibold text-emerald-700 transition hover:bg-emerald-50"
            >
              <Plus size={16} />
              Create workspace
            </Link>
          </div>
        </div>
      )}
    </div>
  );
}

function OrgAvatar({
  name,
  logoUrl,
  size = "md",
}: {
  name: string;
  logoUrl?: string | null;
  size?: "sm" | "md";
}) {
  const dim = size === "sm" ? "h-7 w-7" : "h-9 w-9";
  if (logoUrl) {
    return (
      <span
        className={`flex ${dim} shrink-0 items-center justify-center overflow-hidden rounded-lg border border-slate-200 bg-white`}
      >
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img src={logoUrl} alt="" className="h-full w-full object-contain p-0.5" />
      </span>
    );
  }
  return (
    <span
      className={`flex ${dim} shrink-0 items-center justify-center rounded-lg bg-emerald-50 text-emerald-700`}
      title={name}
    >
      <Building2 size={size === "sm" ? 14 : 18} />
    </span>
  );
}
