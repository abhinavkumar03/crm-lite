"use client";

import { Eye, EyeOff, Lock } from "lucide-react";
import { useState } from "react";

type Props = {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  label?: string;
  required?: boolean;
  autoComplete?: string;
};

export default function PasswordInput({
  value,
  onChange,
  placeholder = "Enter your password",
  label = "Password",
  required = false,
  autoComplete = "current-password",
}: Props) {
  const [visible, setVisible] =
    useState(false);

  return (
    <div className="space-y-2">
      <label className="block text-sm font-semibold text-slate-700">
        {label}
      </label>

      <div className="relative">
        {/* Lock */}

        <Lock
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

        {/* Input */}

        <input
          type={
            visible
              ? "text"
              : "password"
          }
          value={value}
          onChange={(e) =>
            onChange(e.target.value)
          }
          placeholder={placeholder}
          autoComplete={autoComplete}
          required={required}
          className="
          w-full
          rounded-2xl
          border
          border-slate-200
          bg-white
          py-3.5
          pl-12
          pr-12
          text-sm
          outline-none
          transition-all
          duration-200

          placeholder:text-slate-400

          focus:border-emerald-500
          focus:ring-4
          focus:ring-emerald-100
          "
        />

        {/* Toggle */}

        <button
          type="button"
          onClick={() =>
            setVisible(!visible)
          }
          className="
          absolute
          right-4
          top-1/2
          -translate-y-1/2
          text-slate-400
          transition
          hover:text-slate-700
          "
        >
          {visible ? (
            <EyeOff size={18} />
          ) : (
            <Eye size={18} />
          )}
        </button>
      </div>
    </div>
  );
}