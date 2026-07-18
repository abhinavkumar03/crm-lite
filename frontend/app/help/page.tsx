import type { Metadata } from "next";

import Navbar from "@/components/home/Navbar";
import DocsSection from "@/components/help/DocsSection";

export const metadata: Metadata = {
  title: "How it works · CRM Lite",
  description:
    "Learn how CRM Lite works — architecture, data model, import/export, automation, and more.",
};

export default function HelpPage() {
  return (
    <main className="min-h-screen bg-slate-50">
      <Navbar />
      <DocsSection />
    </main>
  );
}
