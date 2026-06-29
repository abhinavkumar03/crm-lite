"use client";

import { useRouter } from "next/navigation";

import {
  LogOut,
  ChevronDown,
} from "lucide-react";

import { useAuth } from "@/context/AuthContext";

export default function UserMenu() {
  const auth = useAuth();

  const router = useRouter();

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

  return (
    <div className="flex items-center gap-3">
      <div
        className="
        flex
        h-11
        w-11
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

      <div className="hidden xl:block">
        <h4 className="font-semibold text-slate-900">
          {auth.user?.name}
        </h4>

        <p className="text-sm text-slate-500">
          {auth.user?.email}
        </p>
      </div>

      <ChevronDown
        size={18}
        className="hidden text-slate-400 xl:block"
      />

      <button
        onClick={logout}
        className="
        ml-2
        rounded-xl
        border
        border-red-200
        p-2.5
        text-red-500
        transition
        hover:bg-red-50
        "
      >
        <LogOut size={18} />
      </button>
    </div>
  );
}