"use client";

import Link from "next/link";
import { X } from "lucide-react";
import { usePathname } from "next/navigation";
import {
    LayoutDashboard,
    Users,
    ContactRound,
    CheckSquare,
    Settings,
    ChevronRight,
} from "lucide-react";


const navigation = [
    {
        name: "Dashboard",
        href: "/dashboard",
        icon: LayoutDashboard,
    },
    {
        name: "Leads",
        href: "/leads",
        icon: Users,
    },
    {
        name: "Contacts",
        href: "/contacts",
        icon: ContactRound,
    },
    {
        name: "Tasks",
        href: "/tasks",
        icon: CheckSquare,
    },
];

type Props = {
  open: boolean;
  onClose: () => void;
};

export default function Sidebar({
  open,
  onClose,
}: Props) {
    const pathname = usePathname();
    

    return (
        <>
            {open && (
                <div
                    onClick={onClose}
                    className="
      fixed
      inset-0
      z-40
      bg-black/40
      backdrop-blur-sm
      lg:hidden
      "
                />
            )}
            <aside
                className={`
      fixed
      left-0
      top-0
      z-50
      flex
      h-screen
      w-72
      flex-col
      border-r
      border-slate-200
      bg-white
      transition-transform
      duration-300

      lg:relative
      lg:translate-x-0

      ${open
                        ? "translate-x-0"
                        : "-translate-x-full"
                    }
    `}
            >
                {/* Logo */}

                <div className="relative border-b border-slate-200 p-4">
                {/* Close Button */}

                <button
                    onClick={onClose}
                    className="absolute right-4 top-4 rounded-lg p-2 text-slate-500 transition hover:bg-slate-100 hover:text-slate-900 lg:hidden"
                    aria-label="Close sidebar"
                >
                    <X size={20} />
                </button>

                <Link
                    href="/dashboard"
                    className="flex items-center gap-4"
                >
                    <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-emerald-500 to-teal-500 text-xl font-bold text-white shadow">
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
                </div>

                {/* Navigation */}

                <div className="flex-1 px-4 py-6">
                    <p className="mb-4 px-3 text-xs font-semibold uppercase tracking-widest text-slate-400">
                        Workspace
                    </p>

                    <nav className="space-y-2">
                        {navigation.map((item) => {
                            const Icon = item.icon;

                            const active =
                                pathname === item.href;

                            return (
                                <Link
                                    key={item.href}
                                    href={item.href}
                                    className={`group flex items-center justify-between rounded-2xl px-4 py-3 transition-all duration-300 ${active
                                        ? "bg-emerald-500 text-white shadow-lg"
                                        : "text-slate-600 hover:bg-slate-100"
                                        }`}
                                >
                                    <div className="flex items-center gap-3">
                                        <Icon size={20} />

                                        <span className="font-medium">
                                            {item.name}
                                        </span>
                                    </div>

                                    <ChevronRight
                                        size={16}
                                        className={`transition ${active
                                            ? "translate-x-1"
                                            : "opacity-0 group-hover:translate-x-1 group-hover:opacity-100"
                                            }`}
                                    />
                                </Link>
                            );
                        })}
                    </nav>
                </div>

                {/* Bottom */}

                <div className="border-t border-slate-200 p-4">
                    <button className="flex w-full items-center gap-3 rounded-2xl px-4 py-3 text-slate-600 transition hover:bg-slate-100">
                        <Settings size={20} />

                        <span className="font-medium">
                            Settings
                        </span>
                    </button>

                    <div className="mt-6 rounded-2xl bg-slate-50 p-4">
                        <p className="text-sm font-semibold text-slate-900">
                            CRM Lite
                        </p>

                        <p className="mt-1 text-xs leading-5 text-slate-500">
                            Modern CRM built with Go, PostgreSQL,
                            Redis and Next.js.
                        </p>
                    </div>
                </div>
            </aside>
        </>

    );
}