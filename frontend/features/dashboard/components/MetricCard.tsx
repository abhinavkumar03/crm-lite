import { LucideIcon, TrendingUp } from "lucide-react";

type Props = {
  title: string;
  value: number | string;
  icon: LucideIcon;
  color: string;
  trend?: string;
};

export default function MetricCard({
  title,
  value,
  icon: Icon,
  color,
  trend = "+12%",
}: Props) {
  return (
    <div className="group rounded-3xl border border-slate-200 bg-white p-6 shadow-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-xl">
      <div className="flex items-start justify-between">
        <div
          className={`flex h-14 w-14 items-center justify-center rounded-2xl ${color}`}
        >
          <Icon className="h-7 w-7 text-white" />
        </div>

        <div className="flex items-center gap-1 rounded-full bg-emerald-50 px-3 py-1 text-xs font-semibold text-emerald-700">
          <TrendingUp size={14} />
          {trend}
        </div>
      </div>

      <div className="mt-8">
        <p className="text-sm font-medium text-slate-500">
          {title}
        </p>

        <h2 className="mt-2 text-4xl font-bold tracking-tight text-slate-900">
          {value}
        </h2>
      </div>

    </div>
  );
}