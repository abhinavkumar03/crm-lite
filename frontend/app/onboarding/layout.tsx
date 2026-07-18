"use client";

import { useEffect, useState, ReactNode } from "react";
import { useRouter } from "next/navigation";

import { useAuth } from "@/context/AuthContext";
import { listMyOrganizations } from "@/features/organization/api";

export default function OnboardingLayout({
  children,
}: {
  children: ReactNode;
}) {
  const auth = useAuth();
  const router = useRouter();
  const [ready, setReady] = useState(false);

  useEffect(() => {
    if (auth.loading) return;
    if (!auth.token) {
      router.replace("/login");
      return;
    }

    let active = true;
    (async () => {
      try {
        const orgs = await listMyOrganizations();
        if (!active) return;
        if (orgs.length > 0) {
          router.replace("/dashboard");
          return;
        }
        setReady(true);
      } catch {
        if (active) setReady(true);
      }
    })();

    return () => {
      active = false;
    };
  }, [auth.loading, auth.token, router]);

  if (auth.loading || !ready) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-50 text-sm text-slate-500">
        Loading…
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-emerald-50">
      {children}
    </div>
  );
}
