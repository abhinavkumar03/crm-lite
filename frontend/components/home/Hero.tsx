"use client";

import Link from "next/link";
import {
  ArrowRight,
  BookOpen,
  CheckCircle2,
  LayoutDashboard,
  MonitorPlay,
} from "lucide-react";

import { useAuth } from "@/context/AuthContext";

import DashboardPreview from "./DashboardPreview";

const badges = [
  "Go Backend",
  "PostgreSQL",
  "Redis Cache",
  "JWT Auth",
];

export default function Hero() {
  const { token, user, loading } = useAuth();
  const signedIn = !loading && Boolean(token);
  const previewHref = signedIn ? "/dashboard" : "/login";

  return (
    <section className="relative overflow-hidden grid-background">
      <div className="hero-gradient" />

      <div className="container-width p-5 relative z-10">
        <div className="grid items-center gap-16 lg:grid-cols-2 lg:gap-20">
          <div>

            <h1 className="text-5xl font-black leading-tight tracking-tight text-slate-900 md:text-6xl xl:text-7xl">
              Modern CRM
              <br />
              Built Like
              <br />
              <span className="text-gradient">Production Software.</span>
            </h1>

            <p className="mt-8 max-w-xl text-lg leading-8 text-slate-600">
              Go, PostgreSQL, Redis, and Next.js — with metadata-driven modules,
              background jobs, RBAC, and a real dashboard you can open in one
              click.
            </p>

            <div className="mt-10 flex flex-wrap gap-3">
              {signedIn ? (
                <Link href="/dashboard" className="primary-btn">
                  <LayoutDashboard size={18} />
                  Open dashboard
                  <ArrowRight size={18} />
                </Link>
              ) : (
                <Link href="/login" className="primary-btn">
                  Enter dashboard
                  <ArrowRight size={18} />
                </Link>
              )}

              <Link href="/help" className="secondary-btn">
                <BookOpen size={18} />
                How it works
              </Link>

              <Link
                href="/crm_lite_walkthrough.html"
                target="_blank"
                rel="noopener noreferrer"
                className="secondary-btn"
              >
                <MonitorPlay size={18} />
                Watch walkthrough
              </Link>
            </div>

            {!signedIn && !loading && (
              <p className="mt-5 text-sm text-slate-500">
                Demo:{" "}
                <span className="font-medium text-slate-700">
                  admin@crmlite.com
                </span>{" "}
                /{" "}
                <span className="font-medium text-slate-700">Admin@12345</span>
                {" · "}
                then you land on the dashboard.
              </p>
            )}

            <div className="mt-10 flex flex-wrap gap-3">
              {badges.map((badge) => (
                <div key={badge} className="badge">
                  <CheckCircle2 size={16} />
                  {badge}
                </div>
              ))}
            </div>

            <div className="mt-12 grid grid-cols-3 gap-6 sm:gap-8">
              <div>
                <h3 className="text-3xl font-black text-slate-900 sm:text-4xl">
                  60+
                </h3>
                <p className="mt-1 text-sm text-slate-500">API routes</p>
              </div>
              <div>
                <h3 className="text-3xl font-black text-slate-900 sm:text-4xl">
                  Live
                </h3>
                <p className="mt-1 text-sm text-slate-500">Dashboard</p>
              </div>
              <div>
                <h3 className="text-3xl font-black text-slate-900 sm:text-4xl">
                  JWT
                </h3>
                <p className="mt-1 text-sm text-slate-500">Secure auth</p>
              </div>
            </div>
          </div>

          <DashboardPreview
            href={previewHref}
            ctaLabel={signedIn ? "Open live dashboard" : "Enter the dashboard"}
          />
        </div>
      </div>
    </section>
  );
}
