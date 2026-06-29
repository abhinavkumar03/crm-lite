"use client";

import Link from "next/link";
import { useState } from "react";
import {
  Menu,
  X,
  ArrowRight,
  LayoutDashboard,
  Cpu,
  Activity,
  ExternalLink,
} from "lucide-react";
import NavbarActions from "./NavbarActions";
import NavbarMobile from "./NavbarMobile";

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

export default function Navbar() {
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <>
      <header className="sticky top-0 z-50 border-b border-slate-200/70 bg-white/80 backdrop-blur-xl">
        <div className="container-width">
          <div className="flex h-[72px] items-center justify-between">
            {/* Logo */}
            <Link
              href="/"
              className="group flex items-center gap-3"
            >
              <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-gradient-to-br from-emerald-500 to-teal-500 text-lg font-bold text-white shadow-md transition-transform duration-300 group-hover:scale-105">
                C
              </div>

              <div>
                <h2 className="text-lg font-bold tracking-tight text-slate-900">
                  CRM Lite
                </h2>

                <p className="text-xs text-slate-500">
                  Production CRM
                </p>
              </div>
            </Link>

            {/* Desktop Navigation */}
            <nav className="hidden items-center gap-2 lg:flex">
              {navigation?.map((item) => {
                const Icon = item.icon;

                return (
                  <a
                    key={item.name}
                    href={item.href}
                    className="flex items-center gap-2 rounded-full px-4 py-2 text-sm font-medium text-slate-600 transition-all duration-300 hover:bg-slate-100 hover:text-slate-900"
                  >
                    <Icon size={16} />
                    {item.name}
                  </a>
                );
              })}
            </nav>

            {/* Desktop Actions */}
            <div className="hidden lg:block">
                <NavbarActions />
            </div>

            {/* Mobile Button */}
            <button
              onClick={() =>
                setMobileOpen(true)
              }
              className="
              rounded-xl
              border
              border-slate-200
              bg-white
              p-2
              transition
              hover:bg-slate-100
              lg:hidden
              "
            >
              <Menu size={22} />
            </button>
          </div>
        </div>
      </header>

      {/* Mobile Navigation */}
      <NavbarMobile
        open={mobileOpen}
        onClose={() =>
          setMobileOpen(false)
        }
      />
    </>
  );
}