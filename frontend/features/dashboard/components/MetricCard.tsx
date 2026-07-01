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
    <div
      className="
      rounded-2xl
      border
      border-slate-200
      bg-white
      p-5
      shadow-sm
      transition-all
      duration-300
      hover:-translate-y-1
      hover:shadow-lg
      "
    >
      <div className="flex items-center justify-between">
        {/* Left */}

        <div className="flex items-center gap-4">
          <div
            className={`flex h-12 w-12 items-center justify-center rounded-2xl ${color}`}
          >
            <Icon
              size={22}
              className="text-white"
            />
          </div>

          <div>
            <p className="text-sm font-medium text-slate-500">
              {title}
            </p>
          </div>
        </div>

        {/* Right */}

        <h2 className="text-3xl font-bold tracking-tight text-slate-900">
          {value}
        </h2>
      </div>

      {/*
      <div className="mt-4 flex items-center gap-2 text-xs font-medium text-emerald-600">
        <TrendingUp size={14} />
        {trend}
        </div>
      */}
    </div>
  );
}