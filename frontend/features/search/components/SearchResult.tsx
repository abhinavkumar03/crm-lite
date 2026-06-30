"use client";

import Link from "next/link";
import {
  ArrowUpRight,
  Building2,
  ClipboardList,
  User,
} from "lucide-react";

type Props = {
  title: string;
  subtitle?: string;
  href: string;
  type: "lead" | "contact" | "task";
  status?: string;
  onClick?: () => void;
  active?: boolean;
};

export default function SearchResult({
  title,
  subtitle,
  href,
  type,
  status,
  onClick,
  active = false,
}: Props) {
  const Icon =
    type === "lead"
      ? Building2
      : type === "contact"
      ? User
      : ClipboardList;

  return (
    <Link
      href={href}
      onClick={onClick}
      className={`
        flex
        items-center
        justify-between
        rounded-2xl
        px-4
        py-3
        transition-all

        ${
        active
        ? "bg-emerald-50 border border-emerald-200"
        : "hover:bg-slate-100"
        }
        `}
    >
      <div className="flex items-center gap-3">
        <div className="rounded-xl bg-slate-100 p-2">
          <Icon
            size={18}
            className="text-slate-600"
          />
        </div>

        <div>
          <h4 className="font-medium text-slate-900">
            {title}
          </h4>

          {subtitle && (
            <p className="text-sm text-slate-500">
              {subtitle}
            </p>
          )}
        </div>
      </div>

      <div className="flex items-center gap-3">
        {status && (
          <span
            className="
            rounded-full
            bg-slate-100
            px-3
            py-1
            text-xs
            font-medium
            text-slate-600
            "
          >
            {status}
          </span>
        )}

        <ArrowUpRight
          size={16}
          className="text-slate-400"
        />
      </div>
    </Link>
  );
}