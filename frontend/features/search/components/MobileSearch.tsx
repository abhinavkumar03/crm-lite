"use client";

import { Search } from "lucide-react";
import { useState } from "react";
import MobileSearchModal from "./MobileSearchModal";


export default function MobileSearch() {
  const [open, setOpen] =
    useState(false);

  return (
    <>
      <button
        onClick={() =>
          setOpen(true)
        }
        className="
        flex
        w-full
        items-center
        gap-3
        rounded-2xl
        border
        border-slate-200
        bg-slate-50
        px-4
        py-3
        text-sm
        text-slate-500
        "
      >
        <Search size={18} />

        Search CRM...
      </button>

      <MobileSearchModal
        open={open}
        onClose={() =>
          setOpen(false)
        }
      />
    </>
  );
}