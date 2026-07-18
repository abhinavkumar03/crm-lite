"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { usePathname, useSearchParams } from "next/navigation";
import {
  X,
  LayoutDashboard,
  LayoutTemplate,
  Table2,
  Upload,
  Download,
  MessageCircle,
  Settings,
  CircleHelp,
  ChevronRight,
  Boxes,
  type LucideIcon,
} from "lucide-react";

import api from "@/services/api";

type NavItem = {
  name: string;
  href: string;
  icon: LucideIcon;
  tourKey: string;
};

type NavModule = {
  id: string;
  api_name: string;
  plural_label: string;
};

const fixedNavigation: NavItem[] = [
  {
    name: "Dashboard",
    href: "/dashboard",
    icon: LayoutDashboard,
    tourKey: "dashboard",
  },
  {
    name: "Forms",
    href: "/forms",
    icon: LayoutTemplate,
    tourKey: "forms",
  },
  {
    name: "Tables",
    href: "/tables",
    icon: Table2,
    tourKey: "tables",
  },
  {
    name: "Import",
    href: "/imports",
    icon: Upload,
    tourKey: "imports",
  },
  {
    name: "Export",
    href: "/exports",
    icon: Download,
    tourKey: "exports",
  },
  {
    name: "Notifications",
    href: "/notifications",
    icon: MessageCircle,
    tourKey: "notifications",
  },
  {
    name: "Settings",
    href: "/settings",
    icon: Settings,
    tourKey: "settings",
  },
  {
    name: "How it works",
    href: "/help",
    icon: CircleHelp,
    tourKey: "help",
  },
];

type Props = {
  open: boolean;
  onClose: () => void;
};

function NavLink({
  name,
  href,
  icon: Icon,
  tourKey,
  active,
  onClose,
}: {
  name: string;
  href: string;
  icon: LucideIcon;
  tourKey: string;
  active: boolean;
  onClose: () => void;
}) {
  return (
    <Link
      href={href}
      data-tour={`nav-${tourKey}`}
      onClick={onClose}
      className={`
      group
      flex
      items-center
      justify-between
      rounded-2xl
      px-4
      py-3
      transition-all
      duration-200

      ${
        active
          ? "bg-emerald-500 text-white shadow-lg"
          : "text-slate-600 hover:bg-slate-100"
      }
    `}
    >
      <div className="flex items-center gap-3">
        <Icon size={20} />
        <span className="font-medium">{name}</span>
      </div>
      <ChevronRight
        size={16}
        className={`
        transition-transform
        ${
          active
            ? "translate-x-1"
            : "opacity-0 group-hover:translate-x-1 group-hover:opacity-100"
        }
      `}
      />
    </Link>
  );
}

export default function Sidebar({ open, onClose }: Props) {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const activeModuleId = searchParams.get("module");
  const [modules, setModules] = useState<NavModule[]>([]);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const res = await api.get("/navigation");
        if (active) setModules(res.data.data ?? []);
      } catch {
        if (active) setModules([]);
      }
    })();
    return () => {
      active = false;
    };
  }, []);

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

        ${open ? "translate-x-0" : "-translate-x-full"}
      `}
      >
        <div className="relative border-b border-slate-200 px-6 py-4">
          <button
            onClick={onClose}
            className="
            absolute
            right-4
            top-4
            rounded-xl
            p-2
            transition
            hover:bg-slate-100
            lg:hidden
            "
          >
            <X size={20} />
          </button>

          <Link
            href="/dashboard"
            onClick={onClose}
            className="flex items-center gap-4"
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
              shadow
              "
            >
              C
            </div>
            <div>
              <h2 className="text-lg font-bold">CRM Lite</h2>
              <p className="text-xs text-slate-500">Production CRM</p>
            </div>
          </Link>
        </div>

        <div className="flex-1 overflow-y-auto px-4 py-6">
          <p className="mb-4 px-3 text-xs font-semibold uppercase tracking-widest text-slate-400">
            Workspace
          </p>

          <nav data-tour="sidebar-nav" className="space-y-2">
            {fixedNavigation.map((item) => {
              const active =
                pathname === item.href ||
                pathname.startsWith(`${item.href}/`);
              return (
                <NavLink
                  key={item.href}
                  {...item}
                  active={active}
                  onClose={onClose}
                />
              );
            })}
          </nav>

          {modules.length > 0 && (
            <>
              <p className="mb-4 mt-8 px-3 text-xs font-semibold uppercase tracking-widest text-slate-400">
                Modules
              </p>
              <nav className="space-y-2">
                {modules.map((m) => {
                  const href = `/tables?module=${m.id}`;
                  const active =
                    pathname.startsWith("/tables") &&
                    activeModuleId === m.id;
                  return (
                    <NavLink
                      key={m.id}
                      name={m.plural_label}
                      href={href}
                      icon={Boxes}
                      tourKey={`module-${m.api_name}`}
                      active={active}
                      onClose={onClose}
                    />
                  );
                })}
              </nav>
            </>
          )}
        </div>
      </aside>
    </>
  );
}
