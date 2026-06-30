"use client";

import {
  FileText,
  Paperclip,
  Phone,
  Activity,
  LayoutDashboard,
} from "lucide-react";

const tabs = [
  {
    id: "overview",
    label: "Overview",
    icon: LayoutDashboard,
  },
  {
    id: "notes",
    label: "Notes",
    icon: FileText,
  },
  {
    id: "attachments",
    label: "Attachments",
    icon: Paperclip,
  },
  {
    id: "calls",
    label: "Call Logs",
    icon: Phone,
  },
  {
    id: "activity",
    label: "Activity",
    icon: Activity,
  },
];

type Props = {
  active: string;

  onChange: (
    value: string
  ) => void;
};

export default function LeadDetailsTabs({
  active,
  onChange,
}: Props) {
  return (
    <div
      className="
      overflow-x-auto
      rounded-3xl
      border
      border-slate-200
      bg-white
      shadow-sm
      backdrop-blur-xl
      bg-white/90
      "
    >
      <div className="flex min-w-max">
        {tabs.map((tab) => {
          const Icon = tab.icon;

          const selected =
            active === tab.id;

          return (
            <button
              key={tab.id}
              onClick={() =>
                onChange(tab.id)
              }
              className={`
              flex
              items-center
              gap-2
              px-6
              py-4
              text-sm
              font-medium
              transition

              ${
                selected
                  ? "border-emerald-500 text-emerald-600"
                  : "border-transparent text-slate-500 hover:text-slate-700"
              }
              `}
            >
              <Icon size={18} />

              {tab.label}
            </button>
          );
        })}
      </div>
    </div>
  );
}