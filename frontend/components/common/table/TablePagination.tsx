"use client";

import {
  ChevronLeft,
  ChevronRight,
} from "lucide-react";

type Props = {
  page: number;
  onPageChange: (page: number) => void;
  hasNext?: boolean;
  pageSize?: number;
  totalItems?: number;
};

export default function TablePagination({
  page,
  onPageChange,
  hasNext = true,
  pageSize = 10,
  totalItems,
}: Props) {
  const start =
    (page - 1) * pageSize + 1;

  const end =
    totalItems
      ? Math.min(
          page * pageSize,
          totalItems
        )
      : page * pageSize;

  return (
    <div
      className="
      flex
      flex-col
      gap-4
      border-t
      border-slate-200
      bg-slate-50
      px-6
      py-4
      lg:flex-row
      lg:items-center
      lg:justify-between
      "
    >
      <div className="text-sm text-slate-500">
        {totalItems ? (
          <>
            Showing{" "}
            <span className="font-semibold text-slate-900">
              {start}
            </span>{" "}
            –
            <span className="font-semibold text-slate-900">
              {" "}
              {end}
            </span>{" "}
            of{" "}
            <span className="font-semibold text-slate-900">
              {totalItems}
            </span>
          </>
        ) : (
          <>Page {page}</>
        )}
      </div>

      <div className="flex items-center gap-3">
        <button
          disabled={page === 1}
          onClick={() =>
            onPageChange(page - 1)
          }
          className="
          inline-flex
          items-center
          gap-2
          rounded-xl
          border
          border-slate-200
          bg-white
          px-4
          py-2
          text-sm
          font-medium
          transition
          hover:bg-slate-100
          disabled:cursor-not-allowed
          disabled:opacity-50
          "
        >
          <ChevronLeft size={18} />

          Previous
        </button>

        <div
          className="
          flex
          h-10
          w-10
          items-center
          justify-center
          rounded-xl
          bg-emerald-500
          font-semibold
          text-white
          "
        >
          {page}
        </div>

        <button
          disabled={!hasNext}
          onClick={() =>
            onPageChange(page + 1)
          }
          className="
          inline-flex
          items-center
          gap-2
          rounded-xl
          border
          border-slate-200
          bg-white
          px-4
          py-2
          text-sm
          font-medium
          transition
          hover:bg-slate-100
          disabled:cursor-not-allowed
          disabled:opacity-50
          "
        >
          Next

          <ChevronRight size={18} />
        </button>
      </div>
    </div>
  );
}