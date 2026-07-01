"use client";

import Link from "next/link";

import {
  ArrowRight,
  ExternalLink,
} from "lucide-react";

import { useAuth } from "@/context/AuthContext";

import UserMenu from "@/components/layout/UserMenu";

type Props = {
  mobile?: boolean;
  onNavigate?: () => void;
};

export default function NavbarActions({
  mobile = false,
  onNavigate,
}: Props) {
  const auth = useAuth();

  if (auth.loading) {
    return (
      <div className="flex items-center gap-3">
        <div className="h-11 w-28 animate-pulse rounded-xl bg-slate-200" />
      </div>
    );
  }

  // Logged In

  if (auth.token) {
    return (
      <div
        className={
          mobile
            ? "space-y-3"
            : "flex items-center gap-3"
        }
      >
        {!mobile && (
          <Link
            href="https://github.com/abhinavkumar03/crm-lite"
            target="_blank"
            className="secondary-btn"
          >
            <ExternalLink size={18} />
            GitHub
          </Link>
        )}

        {mobile ? (
          <>
            <Link
              href="/dashboard"
              onClick={onNavigate}
              className="primary-btn w-full justify-center"
            >
              Dashboard
            </Link>

            <UserMenu mobile />
          </>
        ) : (
          <UserMenu showDashboard />
        )}
      </div>
    );
  }

  // Logged Out

  return (
    <div
      className={
        mobile
          ? "space-y-3"
          : "flex items-center gap-3"
      }
    >
      <Link
        href="https://github.com/abhinavkumar03/crm-lite"
        target="_blank"
        onClick={onNavigate}
        className="secondary-btn"
      >
        <ExternalLink size={18} />
        GitHub
      </Link>

      <Link
        href="/login"
        onClick={onNavigate}
        className="secondary-btn"
      >
        Login
      </Link>

      <Link
        href="/register"
        onClick={onNavigate}
        className="primary-btn"
      >
        Get Started

        <ArrowRight size={18} />
      </Link>
    </div>
  );
}