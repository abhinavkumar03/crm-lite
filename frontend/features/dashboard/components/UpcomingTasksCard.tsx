import {
  ArrowRight,
  CalendarDays,
  CheckCircle2,
  Circle,
  Clock3,
} from "lucide-react";

import Link from "next/link";

import { Task } from "../types";

type Props = {
  tasks: Task[];
};

const priorityColors = {
  pending: "bg-amber-100 text-amber-700",
  "in progress": "bg-blue-100 text-blue-700",
  completed: "bg-emerald-100 text-emerald-700",
};

export default function UpcomingTasksCard({
  tasks,
}: Props) {
  return (
    <section className="rounded-3xl border border-slate-200 bg-white shadow-sm">
      {/* Header */}

      <div className="flex items-center justify-between border-b border-slate-100 p-6">
        <div>
          <p className="text-sm font-medium text-slate-500">
            Productivity
          </p>

          <h2 className="mt-1 text-2xl font-bold text-slate-900">
            Upcoming Tasks
          </h2>
        </div>

        <Link
          href="/tasks"
          className="text-sm font-semibold text-emerald-600 transition hover:text-emerald-700"
        >
          View All
        </Link>
      </div>

      {/* Empty State */}

      {(tasks ?? []).length === 0 && (
        <div className="flex flex-col items-center justify-center py-16">
          <CalendarDays
            size={48}
            className="text-slate-300"
          />

          <p className="mt-4 font-medium text-slate-700">
            No upcoming tasks
          </p>

          <p className="mt-2 text-sm text-slate-500">
            You're all caught up 🎉
          </p>
        </div>
      )}

      {/* Task List */}

      <div className="divide-y divide-slate-100">
        {tasks?.map((task) => {
          const status =
            task.status?.toLowerCase() ?? "pending";

          const badge =
            priorityColors[
              status as keyof typeof priorityColors
            ] ??
            "bg-slate-100 text-slate-700";

          const completed =
            status === "completed";

          return (
            <Link
              href={`/tasks`}
              key={task.id}
              className="group flex items-center justify-between p-5 transition-all duration-300 hover:bg-slate-50"
            >
              <div className="flex items-start gap-4">
                {/* Checkbox */}

                <div className="mt-1">
                  {completed ? (
                    <CheckCircle2
                      size={22}
                      className="text-emerald-500"
                    />
                  ) : (
                    <Circle
                      size={22}
                      className="text-slate-300"
                    />
                  )}
                </div>

                {/* Details */}

                <div>
                  <h3
                    className={`font-semibold ${
                      completed
                        ? "text-slate-400 line-through"
                        : "text-slate-900"
                    }`}
                  >
                    {task.title}
                  </h3>

                  <div className="mt-2 flex flex-wrap items-center gap-3">
                    <span
                      className={`rounded-full px-3 py-1 text-xs font-semibold ${badge}`}
                    >
                      {task.status}
                    </span>

                    <div className="flex items-center gap-1 text-xs text-slate-500">
                      <Clock3 size={14} />

                      Today
                    </div>
                  </div>
                </div>
              </div>

              <ArrowRight
                size={18}
                className="text-slate-400 transition group-hover:translate-x-1"
              />
            </Link>
          );
        })}
      </div>

      {/* Footer */}

      {(tasks ?? []).length > 0 && (
        <div className="border-slate-100 p-5">
          <Link
  href="/tasks"
  className="
    flex
    w-full
    items-center
    justify-center
    gap-2
    rounded-xl
    border
    border-slate-200
    bg-white
    py-3
    font-medium
    text-slate-700
    transition
    hover:bg-slate-100
  "
>
  View Task Board

  <ArrowRight size={18} />
</Link>
        </div>
      )}
    </section>
  );
}