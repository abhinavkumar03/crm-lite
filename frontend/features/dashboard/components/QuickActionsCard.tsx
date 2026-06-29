import {
  Plus,
  UserPlus,
  ClipboardList,
  PhoneCall,
  ArrowRight,
} from "lucide-react";

const actions = [
  {
    title: "New Lead",
    description: "Create a new sales lead",
    icon: Plus,
    color: "bg-emerald-500",
  },
  {
    title: "Add Contact",
    description: "Create a CRM contact",
    icon: UserPlus,
    color: "bg-blue-500",
  },
  {
    title: "Create Task",
    description: "Schedule follow-up work",
    icon: ClipboardList,
    color: "bg-violet-500",
  },
  {
    title: "Log Call",
    description: "Record customer interaction",
    icon: PhoneCall,
    color: "bg-amber-500",
  },
];

export default function QuickActionsCard() {
  return (
    <section className="rounded-3xl border border-slate-200 bg-white shadow-sm">
      <div className="border-b border-slate-100 p-6">
        <p className="text-sm font-medium text-slate-500">
          Productivity
        </p>

        <h2 className="mt-1 text-2xl font-bold">
          Quick Actions
        </h2>
      </div>

      <div className="grid gap-4 p-6">
        {actions.map((action) => {
          const Icon = action.icon;

          return (
            <button
              key={action.title}
              className="group flex items-center justify-between rounded-2xl border border-slate-200 p-4 text-left transition-all duration-300 hover:-translate-y-1 hover:border-emerald-300 hover:bg-emerald-50"
            >
              <div className="flex items-center gap-4">
                <div
                  className={`flex h-12 w-12 items-center justify-center rounded-2xl ${action.color}`}
                >
                  <Icon
                    className="text-white"
                    size={20}
                  />
                </div>

                <div>
                  <h3 className="font-semibold text-slate-900">
                    {action.title}
                  </h3>

                  <p className="text-sm text-slate-500">
                    {action.description}
                  </p>
                </div>
              </div>

              <ArrowRight
                size={18}
                className="text-slate-400 transition group-hover:translate-x-1"
              />
            </button>
          );
        })}
      </div>
    </section>
  );
}