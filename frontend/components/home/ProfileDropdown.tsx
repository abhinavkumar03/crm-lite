"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";

import {
  ChevronDown,
  LayoutDashboard,
  LogOut,
  Settings,
  UserCircle2,
} from "lucide-react";

import { useAuth } from "@/context/AuthContext";

type Props = {
  mobile?: boolean;
};

export default function ProfileDropdown({
  mobile = false,
}: Props) {
  const auth = useAuth();

  const router = useRouter();

  const [open, setOpen] = useState(false);

  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (
        ref.current &&
        !ref.current.contains(e.target as Node)
      ) {
        setOpen(false);
      }
    }

    document.addEventListener(
      "mousedown",
      handleClick
    );

    return () =>
      document.removeEventListener(
        "mousedown",
        handleClick
      );
  }, []);

  function logout() {
    auth.logout();
    router.replace("/login");
  }

  const initials =
    auth.user?.name
      ?.split(" ")
      .map((n) => n[0])
      .join("")
      .substring(0, 2)
      .toUpperCase() ?? "U";

  // Mobile Version
  if (mobile) {
    return (
      <button
        onClick={logout}
        className="
        flex
        w-full
        items-center
        justify-center
        rounded-2xl
        border
        border-red-200
        px-5
        py-3
        font-medium
        text-red-500
        transition
        hover:bg-red-50
        "
      >
        <LogOut
          size={18}
          className="mr-2"
        />

        Logout
      </button>
    );
  }

  return (
    <div
      ref={ref}
      className="relative"
    >
      {/* Trigger */}

      <button
        onClick={() =>
          setOpen(!open)
        }
        className="
        flex
        items-center
        gap-3
        rounded-2xl
        border
        border-slate-200
        bg-white
        px-3
        py-2
        shadow-sm
        transition
        hover:border-emerald-300
        hover:shadow-md
        "
      >
        <div
          className="
          flex
          h-10
          w-10
          items-center
          justify-center
          rounded-full
          bg-gradient-to-br
          from-emerald-500
          to-teal-500
          font-semibold
          text-white
          "
        >
          {initials}
        </div>

        <div className="hidden text-left xl:block">
          <h4 className="max-w-[150px] truncate text-sm font-semibold text-slate-900">
            {auth.user?.name}
          </h4>

          <p className="max-w-[150px] truncate text-xs text-slate-500">
            {auth.user?.email}
          </p>
        </div>

        <ChevronDown
          size={18}
          className={`text-slate-400 transition ${
            open ? "rotate-180" : ""
          }`}
        />
      </button>

      {/* Dropdown */}

      {open && (
        <div
          className="
          absolute
          right-0
          mt-3
          w-72
          overflow-hidden
          rounded-3xl
          border
          border-slate-200
          bg-white
          shadow-2xl
          "
        >
          {/* Header */}

          <div className="border-b border-slate-100 p-5">
            <div className="flex items-center gap-4">
              <div
                className="
                flex
                h-14
                w-14
                items-center
                justify-center
                rounded-full
                bg-gradient-to-br
                from-emerald-500
                to-teal-500
                text-lg
                font-bold
                text-white
                "
              >
                {initials}
              </div>

              <div className="min-w-0">
                <h3 className="truncate font-semibold text-slate-900">
                  {auth.user?.name}
                </h3>

                <p className="truncate text-sm text-slate-500">
                  {auth.user?.email}
                </p>
              </div>
            </div>
          </div>

          {/* Navigation */}

          <div className="p-2">
            <Link
              href="/dashboard"
              onClick={() =>
                setOpen(false)
              }
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
              <LayoutDashboard
                size={18}
              />
              Dashboard
            </Link>

            <button
              className="
              flex
              w-full
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
              <UserCircle2
                size={18}
              />
              Profile
            </button>

            <button
              className="
              flex
              w-full
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
              <Settings
                size={18}
              />
              Settings
            </button>
          </div>

          {/* Footer */}

          <div className="border-t border-slate-100 p-2">
            <button
              onClick={logout}
              className="
              flex
              w-full
              items-center
              gap-3
              rounded-2xl
              px-4
              py-3
              text-red-500
              transition
              hover:bg-red-50
              "
            >
              <LogOut
                size={18}
              />
              Logout
            </button>
          </div>
        </div>
      )}
    </div>
  );
}