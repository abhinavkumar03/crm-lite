"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import {
  SlidersHorizontal,
  Boxes,
  FormInput,
  ShieldCheck,
  Zap,
  Database,
  LucideIcon,
} from "lucide-react";

type NavItem = {
  name: string;
  href: string;
  icon: LucideIcon;
  description: string;
};

const items: NavItem[] = [
  {
    name: "General",
    href: "/settings",
    icon: SlidersHorizontal,
    description: "Organization profile & preferences",
  },
  {
    name: "Modules",
    href: "/settings/modules",
    icon: Boxes,
    description: "Create & manage object types",
  },
  {
    name: "Fields",
    href: "/settings/fields",
    icon: FormInput,
    description: "Define fields per module",
  },
  {
    name: "Validation",
    href: "/settings/validation",
    icon: ShieldCheck,
    description: "Data-quality rules",
  },
  {
    name: "Automation",
    href: "/settings/automation",
    icon: Zap,
    description: "Notification behaviour",
  },
  {
    name: "Data",
    href: "/settings/data",
    icon: Database,
    description: "Import & export",
  },
];

export default function SettingsNav() {
  const pathname = usePathname();

  return (
    <nav className="flex flex-col gap-1">
      {items.map((item) => {
        const Icon = item.icon;
        const active =
          item.href === "/settings"
            ? pathname === "/settings"
            : pathname === item.href || pathname.startsWith(`${item.href}/`);

        return (
          <Link
            key={item.href}
            href={item.href}
            className={`group flex items-start gap-3 rounded-2xl border px-4 py-3 transition ${
              active
                ? "border-emerald-200 bg-emerald-50"
                : "border-transparent hover:bg-slate-100"
            }`}
          >
            <span
              className={`mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-xl ${
                active
                  ? "bg-emerald-500 text-white"
                  : "bg-slate-100 text-slate-500 group-hover:bg-white"
              }`}
            >
              <Icon size={16} />
            </span>

            <span className="min-w-0">
              <span
                className={`block text-sm font-semibold ${
                  active ? "text-emerald-700" : "text-slate-800"
                }`}
              >
                {item.name}
              </span>
              <span className="block truncate text-xs text-slate-500">
                {item.description}
              </span>
            </span>
          </Link>
        );
      })}
    </nav>
  );
}
