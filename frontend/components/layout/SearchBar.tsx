"use client";

import { Search } from "lucide-react";

export default function SearchBar() {
  return (
    <div className="relative hidden w-full max-w-md lg:block">
      <Search
        size={18}
        className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400"
      />

      <input
        type="text"
        placeholder="Search leads, contacts or tasks..."
        className="
        w-full
        rounded-2xl
        border
        border-slate-200
        bg-slate-50
        py-3
        pl-11
        pr-4
        text-sm
        outline-none
        transition
        focus:border-emerald-500
        focus:bg-white
        "
        disabled
      />
    </div>
  );
}