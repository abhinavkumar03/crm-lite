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

  const hasResults =
    results.leads.length > 0 ||
    results.contacts.length > 0 ||
    results.tasks.length > 0;

  let index = 0;

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
          <p className="font-medium text-slate-700">
            No results found
          </p>

          <p className="mt-2 text-sm text-slate-500">
            Try another keyword.
          </p>
        </div>
      )}

      {!loading && hasResults && (
        <div className="max-h-[420px] overflow-y-auto py-2">
          {results.leads.length > 0 && (
            <SearchSection title="Leads">
              {results.leads.map((lead) => {
                const current = index++;

                return (
                  <SearchResult
                    key={lead.id}
                    type="lead"
                    title={lead.name}
                    subtitle={lead.company}
                    status={lead.status}
                    href="/leads"
                    onClick={onClose}
                    query={query}
                    active={activeIndex === current}
                  />
                );
              })}
            </SearchSection>
          )}

          {results.contacts.length > 0 && (
            <SearchSection title="Contacts">
              {results.contacts.map((contact) => {
                const current = index++;

                return (
                  <SearchResult
                    key={contact.id}
                    type="contact"
                    title={contact.name}
                    subtitle={contact.email}
                    href="/contacts"
                    onClick={onClose}
                    query={query}
                    active={activeIndex === current}
                  />
                );
              })}
            </SearchSection>
          )}

          {results.tasks.length > 0 && (
            <SearchSection title="Tasks">
              {results.tasks.map((task) => {
                const current = index++;

                return (
                  <SearchResult
                    key={task.id}
                    type="task"
                    title={task.title}
                    status={task.status}
                    href="/tasks"
                    onClick={onClose}
                    query={query}
                    query={query}
                    active={activeIndex === current}
                  />
                );
              })}
            </SearchSection>
          )}
        </div>
      )}
    </div>
  );
}