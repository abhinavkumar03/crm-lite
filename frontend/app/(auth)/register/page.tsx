"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import {
  Loader2,
  Mail,
  User,
} from "lucide-react";

import { register } from "@/features/auth/api";

import AuthLayout from "@/components/auth/AuthLayout";
import AuthCard from "@/components/auth/AuthCard";
import PasswordInput from "@/components/auth/PasswordInput";
import { toast } from "sonner";
import PasswordStrength from "@/components/auth/PasswordStrength";

export default function RegisterPage() {
  const router = useRouter();

  const [name, setName] =
    useState("");

  const [email, setEmail] =
    useState("");

  const [password, setPassword] =
    useState("");

  const [
    confirmPassword,
    setConfirmPassword,
  ] = useState("");

  const [loading, setLoading] =
    useState(false);

  const [error, setError] =
    useState("");

  function validatePassword(
    value: string
  ) {
    return (
      value.length >= 8 &&
      /[A-Z]/.test(value) &&
      /[a-z]/.test(value) &&
      /\d/.test(value) &&
      /[^A-Za-z0-9]/.test(value)
    );
  }

  async function handleSubmit(
    e: React.FormEvent
  ) {
    e.preventDefault();

    setError("");

    if (
      !validatePassword(password)
    ) {
      setError(
        "Password must contain at least 8 characters, one uppercase letter, one lowercase letter, one number and one special character."
      );
      return;
    }

    if (
      password !==
      confirmPassword
    ) {
      setError(
        "Passwords do not match."
      );
      return;
    }

    try {
      setLoading(true);

      await register({
        name,
        email,
        password,
      });

      toast.success(
        "Account created. Sign in to set up your workspace."
      );

      router.replace("/login");
    } catch (err: any) {
      console.error(err);

      setError(
        err?.response?.data
          ?.message ??
          "Registration failed."
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthLayout
      // title="Create Account"
      // subtitle="Start managing leads, contacts and tasks with CRM Lite."
      title=""
      subtitle=""
    >
      <AuthCard
        footer={
          <p className="text-sm text-slate-600">
            Already have an account?{" "}
            <Link
              href="/login"
              className="
              font-semibold
              text-emerald-600
              hover:text-emerald-700
              "
            >
              Sign In
            </Link>
          </p>
        }
      >
        <form
          onSubmit={handleSubmit}
          className="space-y-4"
        >
          {/* Name */}

          <div className="space-y-1">
            <label className="text-sm font-semibold text-slate-700">
              Full Name
            </label>

            <div className="relative">
              <User
                size={18}
                className="
                absolute
                left-4
                top-1/2
                -translate-y-1/2
                text-slate-400
                "
              />

              <input
                required
                autoComplete="name"
                value={name}
                onChange={(e) =>
                  setName(
                    e.target.value
                  )
                }
                placeholder="John Doe"
                className="
                w-full
                rounded-2xl
                border
                border-slate-200
                py-3
                pl-12
                pr-4
                text-sm
                outline-none

                focus:border-emerald-500
                focus:ring-4
                focus:ring-emerald-100
                "
              />
            </div>
          </div>

          {/* Email */}

          <div className="space-y-1">
            <label className="text-sm font-semibold text-slate-700">
              Email Address
            </label>

            <div className="relative">
              <Mail
                size={18}
                className="
                absolute
                left-4
                top-1/2
                -translate-y-1/2
                text-slate-400
                "
              />

              <input
                type="email"
                required
                autoComplete="email"
                value={email}
                onChange={(e) =>
                  setEmail(
                    e.target.value
                  )
                }
                placeholder="john@example.com"
                className="
                w-full
                rounded-2xl
                border
                border-slate-200
                py-3
                pl-12
                pr-4
                text-sm
                outline-none

                focus:border-emerald-500
                focus:ring-4
                focus:ring-emerald-100
                "
              />
            </div>
          </div>

          {/* Password */}

          <PasswordInput
            label="Password"
            value={password}
            onChange={setPassword}
            autoComplete="new-password"
            required
          />

          {/* Confirm Password */}

          <PasswordInput
            label="Confirm Password"
            value={confirmPassword}
            onChange={
              setConfirmPassword
            }
            autoComplete="new-password"
            required
          />

          {confirmPassword.length > 0 && (
            <div
              className={`rounded-xl px-4 py-3 text-sm ${
                password === confirmPassword
                  ? "bg-emerald-50 text-emerald-700"
                  : "bg-red-50 text-red-700"
              }`}
            >
              {password === confirmPassword
                ? "✓ Passwords match"
                : "✗ Passwords do not match"}
            </div>
          )}

          {/* Password Helper */}

          <PasswordStrength
            password={password}
          />

          <p className="rounded-2xl border border-slate-100 bg-slate-50 px-4 py-3 text-xs text-slate-500">
            After you sign in, we&apos;ll guide you through creating your
            organization workspace (roles and modules are set up automatically).
          </p>

          {/* Error */}

          {error && (
            <div className="rounded-2xl border border-red-200 bg-red-50 p-4 text-sm text-red-600">
              {error}
            </div>
          )}

          {/* Submit */}

          <button
            disabled={loading}
            className="
            flex
            w-full
            items-center
            justify-center
            gap-2
            rounded-2xl
            bg-emerald-500
            py-3
            font-semibold
            text-white
            transition

            hover:bg-emerald-600

            disabled:cursor-not-allowed
            disabled:opacity-60
            "
          >
            {loading && (
              <Loader2
                size={18}
                className="animate-spin"
              />
            )}

            {loading
              ? "Creating Account..."
              : "Create Account"}
          </button>
        </form>
      </AuthCard>
    </AuthLayout>
  );
}