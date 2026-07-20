"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import {
  SlidersHorizontal,
  Boxes,
  FormInput,
  ShieldCheck,
  Zap,
  Upload,
  Download,
  LayoutTemplate,
  Table2,
  Shield,
  Users,
  Building2,
  Network,
  MessageSquareText,
  Radio,
  LucideIcon,
} from "lucide-react";

type NavItem = {
  name: string;
  href: string;
  icon: LucideIcon;
  description: string;
  tutorialAction?: string;
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
    tutorialAction: "open-modules",
  },
  {
    name: "Fields",
    href: "/settings/fields",
    icon: FormInput,
    description: "Define fields per module",
    tutorialAction: "open-fields",
  },
  {
    name: "Form Designer",
    href: "/settings/forms",
    icon: LayoutTemplate,
    description: "Preview metadata-driven forms",
    tutorialAction: "open-forms",
  },
  {
    name: "Listing Columns",
    href: "/settings/tables",
    icon: Table2,
    description: "Table column visibility & order",
  },
  {
    name: "Validation",
    href: "/settings/validation",
    icon: ShieldCheck,
    description: "Data-quality rules",
    tutorialAction: "open-validation",
  },
  {
    name: "Import",
    href: "/settings/imports",
    icon: Upload,
    description: "CSV / Excel into modules",
    tutorialAction: "open-imports",
  },
  {
    name: "Export",
    href: "/settings/exports",
    icon: Download,
    description: "CSV / Excel out of modules",
    tutorialAction: "open-exports",
  },
  {
    name: "Members",
    href: "/settings/members",
    icon: Users,
    description: "Invite & manage people",
  },
  {
    name: "Roles",
    href: "/settings/roles",
    icon: Shield,
    description: "Permission matrix & ACL",
    tutorialAction: "open-roles",
  },
  {
    name: "Departments",
    href: "/settings/departments",
    icon: Building2,
    description: "Org structure departments",
  },
  {
    name: "Teams",
    href: "/settings/teams",
    icon: Network,
    description: "Teams under departments",
  },
  {
    name: "Automation",
    href: "/settings/automation",
    icon: Zap,
    description: "Notification behaviour",
    tutorialAction: "open-automation",
  },
  {
    name: "Providers",
    href: "/settings/communication-providers",
    icon: Radio,
    description: "Email & WhatsApp delivery",
  },
  {
    name: "Message Templates",
    href: "/settings/notification-templates",
    icon: MessageSquareText,
    description: "Email & WhatsApp templates",
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
            data-tutorial-action={item.tutorialAction}
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
