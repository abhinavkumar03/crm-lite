"use client";

import {
  Search,
  X,
} from "lucide-react";

import { useEffect, useRef, useState } from "react";

import useGlobalSearch from "../hooks/useGlobalSearch";

import SearchDropdown from "./SearchDropdown";

type Props = {
  open: boolean;
  onClose: () => void;
};

export default function MobileSearchModal({
  open,
  onClose,
}: Props) {
  const inputRef =
    useRef<HTMLInputElement>(null);

  const [query, setQuery] =
    useState("");

  const {
    loading,
    results,
  } = useGlobalSearch(query);

  useEffect(() => {
    if (open) {
      setTimeout(() => {
        inputRef.current?.focus();
      }, 100);
    }
  }, [open]);

  if (!open) return null;

  return (
    <div
      className="
      fixed
      inset-0
      z-[100]
      bg-white
      "
    >
      <div className="border-b border-slate-200 p-4">
        <div className="flex items-center gap-3">
          <Search
            size={20}
            className="text-slate-400"
          />

          <input
            ref={inputRef}
            value={query}
            onChange={(e) =>
              setQuery(e.target.value)
            }
            placeholder="Search CRM..."
            className="
            flex-1
            outline-none
            text-lg
            "
          />

          <button
            onClick={onClose}
          >
            <X size={22} />
          </button>
        </div>
      </div>

      <div className="p-4">
        <SearchDropdown
          loading={loading}
          results={results}
          open={query.length >= 2}
          query={query}
          activeIndex={0}
          onClose={onClose}
        />
      </div>
    </div>
  );
}