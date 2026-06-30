"use client";

import {
  useEffect,
  useRef,
  useState,
} from "react";

import { Search, Command } from "lucide-react";

import useGlobalSearch from "../hooks/useGlobalSearch";

import SearchDropdown from "./SearchDropdown";

export default function SearchBar() {
  const [query, setQuery] =
    useState("");

  const [open, setOpen] =
    useState(false);

  const wrapperRef =
    useRef<HTMLDivElement>(null);

  const {
    loading,
    results,
  } = useGlobalSearch(query);

  // Close when clicking outside

  useEffect(() => {
    function handleClickOutside(
      event: MouseEvent
    ) {
      if (
        wrapperRef.current &&
        !wrapperRef.current.contains(
          event.target as Node
        )
      ) {
        setOpen(false);
      }
    }

    document.addEventListener(
      "mousedown",
      handleClickOutside
    );

    return () =>
      document.removeEventListener(
        "mousedown",
        handleClickOutside
      );
  }, []);

  // ESC closes search

  useEffect(() => {
    function handleKeyDown(
      e: KeyboardEvent
    ) {
      if (e.key === "Escape") {
        setOpen(false);
      }
    }

    window.addEventListener(
      "keydown",
      handleKeyDown
    );

    return () =>
      window.removeEventListener(
        "keydown",
        handleKeyDown
      );
  }, []);

  useEffect(() => {
    setOpen(query.length >= 2);
  }, [query]);

  return (
    <div
      ref={wrapperRef}
      className="
      relative
      hidden
      w-full
      max-w-xl
      lg:block
      "
    >
      <Search
        size={18}
        className="
        absolute
        left-4
        top-1/2
        -translate-y-1/2
        text-slate-400
        "
      />

      <input
        value={query}
        onChange={(e) =>
          setQuery(e.target.value)
        }
        placeholder="Search leads, contacts or tasks..."
        className="
        w-full
        rounded-2xl
        border
        border-slate-200
        bg-slate-50
        py-3
        pl-11
        pr-16
        text-sm
        outline-none
        transition

        focus:border-emerald-500
        focus:bg-white
        focus:ring-4
        focus:ring-emerald-100
        "
      />

      {/* Shortcut */}

      <div
        className="
        absolute
        right-4
        top-1/2
        flex
        -translate-y-1/2
        items-center
        gap-1
        rounded-lg
        border
        border-slate-200
        bg-white
        px-2
        py-1
        text-xs
        text-slate-400
        "
      >
        <Command size={12} />

        <span>K</span>
      </div>

      <SearchDropdown
        loading={loading}
        results={results}
        open={open}
        onClose={() => {
          setOpen(false);
          setQuery("");
        }}
      />
    </div>
  );
}