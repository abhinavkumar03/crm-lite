"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import {
  Mail,
  Loader2,
} from "lucide-react";

import { login } from "@/features/auth/api";
import { useAuth } from "@/context/AuthContext";

import AuthCard from "@/components/auth/AuthCard";
import PasswordInput from "@/components/auth/PasswordInput";
import AuthLayout from "@/components/auth/AuthLayout";
import { toast } from "sonner";

export default function LoginPage() {
  const router = useRouter();

  const auth = useAuth();

  const [email, setEmail] =
    useState("");

  const [password, setPassword] =
    useState("");

  const [rememberMe, setRememberMe] =
    useState(true);

  const [loading, setLoading] =
    useState(false);

  async function handleLogin(
    e: React.FormEvent
  ) {
    e.preventDefault();

    try {
      setLoading(true);

      const res = await login({
        email,
        password,
      });

      await auth.login(
        res.data.access_token
      );
      toast.success("Welcome back!");
      router.replace("/dashboard");
    } catch (err) {
      console.error(err);

      toast.error("Invalid email or password.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthLayout
      title="Welcome Back"
      subtitle="Sign in to continue managing your CRM workspace."
    >
      <AuthCard
        footer={
          <p className="text-sm text-slate-600">
            Don't have an account?{" "}
            <Link
              href="/register"
              className="
              font-semibold
              text-emerald-600
              transition
              hover:text-emerald-700
              "
            >
              Create Account
            </Link>
          </p>
        }
      >
        <form
          onSubmit={handleLogin}
          className="space-y-4"
        >
          {/* Email */}

          <div className="space-y-1">
            <label className="block text-sm font-semibold text-slate-700">
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
                placeholder="Enter your email"
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
                transition

                focus:border-emerald-500
                focus:ring-4
                focus:ring-emerald-100
                "
              />
            </div>
          </div>

          {/* Password */}

          <PasswordInput
            value={password}
            onChange={setPassword}
            required
          />

          {/* Remember */}

          <div className="flex items-center justify-between">
            <label className="flex items-center gap-3 text-sm text-slate-600">
              <input
                type="checkbox"
                checked={rememberMe}
                onChange={() =>
                  setRememberMe(
                    !rememberMe
                  )
                }
                className="
                h-4
                w-4
                rounded
                border-slate-300
                text-emerald-600
                "
              />

              Remember me
            </label>

            {/* <button
              type="button"
              className="
              text-sm
              font-medium
              text-emerald-600
              hover:text-emerald-700
              "
            >
              Forgot password?
            </button> */}
          </div>

          {/* Button */}

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
              ? "Signing In..."
              : "Sign In"}
          </button>
        </form>
      </AuthCard>
    </AuthLayout>
  );
}