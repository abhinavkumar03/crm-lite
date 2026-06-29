import {
  Users,
  BriefcaseBusiness,
  CheckCircle2,
  TrendingUp,
  CalendarDays,
  Database,
} from "lucide-react";

const metrics = [
  {
    title: "Leads",
    value: "128",
    change: "+18%",
    icon: Users,
  },
  {
    title: "Deals",
    value: "$42K",
    change: "+12%",
    icon: BriefcaseBusiness,
  },
  {
    title: "Tasks",
    value: "24",
    change: "Today",
    icon: CheckCircle2,
  },
];

const pipeline = [
  {
    stage: "Prospect",
    value: 82,
    color: "bg-emerald-500",
  },
  {
    stage: "Qualified",
    value: 54,
    color: "bg-teal-500",
  },
  {
    stage: "Proposal",
    value: 31,
    color: "bg-cyan-500",
  },
  {
    stage: "Won",
    value: 18,
    color: "bg-sky-500",
  },
];

const leads = [
  "John Smith",
  "Sarah Wilson",
  "Alex Brown",
  "Emily Clark",
];

const tasks = [
  "Follow Up",
  "Prepare Proposal",
  "Client Demo",
  "Sales Meeting",
];

export default function DashboardPreview() {
  return (
    <div className="relative">

      {/* Floating Card */}

      <div className="glass shadow-soft absolute -left-6 top-16 hidden rounded-2xl p-4 lg:block">
        <div className="flex items-center gap-3">
          <div className="rounded-xl bg-emerald-100 p-2">
            <TrendingUp
              className="text-emerald-600"
              size={18}
            />
          </div>

          <div>
            <p className="text-xs text-slate-500">
              Conversion
            </p>

            <h4 className="font-bold">
              97%
            </h4>
          </div>
        </div>
      </div>

      <div className="glass shadow-soft absolute -right-5 bottom-20 hidden rounded-2xl p-4 lg:block">
        <div className="flex items-center gap-3">
          <div className="rounded-xl bg-cyan-100 p-2">
            <Database
              className="text-cyan-600"
              size={18}
            />
          </div>

          <div>
            <p className="text-xs text-slate-500">
              Redis
            </p>

            <h4 className="font-bold">
              Enabled
            </h4>
          </div>
        </div>
      </div>

      {/* Browser */}

      <div className="card shadow-soft overflow-hidden">

        {/* Browser Header */}

        <div className="flex items-center gap-2 border-b border-slate-200 bg-slate-50 px-5 py-4">
          <div className="h-3 w-3 rounded-full bg-red-400" />
          <div className="h-3 w-3 rounded-full bg-yellow-400" />
          <div className="h-3 w-3 rounded-full bg-green-400" />

          <div className="ml-6 rounded-full bg-white px-4 py-1 text-xs text-slate-400 shadow-sm">
            crm-lite.app/dashboard
          </div>
        </div>

        <div className="space-y-6 p-6">

          {/* Metrics */}

          <div className="grid grid-cols-3 gap-4">

            {metrics?.map((item) => {
              const Icon = item.icon;

              return (
                <div
                  key={item.title}
                  className="rounded-2xl border border-slate-200 bg-white p-4 transition hover:-translate-y-1 hover:shadow-md"
                >
                  <div className="flex items-center justify-between">

                    <div className="rounded-xl bg-emerald-50 p-2">
                      <Icon
                        size={18}
                        className="text-emerald-600"
                      />
                    </div>

                    <span className="text-xs font-semibold text-emerald-600">
                      {item.change}
                    </span>

                  </div>

                  <p className="mt-4 text-xs text-slate-500">
                    {item.title}
                  </p>

                  <h3 className="mt-1 text-2xl font-bold text-slate-900">
                    {item.value}
                  </h3>

                </div>
              );
            })}

          </div>

          {/* Two Column */}

          <div className="grid gap-5 lg:grid-cols-[1.2fr_.8fr]">

            {/* Pipeline */}

            <div className="rounded-2xl border border-slate-200 p-5">

              <div className="mb-5 flex items-center justify-between">

                <h3 className="font-semibold">
                  Sales Pipeline
                </h3>

                <TrendingUp
                  size={18}
                  className="text-emerald-500"
                />

              </div>

              <div className="space-y-4">

                {pipeline?.map((item) => (
                  <div key={item.stage}>

                    <div className="mb-2 flex justify-between text-sm">

                      <span>{item.stage}</span>

                      <span>{item.value}</span>

                    </div>

                    <div className="h-2 rounded-full bg-slate-100">

                      <div
                        className={`${item.color} h-2 rounded-full`}
                        style={{
                          width: `${item.value}%`,
                        }}
                      />

                    </div>

                  </div>
                ))}

              </div>

            </div>

            {/* Recent Leads */}

            <div className="rounded-2xl border border-slate-200 p-5">

              <h3 className="mb-5 font-semibold">
                Recent Leads
              </h3>

              <div className="space-y-4">

                {leads?.map((lead) => (

                  <div
                    key={lead}
                    className="flex items-center gap-3"
                  >

                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-emerald-500 to-teal-500 text-sm font-bold text-white">
                      {lead.charAt(0)}
                    </div>

                    <div>

                      <h4 className="text-sm font-semibold">
                        {lead}
                      </h4>

                      <p className="text-xs text-slate-500">
                        New Customer
                      </p>

                    </div>

                  </div>

                ))}

              </div>

            </div>

          </div>

          {/* Tasks */}

          <div className="rounded-2xl border border-slate-200 p-5">

            <div className="mb-5 flex items-center gap-2">

              <CalendarDays
                className="text-emerald-500"
                size={18}
              />

              <h3 className="font-semibold">
                Upcoming Tasks
              </h3>

            </div>

            <div className="grid gap-3 md:grid-cols-2">

              {tasks?.map((task) => (

                <div
                  key={task}
                  className="flex items-center gap-3 rounded-xl bg-slate-50 p-3"
                >

                  <CheckCircle2
                    className="text-emerald-500"
                    size={18}
                  />

                  <span className="text-sm">
                    {task}
                  </span>

                </div>

              ))}

            </div>

          </div>

        </div>

      </div>

    </div>
  );
}