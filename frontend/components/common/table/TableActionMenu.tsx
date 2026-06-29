"use client";

import { useEffect, useRef, useState } from "react";
import {
  MoreHorizontal,
  Pencil,
  Trash2,
  Eye,
} from "lucide-react";

type Props = {
  onView?: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
};

export default function TableActionMenu({
  onView,
  onEdit,
  onDelete,
}: Props) {
  const [open, setOpen] =
    useState(false);

  const ref =
    useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(
      event: MouseEvent
    ) {
      if (
        ref.current &&
        !ref.current.contains(
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

  return (
    <div
      ref={ref}
      className="relative"
    >
      <button
        onClick={() =>
          setOpen((prev) => !prev)
        }
        className="
          rounded-xl
          p-2
          transition
          hover:bg-slate-100
        "
      >
        <MoreHorizontal
          size={20}
        />
      </button>

      {open && (
        <div
          className="
            absolute
            right-0
            top-12
            z-50
            w-44
            overflow-hidden
            rounded-2xl
            border
            border-slate-200
            bg-white
            shadow-xl
          "
        >
          {onView && (
            <button
              onClick={() => {
                onView();
                setOpen(false);
              }}
              className="
                flex
                w-full
                items-center
                gap-3
                px-4
                py-3
                text-left
                transition
                hover:bg-slate-50
              "
            >
              <Eye size={18} />

              View
            </button>
          )}

          {onEdit && (
            <button
              onClick={() => {
                onEdit();
                setOpen(false);
              }}
              className="
                flex
                w-full
                items-center
                gap-3
                px-4
                py-3
                text-left
                transition
                hover:bg-slate-50
              "
            >
              <Pencil
                size={18}
              />

              Edit
            </button>
          )}

          {onDelete && (
            <button
              onClick={() => {
                onDelete();
                setOpen(false);
              }}
              className="
                flex
                w-full
                items-center
                gap-3
                px-4
                py-3
                text-left
                text-red-600
                transition
                hover:bg-red-50
              "
            >
              <Trash2
                size={18}
              />

              Delete
            </button>
          )}
        </div>
      )}
    </div>
  );
}