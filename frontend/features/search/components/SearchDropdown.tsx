"use client";

import SearchSkeleton from "./SearchSkeleton";
import SearchSection from "./SearchSection";
import SearchResult from "./SearchResult";
import { SearchResponse } from "../types";

type Props = {
  loading: boolean;
  results: SearchResponse;
  open: boolean;
  onClose: () => void;
  activeIndex: number;
  query: string;
};

export default function SearchDropdown({
  loading,
  results,
  open,
  onClose,
  activeIndex,
  query,
}: Props) {
  if (!open) return null;

  const hits = results.results ?? [];
  const hasResults = hits.length > 0;

  return (
    <div
      className="
      absolute
      left-0
      right-0
      top-full
      z-50
      mt-3
      overflow-hidden
      rounded-3xl
      border
      border-slate-200
      bg-white
      shadow-2xl
      "
    >
      {loading && <SearchSkeleton />}

      {!loading && !hasResults && (
        <div className="p-10 text-center">
          <p className="font-medium text-slate-700">No results found</p>
          <p className="mt-2 text-sm text-slate-500">Try another keyword.</p>
        </div>
      )}

      {!loading && hasResults && (
        <div className="max-h-[420px] overflow-y-auto py-2">
          <SearchSection title="Records">
            {hits.map((hit, index) => (
              <SearchResult
                key={hit.id}
                type="record"
                title={hit.title}
                subtitle={hit.subtitle || hit.module_label}
                href={`/tables?module=${hit.module_id}`}
                onClick={onClose}
                query={query}
                active={activeIndex === index}
              />
            ))}
          </SearchSection>
        </div>
      )}
    </div>
  );
}
