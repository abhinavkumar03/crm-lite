"use client";

import { useEffect, useRef, useState } from "react";
import {
  BookOpen,
  Boxes,
  ChevronLeft,
  ChevronRight,
  Database,
  FileUp,
  GitBranch,
  Layers3,
  Route,
  Sparkles,
  Terminal,
} from "lucide-react";

import { DOC_GUIDES } from "@/features/docs/content";

import DocBlocks from "./DocBlocks";

const ICONS = {
  architecture: Boxes,
  erd: Database,
  sequences: Route,
  migrations: GitBranch,
  metadata: Layers3,
  "import-export": FileUp,
  automation: Sparkles,
  onboarding: Terminal,
} as const;

export default function DocsSection() {
  const [activeId, setActiveId] = useState(DOC_GUIDES[0]?.id ?? "architecture");
  const articleRef = useRef<HTMLElement>(null);
  const tabRefs = useRef<Record<string, HTMLButtonElement | null>>({});

  const activeIndex = Math.max(
    0,
    DOC_GUIDES.findIndex((guide) => guide.id === activeId)
  );
  const active = DOC_GUIDES[activeIndex] ?? DOC_GUIDES[0];
  const prev = activeIndex > 0 ? DOC_GUIDES[activeIndex - 1] : null;
  const next =
    activeIndex < DOC_GUIDES.length - 1 ? DOC_GUIDES[activeIndex + 1] : null;

  useEffect(() => {
    const hash = window.location.hash.replace("#", "");
    if (hash.startsWith("doc-")) {
      const id = hash.slice(4);
      if (DOC_GUIDES.some((g) => g.id === id)) {
        setActiveId(id);
      }
    }
  }, []);

  useEffect(() => {
    tabRefs.current[activeId]?.scrollIntoView({
      behavior: "smooth",
      inline: "center",
      block: "nearest",
    });
  }, [activeId]);

  function selectGuide(id: string) {
    setActiveId(id);
    window.history.replaceState(null, "", `/help#doc-${id}`);
    articleRef.current?.scrollIntoView({ behavior: "smooth", block: "start" });
  }

  if (!active) {
    return null;
  }

  const ActiveIcon = ICONS[active.id as keyof typeof ICONS] ?? BookOpen;
  const progress = ((activeIndex + 1) / DOC_GUIDES.length) * 100;

  return (
    <section id="help-docs" className="relative bg-slate-50">
      <div className="container-width section-padding !pt-0">
        {/* Mobile / tablet horizontal guide chips */}
        <div className="mt-8 -mx-1 overflow-x-auto pb-2 lg:hidden">
          <div className="flex min-w-max gap-2 px-1">
            {DOC_GUIDES.map((guide, index) => {
              const selected = guide.id === active.id;
              return (
                <button
                  key={guide.id}
                  type="button"
                  ref={(el) => {
                    tabRefs.current[guide.id] = el;
                  }}
                  onClick={() => selectGuide(guide.id)}
                  className={`rounded-full px-4 py-2 text-sm font-semibold transition ${
                    selected
                      ? "bg-emerald-500 text-white shadow-sm"
                      : "border border-slate-200 bg-white text-slate-600 hover:border-emerald-200 hover:text-emerald-700"
                  }`}
                >
                  {index + 1}. {guide.title}
                </button>
              );
            })}
          </div>
        </div>

        <div className="mt-6 grid gap-6 lg:mt-12 lg:grid-cols-[260px_minmax(0,1fr)]">
          <nav
            aria-label="How it works guides"
            className="hidden h-fit rounded-3xl border border-slate-200 bg-white p-3 shadow-sm lg:sticky lg:top-24 lg:block"
          >
            <div className="mb-3 px-3 pt-2">
              <p className="text-xs font-semibold uppercase tracking-wider text-slate-400">
                Topics
              </p>
              <div className="mt-2 h-1.5 overflow-hidden rounded-full bg-slate-100">
                <div
                  className="h-full rounded-full bg-emerald-500 transition-all duration-300"
                  style={{ width: `${progress}%` }}
                />
              </div>
              <p className="mt-2 text-xs text-slate-500">
                {activeIndex + 1} of {DOC_GUIDES.length}
              </p>
            </div>

            <ul className="space-y-1">
              {DOC_GUIDES.map((guide) => {
                const Icon = ICONS[guide.id as keyof typeof ICONS] ?? BookOpen;
                const selected = guide.id === active.id;

                return (
                  <li key={guide.id}>
                    <button
                      type="button"
                      onClick={() => selectGuide(guide.id)}
                      className={`flex w-full items-start gap-3 rounded-2xl px-3 py-3 text-left transition ${
                        selected
                          ? "bg-emerald-50 text-slate-900 ring-1 ring-emerald-200"
                          : "text-slate-600 hover:bg-slate-50 hover:text-slate-900"
                      }`}
                    >
                      <span
                        className={`mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-xl ${
                          selected
                            ? "bg-emerald-500 text-white"
                            : "bg-slate-100 text-slate-600"
                        }`}
                      >
                        <Icon size={18} />
                      </span>
                      <span className="min-w-0">
                        <span className="block text-sm font-semibold">
                          {guide.title}
                        </span>
                        <span className="mt-0.5 block text-xs text-slate-500">
                          {guide.readingTime}
                        </span>
                      </span>
                    </button>
                  </li>
                );
              })}
            </ul>
          </nav>

          <article
            ref={articleRef}
            className="rounded-3xl border border-slate-200 bg-white p-5 shadow-soft md:p-9"
          >
            <header className="border-b border-slate-100 pb-6">
              <div className="flex flex-wrap items-start gap-4">
                <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-emerald-500 text-white shadow-sm">
                  <ActiveIcon size={22} />
                </span>
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <p className="text-xs font-semibold uppercase tracking-wider text-emerald-700">
                      {active.eyebrow}
                    </p>
                    <span className="rounded-full bg-slate-100 px-2.5 py-0.5 text-xs font-semibold text-slate-500">
                      {active.readingTime}
                    </span>
                    <span className="rounded-full bg-slate-100 px-2.5 py-0.5 text-xs font-semibold text-slate-500 lg:hidden">
                      {activeIndex + 1}/{DOC_GUIDES.length}
                    </span>
                  </div>
                  <h2 className="mt-1 text-2xl font-black tracking-tight text-slate-900 md:text-3xl">
                    {active.title}
                  </h2>
                  <p className="mt-3 max-w-3xl text-base leading-7 text-slate-600">
                    {active.summary}
                  </p>
                </div>
              </div>
            </header>

            <div className="mt-8">
              <DocBlocks blocks={active.blocks} />
            </div>

            <footer className="mt-10 flex flex-col gap-3 border-t border-slate-100 pt-6 sm:flex-row sm:items-center sm:justify-between">
              {prev ? (
                <button
                  type="button"
                  onClick={() => selectGuide(prev.id)}
                  className="inline-flex items-center gap-2 rounded-2xl border border-slate-200 bg-white px-4 py-3 text-left text-sm font-semibold text-slate-700 transition hover:border-emerald-200 hover:bg-emerald-50 hover:text-emerald-800"
                >
                  <ChevronLeft size={18} />
                  <span>
                    <span className="block text-xs font-medium text-slate-400">
                      Previous
                    </span>
                    {prev.title}
                  </span>
                </button>
              ) : (
                <span />
              )}

              {next ? (
                <button
                  type="button"
                  onClick={() => selectGuide(next.id)}
                  className="inline-flex items-center justify-end gap-2 rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-right text-sm font-semibold text-emerald-800 transition hover:bg-emerald-100 sm:ml-auto"
                >
                  <span>
                    <span className="block text-xs font-medium text-emerald-600/80">
                      Next
                    </span>
                    {next.title}
                  </span>
                  <ChevronRight size={18} />
                </button>
              ) : (
                <span className="rounded-2xl bg-slate-50 px-4 py-3 text-sm font-medium text-slate-500">
                  You have finished the guides
                </span>
              )}
            </footer>
          </article>
        </div>
      </div>
    </section>
  );
}
