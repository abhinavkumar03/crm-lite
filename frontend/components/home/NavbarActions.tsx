"use client";

import Link from "next/link";
import { ArrowRight, ExternalLink, LayoutDashboard } from "lucide-react";

import { useAuth } from "@/context/AuthContext";
import UserMenu from "@/components/layout/UserMenu";

type Props = {
  mobile?: boolean;
  onNavigate?: () => void;
};

export default function NavbarActions({ mobile = false, onNavigate }: Props) {
  const auth = useAuth();

  if (auth.loading) {
    return (
      <div className="flex items-center gap-3">
        <div className="h-11 w-28 animate-pulse rounded-xl bg-slate-200" />
      </div>
    );
  }

  if (auth.token) {
    return (
      <div className={mobile ? "space-y-3" : "flex items-center gap-3"}>
        {!mobile && (
          <Link
            href="https://github.com/abhinavkumar03/crm-lite"
            target="_blank"
            rel="noopener noreferrer"
            className="secondary-btn !px-4 !py-2.5 text-sm"
          >
            <ExternalLink size={16} />
            GitHub
          </Link>
        )}

        <Link
          href="/dashboard"
          onClick={onNavigate}
          className={
            mobile
              ? "primary-btn w-full justify-center"
              : "primary-btn !px-4 !py-2.5 text-sm"
          }
        >
          <LayoutDashboard size={16} />
          Open dashboard
          {!mobile && <ArrowRight size={16} />}
        </Link>

        {mobile ? <UserMenu mobile /> : <UserMenu showDashboard />}
      </div>
    );
  }

  return (
    <div className={mobile ? "space-y-3" : "flex items-center gap-3"}>
      <Link
        href="https://github.com/abhinavkumar03/crm-lite"
        target="_blank"
        rel="noopener noreferrer"
        onClick={onNavigate}
        className={
          mobile ? "secondary-btn w-full justify-center" : "secondary-btn !px-4 !py-2.5 text-sm"
        }
      >
        <ExternalLink size={16} />
        GitHub
      </Link>

      <Link
        href="/login"
        onClick={onNavigate}
        className={
          mobile ? "secondary-btn w-full justify-center" : "secondary-btn !px-4 !py-2.5 text-sm"
        }
      >
        Sign in
      </Link>

      <Link
        href="/login"
        onClick={onNavigate}
        className={
          mobile ? "primary-btn w-full justify-center" : "primary-btn !px-4 !py-2.5 text-sm"
        }
      >
        Enter dashboard
        <ArrowRight size={16} />
      </Link>
    </div>
  );
}
