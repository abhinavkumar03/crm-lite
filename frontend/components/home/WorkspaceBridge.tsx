"use client";

import Link from "next/link";
import {
  ArrowRight,
  BookOpen,
  LayoutDashboard,
  LogIn,
} from "lucide-react";

import { useAuth } from "@/context/AuthContext";

const stepsSignedOut = [
  {
    step: "1",
    title: "Sign in",
    detail: "Use the demo admin account or create your own.",
    href: "/login",
    icon: LogIn,
  },
  {
    step: "2",
    title: "Open dashboard",
    detail: "Land on live module counts, recent records, and tools.",
    href: "/login",
    icon: LayoutDashboard,
  },
  {
    step: "3",
    title: "Learn the system",
    detail: "Read How it works anytime from the sidebar or home.",
    href: "/help",
    icon: BookOpen,
  },
];

const stepsSignedIn = [
  {
    step: "1",
    title: "Your dashboard",
    detail: "Module counts, recent records, and quick actions are waiting.",
    href: "/dashboard",
    icon: LayoutDashboard,
  },
  {
    step: "2",
    title: "Work the CRM",
    detail: "Jump into forms, tables, imports, and settings.",
    href: "/tables",
    icon: LogIn,
  },
  {
    step: "3",
    title: "How it works",
    detail: "Architecture, import/export, and RBAC explained visually.",
    href: "/help",
    icon: BookOpen,
  },
];

export default function WorkspaceBridge() {
  const { token, loading } = useAuth();
  const signedIn = !loading && Boolean(token);
  const steps = signedIn ? stepsSignedIn : stepsSignedOut;

  return (
    <section className="border-t border-slate-200 bg-white">
      <div className="container-width section-padding !py-16 md:!py-20">
        <div className="flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
          <div className="max-w-2xl">
            <p className="text-sm font-semibold uppercase tracking-wider text-emerald-700">
              From home to workspace
            </p>
            <h2 className="mt-2 text-3xl font-black tracking-tight text-slate-900 md:text-4xl">
              {signedIn
                ? "Continue where you left off"
                : "Get into the dashboard in under a minute"}
            </h2>
            <p className="mt-3 text-base leading-7 text-slate-600">
              {signedIn
                ? "You are already signed in. Open the dashboard or explore a module — Help stays one click away."
                : "Sign in once, and every visit from this page can take you straight to your live CRM workspace."}
            </p>
          </div>

          <Link
            href={signedIn ? "/dashboard" : "/login"}
            className="primary-btn shrink-0 self-start lg:self-auto"
          >
            {signedIn ? (
              <>
                <LayoutDashboard size={18} />
                Open dashboard
              </>
            ) : (
              <>
                Enter dashboard
                <ArrowRight size={18} />
              </>
            )}
          </Link>
        </div>

        <ol className="mt-10 grid gap-4 md:grid-cols-3">
          {steps.map((item) => {
            const Icon = item.icon;
            return (
              <li key={item.title}>
                <Link
                  href={item.href}
                  className="flex h-full flex-col rounded-3xl border border-slate-200 bg-slate-50/80 p-6 transition hover:border-emerald-200 hover:bg-emerald-50/40 hover:shadow-sm"
                >
                  <div className="flex items-center justify-between">
                    <span className="flex h-10 w-10 items-center justify-center rounded-2xl bg-emerald-500 text-sm font-bold text-white">
                      {item.step}
                    </span>
                    <Icon size={20} className="text-emerald-600" />
                  </div>
                  <h3 className="mt-5 text-lg font-bold text-slate-900">
                    {item.title}
                  </h3>
                  <p className="mt-2 flex-1 text-sm leading-6 text-slate-600">
                    {item.detail}
                  </p>
                  <span className="mt-4 inline-flex items-center gap-1 text-sm font-semibold text-emerald-700">
                    Continue
                    <ArrowRight size={14} />
                  </span>
                </Link>
              </li>
            );
          })}
        </ol>
      </div>
    </section>
  );
}
