"use client";

import { useEffect, useState, ReactNode } from "react";
import { usePathname, useRouter } from "next/navigation";

import { listMyOrganizations } from "@/features/organization/api";

type Props = {
  children: ReactNode;
};

/**
 * Ensures authenticated dashboard users have at least one organization.
 * Zero-org users are sent to the workspace setup wizard.
 */
export default function OrgGate({ children }: Props) {
  const router = useRouter();
  const pathname = usePathname();
  const [ready, setReady] = useState(false);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const orgs = await listMyOrganizations();
        if (!active) return;
        if (orgs.length === 0) {
          router.replace("/onboarding/organization");
          return;
        }
        setReady(true);
      } catch {
        if (!active) return;
        // Auth failures are handled by parent layout; treat empty as setup.
        router.replace("/onboarding/organization");
      }
    })();
    return () => {
      active = false;
    };
  }, [router, pathname]);

  if (!ready) {
    return (
      <div className="flex h-screen items-center justify-center text-sm text-slate-500">
        Checking workspace…
      </div>
    );
  }

  return <>{children}</>;
}
