import { InputHTMLAttributes } from "react";
import { CalendarDays } from "lucide-react";

type Props = InputHTMLAttributes<HTMLInputElement> & {
  label: string;
  requiredMark?: boolean;
  helperText?: string;
};

export default function FormDateTime({
  label,
  requiredMark,
  helperText,
  className = "",
  ...props
}: Props) {
  return (
    <div className="space-y-2">
      <label className="block text-sm font-semibold text-slate-700">
        {label}

        {requiredMark && (
          <span className="ml-1 text-red-500">*</span>
        )}
      </label>

      <div className="relative">
        <CalendarDays
          size={18}
          className="
            pointer-events-none
            absolute
            left-4
            top-1/2
            -translate-y-1/2
            text-slate-400
          "
        />

        <input
          type="datetime-local"
          {...props}
          className={`
            w-full
            rounded-2xl
            border
            border-slate-300
            bg-white
            py-3
            pl-12
            pr-4
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