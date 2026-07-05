import Link from "next/link";
import {
  ArrowRight,
  PlayCircle,
  CheckCircle2,
  MonitorPlay
} from "lucide-react";

import DashboardPreview from "./DashboardPreview";

const badges = [
  "Go Backend",
  "PostgreSQL",
  "Redis Cache",
  "JWT Auth",
];

export default function Hero() {
  return (
    <section className="relative overflow-hidden grid-background">
      <div className="hero-gradient" />

      <div className="container-width section-padding relative z-10">
        <div className="grid items-center gap-20 lg:grid-cols-2">
          {/* Left */}

          <div>
            <div className="mb-6 inline-flex items-center rounded-full border border-emerald-200 bg-emerald-50 px-4 py-2 text-sm font-semibold text-emerald-700">
              🚀 Production Ready CRM Portfolio
            </div>

            <h1 className="text-5xl font-black leading-tight tracking-tight text-slate-900 md:text-6xl xl:text-7xl">
              Modern CRM
              <br />
              Built Like
              <br />
              <span className="text-gradient">
                Production Software.
              </span>
            </h1>

            <p className="mt-8 max-w-2xl text-lg leading-8 text-slate-600">
              A production-ready CRM built using Go,
              PostgreSQL, Redis, Docker and Next.js.

              <br />
              <br />

              Designed for recruiters, hiring managers
              and backend engineers to showcase
              scalable software engineering rather than
              simple CRUD operations.
            </p>

            {/* Buttons */}

            <div className="mt-10 flex flex-wrap gap-4">
              <Link
                href="/login"
                className="primary-btn"
              >
                Launch CRM

                <ArrowRight size={18} />
              </Link>

              {/* <Link
                href="#architecture"
                className="secondary-btn"
              >
                <PlayCircle size={18} />

                Explore Architecture
              </Link> */}

              <Link
                href="/crm_lite_walkthrough.html"
                target="_blank"
                rel="noopener noreferrer"
                className="secondary-btn"
              >
                <MonitorPlay size={18} />

                Product Walkthrough
              </Link>
            </div>

            {/* Badges */}

            <div className="mt-10 flex flex-wrap gap-3">
              {badges?.map((badge) => (
                <div
                  key={badge}
                  className="badge"
                >
                  <CheckCircle2 size={16} />

                  {badge}
                </div>
              ))}
            </div>

            {/* Stats */}

            <div className="mt-14 grid grid-cols-3 gap-8">
              <div>
                <h3 className="text-4xl font-black text-slate-900">
                  30+
                </h3>

                <p className="mt-2 text-sm text-slate-500">
                  REST APIs
                </p>
              </div>

              <div>
                <h3 className="text-4xl font-black text-slate-900">
                  Redis
                </h3>

                <p className="mt-2 text-sm text-slate-500">
                  Cached Queries
                </p>
              </div>

              <div>
                <h3 className="text-4xl font-black text-slate-900">
                  JWT
                </h3>

                <p className="mt-2 text-sm text-slate-500">
                  Secure Authentication
                </p>
              </div>
            </div>
          </div>

          {/* Right */}

          <DashboardPreview />
        </div>
      </div>
    </section>
  );
}