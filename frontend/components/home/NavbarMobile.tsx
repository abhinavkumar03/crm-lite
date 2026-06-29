"use client";

import Link from "next/link";

import {
  X,
  LayoutDashboard,
  Cpu,
  Activity,
} from "lucide-react";

import NavbarActions from "./NavbarActions";

const navigation = [
  {
    name: "Features",
    href: "#features",
    icon: LayoutDashboard,
  },
  {
    name: "Architecture",
    href: "#architecture",
    icon: Cpu,
  },
  {
    name: "Performance",
    href: "#performance",
    icon: Activity,
  },
];

type Props = {
  open: boolean;
  onClose: () => void;
};

export default function NavbarMobile({
  open,
  onClose,
}: Props) {
  return (
    <>
      {/* Overlay */}

      <div
        onClick={onClose}
        className={`
          fixed
          inset-0
          z-40
          bg-black/40
          backdrop-blur-sm
          transition-opacity
          duration-300

          ${
            open
              ? "opacity-100"
              : "pointer-events-none opacity-0"
          }
        `}
      />

      {/* Drawer */}

      <aside
        className={`
          fixed
          left-0
          top-0
          z-50
          flex
          h-screen
          w-[320px]
          max-w-[85vw]
          flex-col
          bg-white
          shadow-2xl
          transition-transform
          duration-300

          ${
            open
              ? "translate-x-0"
              : "-translate-x-full"
          }
        `}
      >
        {/* Header */}

        <div className="flex items-center justify-between border-b border-slate-200 p-6">
          <Link
            href="/"
            onClick={onClose}
            className="flex items-center gap-3"
          >
            <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-gradient-to-br from-emerald-500 to-teal-500 text-lg font-bold text-white">
              C
            </div>

            <div>
              <h2 className="font-bold text-slate-900">
                CRM Lite
              </h2>

              <p className="text-xs text-slate-500">
                Production CRM
              </p>
            </div>
          </Link>

          <button
            onClick={onClose}
            className="rounded-xl p-2 transition hover:bg-slate-100"
          >
            <X size={22} />
          </button>
        </div>

        {/* Navigation */}

        <div className="flex-1 overflow-y-auto px-5 py-6">
          <p className="mb-4 px-3 text-xs font-semibold uppercase tracking-widest text-slate-400">
            Navigation
          </p>

          <div className="space-y-2">
            {navigation.map((item) => {
              const Icon = item.icon;

              return (
                <a
                  key={item.name}
                  href={item.href}
                  onClick={onClose}
                  className="
                  flex
                  items-center
                  gap-3
                  rounded-2xl
                  px-4
                  py-3
                  text-slate-700
                  transition
                  hover:bg-slate-100
                  "
                >
                  <Icon size={20} />

                  {item.name}
                </a>
              );
            })}
          </div>
        </div>

        {/* Footer */}

        <div className="border-t border-slate-200 p-5">
          <NavbarActions
            mobile
            onNavigate={onClose}
          />
        </div>
      </aside>
    </>
  );
}