import { ReactNode } from "react";

import Link from "next/link";

import {
  ShieldCheck,
  Database,
  Activity,
} from "lucide-react";

type Props = {
  title: string;
  subtitle: string;
  children: ReactNode;
};

const features = [
  {
    icon: ShieldCheck,
    title: "Secure Authentication",
    description:
      "JWT based authentication with protected routes.",
  },
  {
    icon: Database,
    title: "Production Backend",
    description:
      "Go, PostgreSQL and Redis powered architecture.",
  },
  {
    icon: Activity,
    title: "Modern CRM",
    description:
      "Built for recruiters, engineers and students.",
  },
];

export default function AuthLayout({
  title,
  subtitle,
  children,
}: Props) {
  return (
    <div className="min-h-screen bg-slate-50">
      <div className="grid min-h-screen lg:grid-cols-2">
        {/* Left Panel */}

        <div
          className="
          hidden
          lg:flex
          flex-col
          justify-between
          bg-gradient-to-br
          from-emerald-600
          via-teal-600
          to-cyan-700
          p-12
          text-white
          "
        >
          {/* Logo */}

          <Link
            href="/"
            className="flex items-center gap-4"
          >
            <div
              className="
              flex
              h-14
              w-14
              items-center
              justify-center
              rounded-2xl
              bg-white/20
              text-2xl
              font-bold
              backdrop-blur
              "
            >
              C
            </div>

            <div>
              <h1 className="text-2xl font-bold">
                CRM Lite
              </h1>

              <p className="text-sm text-emerald-100">
                Modern CRM Platform
              </p>
            </div>
          </Link>

          {/* Hero */}

          <div className="max-w-lg">
            <h2 className="text-5xl font-bold leading-tight">
              Manage customers.
              <br />
              Close deals.
              <br />
              Grow faster.
            </h2>

            <p className="mt-6 text-lg leading-8 text-emerald-100">
              CRM Lite helps you organize leads,
              manage contacts and track tasks with
              a modern, production-ready workflow.
            </p>
          </div>

          {/* Features */}

          <div className="space-y-6">
            {features.map((feature) => {
              const Icon = feature.icon;

              return (
                <div
                  key={feature.title}
                  className="flex gap-4"
                >
                  <div
                    className="
                    flex
                    h-12
                    w-12
                    items-center
                    justify-center
                    rounded-xl
                    bg-white/15
                    "
                  >
                    <Icon size={22} />
                  </div>

                  <div>
                    <h3 className="font-semibold">
                      {feature.title}
                    </h3>

                    <p className="mt-1 text-sm text-emerald-100">
                      {feature.description}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Right Panel */}

        <div
          className="
          flex
          items-center
          justify-center
          px-6
          py-12
          lg:px-16
          "
        >
          <div className="w-full max-w-md">
            {/* Mobile Logo */}

            <Link
              href="/"
              className="mb-10 flex items-center gap-4 lg:hidden"
            >
              <div
                className="
                flex
                h-12
                w-12
                items-center
                justify-center
                rounded-2xl
                bg-gradient-to-br
                from-emerald-500
                to-teal-500
                text-xl
                font-bold
                text-white
                "
              >
                C
              </div>

              <div>
                <h2 className="text-xl font-bold text-slate-900">
                  CRM Lite
                </h2>

                <p className="text-sm text-slate-500">
                  Modern CRM
                </p>
              </div>
            </Link>

            {/* Header */}

            <div className="mb-8">
              <h1 className="text-3xl font-bold text-slate-900">
                {title}
              </h1>

              <p className="mt-2 text-slate-500">
                {subtitle}
              </p>
            </div>

            {children}
          </div>
        </div>
      </div>
    </div>
  );
}