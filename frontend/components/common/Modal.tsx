"use client";

import { X } from "lucide-react";

interface ModalProps {
  open: boolean;
  title: string;
  onClose: () => void;
  children: React.ReactNode;
}

export default function Modal({
  open,
  title,
  onClose,
  children,
}: ModalProps) {
  if (!open) return null;

  return (
    <div
      data-demo-modal="true"
      className="fixed inset-0 z-[70] flex items-center justify-center bg-black/40 p-4"
    >
      <div className="flex max-h-[min(90vh,720px)] w-full max-w-4xl flex-col overflow-hidden rounded-3xl bg-white shadow-2xl">
        {/* Header */}

        <div className="flex shrink-0 items-center justify-between border-b border-slate-200 px-6 py-5">
          <h2 className="text-xl font-semibold text-slate-900">
            {title}
          </h2>

          <button
            type="button"
            onClick={onClose}
            className="
              flex
              h-10
              w-10
              items-center
              justify-center
              rounded-xl
              text-slate-500
              transition
              hover:bg-red-50
              hover:text-red-500
            "
          >
            <X size={20} />
          </button>
        </div>

        {/* Body — min-h-0 required for flex child overflow scrolling */}
        <div className="min-h-0 flex-1 overflow-y-auto overscroll-contain p-6">
          {children}
        </div>
      </div>
    </div>
  );
}