"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  X,
  LayoutDashboard,
  MessageCircle,
  Settings,
  CircleHelp,
  ChevronRight,
  Boxes,
  Compass,
  type LucideIcon,
} from "lucide-react";

import api from "@/services/api";
import { useDemo } from "@/features/demo/DemoProvider";
import WorkspaceSwitcher from "@/features/organization/components/WorkspaceSwitcher";
import { useOrgOptional } from "@/features/organization/OrgContext";

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

/** Fixed workspace items (modules are inserted after Dashboard). */
const beforeModules: NavItem[] = [
  {
    name: "Dashboard",
    href: "/dashboard",
    icon: LayoutDashboard,
    tourKey: "dashboard",
  },
];

const afterModules: NavItem[] = [
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

function ExploreCrmButton({ onClose }: { onClose: () => void }) {
  const demo = useDemo();
  if (!demo) return null;
  return (
    <button
      type="button"
      data-demo="sidebar-explore"
      onClick={() => {
        onClose();
        demo.openLauncher();
      }}
      className="
        mt-6 flex w-full items-center gap-3 rounded-2xl
        border border-emerald-200 bg-emerald-50 px-4 py-3
        text-left text-sm font-semibold text-emerald-800
        transition hover:bg-emerald-100
      "
    >
      <Compass size={18} />
      <span>
        Explore CRM
        <span className="mt-0.5 block text-xs font-normal text-emerald-700/80">
          Interactive sandbox tutorial
        </span>
      </span>
    </button>
  );
}

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
  const [modules, setModules] = useState<NavModule[]>([]);
  const orgCtx = useOrgOptional();

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
  }, [orgCtx?.activeOrg?.id]);

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
        <div className="relative border-b border-slate-200 px-4 py-4">
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

          <WorkspaceSwitcher variant="sidebar" className="pr-8 lg:pr-0" />
        </div>

        <div className="flex-1 overflow-y-auto px-4 py-6">
          <p className="mb-4 px-3 text-xs font-semibold uppercase tracking-widest text-slate-400">
            Navigation
          </p>

          <nav data-tour="sidebar-nav" className="space-y-2">
            {beforeModules.map((item) => {
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

            {modules.map((m) => {
              const href = `/m/${m.api_name}`;
              const active =
                pathname === href ||
                pathname.startsWith(`${href}/`);
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

            {afterModules.map((item) => {
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

          <ExploreCrmButton onClose={onClose} />
        </div>
      </aside>
    </>
  );
}
