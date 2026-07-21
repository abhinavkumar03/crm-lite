import {
  SelectHTMLAttributes,
} from "react";

import { ChevronDown } from "lucide-react";

type Option = { value: string; label: string };

type Props =
  SelectHTMLAttributes<HTMLSelectElement> & {
    label: string;
    requiredMark?: boolean;
    helperText?: string;
    options?: Option[];
  };

export default function FormSelect({
  label,
  requiredMark,
  helperText,
  className = "",
  children,
  options,
  ...props
}: Props) {
  return (
    <div className="space-y-1">
      <label className="block text-sm font-semibold text-slate-700">
        {label}

        {requiredMark && (
          <span className="ml-1 text-red-500">
            *
          </span>
        )}
      </label>

      <div className="relative">
        <select
          {...props}
          className={`
            w-full
            appearance-none
            rounded-2xl
            border
            border-slate-300
            bg-white
            px-4
            py-3
            pr-10
            text-sm
            text-slate-900
            shadow-sm
            transition-all
            duration-200
            focus:border-emerald-500
            focus:ring-4
            focus:ring-emerald-100
            focus:outline-none
            disabled:bg-slate-100
            disabled:text-slate-400
            ${className}
          `}
        >
          {options
            ? options.map((o) => (
                <option key={o.value} value={o.value}>
                  {o.label}
                </option>
              ))
            : children}
        </select>

        <ChevronDown
          size={18}
          className="
            pointer-events-none
            absolute
            right-4
            top-1/2
            -translate-y-1/2
            text-slate-400
          "
        />
      </div>

      {helperText && (
        <p className="text-xs text-slate-500">
          {helperText}
        </p>
      )}
    </div>
  );
}
