"use client";

import { Compass } from "lucide-react";

import { useDemo } from "./DemoProvider";

export default function DemoLauncher() {
  const demo = useDemo();
  if (!demo) return null;

  const running = demo.mode === "running" || demo.session?.status === "active";

  return (
    <button
      type="button"
      data-demo="launcher"
      onClick={() => demo.openLauncher()}
      className="
        fixed bottom-6 right-6 z-[55]
        flex items-center gap-2
        rounded-2xl
        bg-emerald-600 px-4 py-3
        text-sm font-semibold text-white
        shadow-lg shadow-emerald-900/20
        transition hover:bg-emerald-700
        focus:outline-none focus:ring-2 focus:ring-emerald-400 focus:ring-offset-2
      "
    >
      <Compass size={18} />
      {running ? "Demo in progress" : "Explore CRM"}
      {demo.session?.status === "active" && (
        <span className="rounded-lg bg-white/20 px-2 py-0.5 text-xs font-medium">
          {demo.session.progress_percent}%
        </span>
      )}
    </button>
  );
}
