"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import {
  X,
  LayoutDashboard,
  Users,
  ContactRound,
  CheckSquare,
  LayoutTemplate,
  Table2,
  MessageCircle,
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
  {
    name: "Dynamic Forms",
    href: "/forms",
    icon: LayoutTemplate,
  },
  {
    name: "Dynamic Tables",
    href: "/tables",
    icon: Table2,
  },
  {
    name: "Notifications",
    href: "/notifications",
    icon: MessageCircle,
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
      {/* Overlay */}

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

        ${
          open
            ? "translate-x-0"
            : "-translate-x-full"
        }
      `}
      >
        {/* Header */}

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
              <h2 className="text-lg font-bold">
                CRM Lite
              </h2>

              <p className="text-xs text-slate-500">
                Production CRM
              </p>
            </div>
          </Link>
        </div>

        {/* Navigation */}

        <div className="flex-1 overflow-y-auto px-4 py-6">
          <p
            className="
            mb-4
            px-3
            text-xs
            font-semibold
            uppercase
            tracking-widest
            text-slate-400
            "
          >
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

                    <span className="font-medium">
                      {item.name}
                    </span>
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
            })}
          </nav>
        </div>

        {/* Footer */}

        {/* <div className="border-slate-200 p-5">
          <div
            className="
            rounded-2xl
            bg-gradient-to-br
            from-slate-50
            to-slate-100
            p-5
            "
          >
            <p className="mt-2 text-sm leading-6 text-slate-500">
              Built with Go, PostgreSQL, Redis and
              Next.js for modern CRM workflows.
            </p>

            <div className="mt-4 flex flex-wrap gap-2">
              <span className="rounded-full bg-white px-3 py-1 text-xs text-slate-600 shadow-sm">
                Go
              </span>

              <span className="rounded-full bg-white px-3 py-1 text-xs text-slate-600 shadow-sm">
                PostgreSQL
              </span>

              <span className="rounded-full bg-white px-3 py-1 text-xs text-slate-600 shadow-sm">
                Redis
              </span>

              <span className="rounded-full bg-white px-3 py-1 text-xs text-slate-600 shadow-sm">
                Next.js
              </span>
            </div>
          </div>
        </div> */}
      </aside>
    </>
  );
}