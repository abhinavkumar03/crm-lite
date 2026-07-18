"use client";

import Link from "next/link";
import { Home, Menu } from "lucide-react";

import SearchBar from "./SearchBar";
import NotificationBell from "./NotificationBell";
import UserMenu from "./UserMenu";
import MobileSearch from "@/features/search/components/MobileSearch";

type Props = {
  onMenuClick: () => void;
};

export default function Topbar({
  onMenuClick,
}: Props) {
  return (
    <header
      className="
      sticky
      top-0
      z-30
      border-b
      border-slate-200
      bg-white/80
      backdrop-blur-xl
      "
    >
      <div
        className="
        flex
        h-20
        items-center
        justify-between
        gap-4
        px-4
        lg:px-8
        "
      >
        {/* Left */}

        <div
          className="
          flex
          min-w-0
          flex-1
          items-center
          gap-4
          "
        >
          {/* Mobile Menu */}

          <button
            onClick={onMenuClick}
            className="
            flex
            h-11
            w-11
            items-center
            justify-center
            rounded-2xl
            border
            border-slate-200
            bg-white
            transition

            hover:bg-slate-100

            lg:hidden
            "
          >
            <Menu size={20} />
          </button>

          {/* Desktop Search */}

          <div
            data-tour="global-search"
            className="hidden flex-1 lg:block"
          >
            <SearchBar />
          </div>

          {/* Mobile Search */}

          <div className="flex-1 lg:hidden">
            <MobileSearch />
          </div>
        </div>

        {/* Right */}

        <div
          className="
          flex
          flex-shrink-0
          items-center
          gap-3
          "
        >
          <Link
            href="/"
            aria-label="Home"
            className="
            rounded-2xl
            bg-white
            p-3
            transition
            hover:bg-slate-100
            "
          >
            <Home size={20} className="text-slate-700" />
          </Link>

          <div data-tour="notification-bell">
            <NotificationBell />
          </div>

          <div data-tour="user-menu">
            <UserMenu />
          </div>
        </div>
      </div>
    </header>
  );
}