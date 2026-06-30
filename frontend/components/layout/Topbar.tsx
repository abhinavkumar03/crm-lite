"use client";

import { Menu } from "lucide-react";

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
      flex
      h-20
      items-center
      justify-between
      border-b
      border-slate-200
      bg-white/80
      px-4
      backdrop-blur-xl
      lg:px-8
      "
    >
      <div className="flex w-full items-center gap-4">
  <button
    onClick={onMenuClick}
    className="
    rounded-xl
    border
    border-slate-200
    bg-white
    p-2
    transition
    hover:bg-slate-100
    lg:hidden
    "
  >
    <Menu size={20} />
  </button>

  {/* Desktop */}

  <div className="hidden flex-1 lg:block">
    <SearchBar />
  </div>

  {/* Mobile */}

  <div className="flex-1 lg:hidden">
    <MobileSearch />
  </div>
</div>

      <div className="flex items-center gap-3 lg:gap-4">
        <NotificationBell />

        <UserMenu />
      </div>
    </header>
  );
}