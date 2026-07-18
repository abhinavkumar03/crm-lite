"use client";

import { Suspense, useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import { useAuth } from "@/context/AuthContext";

import Sidebar from "@/components/layout/Sidebar";
import Topbar from "@/components/layout/Topbar";
import OrgGate from "@/components/organization/OrgGate";
import { TourProvider } from "@/features/tour/TourProvider";
import TourOverlay from "@/features/tour/TourOverlay";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const auth = useAuth();

  const router = useRouter();

  const [sidebarOpen, setSidebarOpen] = useState(false);

  useEffect(() => {
    if (!auth.loading && !auth.token) {
      router.replace("/login");
    }
  }, [auth.loading, auth.token, router]);

  if (auth.loading) {
    return (
      <div className="flex h-screen items-center justify-center">
        Loading...
      </div>
    );
  }

  if (!auth.token) return null;

  return (
    <OrgGate>
      <TourProvider>
        <div className="flex h-screen overflow-hidden bg-slate-50">
          <Suspense fallback={null}>
            <Sidebar
              open={sidebarOpen}
              onClose={() => setSidebarOpen(false)}
            />
          </Suspense>

          <div className="flex min-w-0 flex-1 flex-col">
            <Topbar onMenuClick={() => setSidebarOpen(true)} />

            <main className="flex-1 overflow-auto p-6 lg:p-8">{children}</main>
          </div>
        </div>

        <TourOverlay />
      </TourProvider>
    </OrgGate>
  );
}
