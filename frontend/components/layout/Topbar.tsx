"use client";

import { Menu } from "lucide-react";

import SearchBar from "./SearchBar";
import NotificationBell from "./NotificationBell";
import UserMenu from "./UserMenu";

export default function Topbar() {
  return (
    <header
      className="
      sticky
      top-0
      z-30
      flex
      h-20
      items-center
      justify-between
      border-b
      border-slate-200
      bg-white/80
      px-8
      backdrop-blur-xl
      "
    >
      {/* Left */}

      <div className="flex items-center gap-4">
        <button
          className="
          rounded-xl
          border
          border-slate-200
          p-2
          lg:hidden
          "
        >
          <Menu size={20} />
        </button>

        <SearchBar />
      </div>

      {/* Right */}

      <div className="flex items-center gap-4">
        <NotificationBell />

        <UserMenu />
      </div>
    </header>
  );
}