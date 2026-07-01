"use client";

import {
  CheckCircle2,
  Circle,
} from "lucide-react";

type Props = {
  password: string;
};

export default function PasswordStrength({
  password,
}: Props) {
  const rules = [
    {
      label: "At least 8 characters",
      valid: password.length >= 8,
    },
    {
      label: "One uppercase letter",
      valid: /[A-Z]/.test(password),
    },
    {
      label: "One lowercase letter",
      valid: /[a-z]/.test(password),
    },
    {
      label: "One number",
      valid: /\d/.test(password),
    },
    {
      label: "One special character",
      valid: /[^A-Za-z0-9]/.test(password),
    },
  ];

  const passed = rules.filter(
    (r) => r.valid
  ).length;

  const strength =
    passed <= 2
      ? "Weak"
      : passed === 3 || passed === 4
      ? "Medium"
      : "Strong";

  const strengthColor =
    strength === "Weak"
      ? "bg-red-500"
      : strength === "Medium"
      ? "bg-amber-500"
      : "bg-emerald-500";

  return (
    <div className="rounded-2xl border border-slate-200 bg-slate-50 p-4">
      <div className="mb-3 flex items-center justify-between">
        <span className="text-sm font-semibold text-slate-700">
          Password Strength
        </span>

        <span
          className={`rounded-full px-3 py-1 text-xs font-semibold text-white ${strengthColor}`}
        >
          {strength}
        </span>
      </div>

      <div className="mb-4 h-2 overflow-hidden rounded-full bg-slate-200">
        <div
          className={`h-full transition-all duration-300 ${strengthColor}`}
          style={{
            width: `${(passed / 5) * 100}%`,
          }}
        />
      </div>

      <div className="space-y-2">
        {rules.map((rule) => (
          <div
            key={rule.label}
            className="flex items-center gap-2 text-sm"
          >
            {rule.valid ? (
              <CheckCircle2
                size={18}
                className="text-emerald-500"
              />
            ) : (
              <Circle
                size={18}
                className="text-slate-400"
              />
            )}

            <span
              className={
                rule.valid
                  ? "text-emerald-700"
                  : "text-slate-500"
              }
            >
              {rule.label}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}