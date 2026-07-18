"use client";

import {
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";

import { useRouter } from "next/navigation";

import { Search, Command } from "lucide-react";

import useGlobalSearch from "../hooks/useGlobalSearch";

import SearchDropdown from "./SearchDropdown";

export default function SearchBar() {
  const router = useRouter();

  const [query, setQuery] = useState("");

  const [open, setOpen] = useState(false);

  const [activeIndex, setActiveIndex] =
    useState(0);

  const wrapperRef =
    useRef<HTMLDivElement>(null);

  const inputRef =
    useRef<HTMLInputElement>(null);

  const {
    loading,
    results,
  } = useGlobalSearch(query);

  const flatResults = useMemo(
    () =>
      (results.results ?? []).map((hit) => ({
        href: `/tables?module=${hit.module_id}`,
        ...hit,
      })),
    [results]
  );

  // Close on outside click

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

  // Open / close dropdown

  useEffect(() => {
    setOpen(query.trim().length >= 2);
  }, [query]);

  // Reset selection

  useEffect(() => {
    setActiveIndex(0);
  }, [query]);

  // Keyboard shortcuts

  useEffect(() => {
    function handleKeyDown(
      e: KeyboardEvent
    ) {
      // Ctrl+K / Cmd+K

      if (
        (e.ctrlKey || e.metaKey) &&
        e.key.toLowerCase() === "k"
      ) {
        e.preventDefault();

        inputRef.current?.focus();

        if (query.length >= 2) {
          setOpen(true);
        }

        return;
      }

      if (!open) return;

      if (e.key === "Escape") {
        setOpen(false);
        return;
      }

      if (e.key === "ArrowDown") {
        e.preventDefault();

        setActiveIndex((prev) =>
          Math.min(
            flatResults.length - 1,
            prev + 1
          )
        );

        return;
      }

      if (e.key === "ArrowUp") {
        e.preventDefault();

        setActiveIndex((prev) =>
          Math.max(0, prev - 1)
        );

        return;
      }

      if (e.key === "Enter") {
        const selected =
          flatResults[activeIndex];

        if (!selected) return;

        setOpen(false);
        setQuery("");

        router.push(selected.href);
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
  }, [
    open,
    query,
    activeIndex,
    flatResults,
    router,
  ]);

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
        ref={inputRef}
        value={query}
        onChange={(e) =>
          setQuery(e.target.value)
        }
        placeholder="Search records..."
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
        query={query}
        activeIndex={activeIndex}
        onClose={() => {
          setOpen(false);
          setQuery("");
        }}
      />
    </div>
  );
}