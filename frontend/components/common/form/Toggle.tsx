"use client";

type Props = {
  checked: boolean;
  onChange: (checked: boolean) => void;
  label?: string;
  description?: string;
  disabled?: boolean;
};

// Toggle is an accessible switch used across the Settings Center for boolean
// preferences and field flags.
export default function Toggle({
  checked,
  onChange,
  label,
  description,
  disabled,
}: Props) {
  const button = (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      disabled={disabled}
      onClick={() => onChange(!checked)}
      className={`relative inline-flex h-6 w-11 shrink-0 items-center rounded-full transition ${
        checked ? "bg-emerald-500" : "bg-slate-300"
      } ${disabled ? "cursor-not-allowed opacity-50" : "cursor-pointer"}`}
    >
      <span
        className={`inline-block h-5 w-5 transform rounded-full bg-white shadow transition ${
          checked ? "translate-x-5" : "translate-x-0.5"
        }`}
      />
    </button>
  );

  if (!label && !description) return button;

  return (
    <div className="flex items-center justify-between gap-4">
      <div className="min-w-0">
        {label && (
          <p className="text-sm font-semibold text-slate-800">{label}</p>
        )}
        {description && (
          <p className="text-xs text-slate-500">{description}</p>
        )}
      </div>
      {button}
    </div>
  );
}
