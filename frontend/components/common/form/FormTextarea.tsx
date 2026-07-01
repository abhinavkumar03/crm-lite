import { TextareaHTMLAttributes } from "react";

type Props = TextareaHTMLAttributes<HTMLTextAreaElement> & {
  label: string;
  requiredMark?: boolean;
  helperText?: string;
};

export default function FormTextarea({
  label,
 requiredMark,
  helperText,
  className = "",
  rows = 4,
  ...props
}: Props) {
  return (
    <div className="space-y-1">
      <label className="block text-sm font-semibold text-slate-700">
        {label}

        {requiredMark && (
          <span className="ml-1 text-red-500">*</span>
        )}
      </label>

      <textarea
        rows={rows}
        {...props}
        className={`
          w-full
          rounded-2xl
          border
          border-slate-300
          bg-white
          px-4
          py-3
          text-sm
          text-slate-900
          shadow-sm
          transition-all
          duration-200
          placeholder:text-slate-400
          focus:border-emerald-500
          focus:ring-4
          focus:ring-emerald-100
          focus:outline-none
          disabled:bg-slate-100
          disabled:text-slate-400
          resize-y
          ${className}
        `}
      />

      {helperText && (
        <p className="text-xs text-slate-500">
          {helperText}
        </p>
      )}
    </div>
  );
}