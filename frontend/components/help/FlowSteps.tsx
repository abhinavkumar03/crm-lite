import { ArrowRight } from "lucide-react";

import type { FlowStep } from "@/features/docs/content";

type Props = {
  title: string;
  steps: FlowStep[];
};

export default function FlowSteps({ title, steps }: Props) {
  return (
    <div className="space-y-4">
      <h4 className="text-sm font-semibold uppercase tracking-wider text-emerald-700">
        {title}
      </h4>

      <ol className="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
        {steps.map((step, index) => (
          <li
            key={`${title}-${step.title}`}
            className="relative flex flex-col rounded-2xl border border-slate-200 bg-white p-4 shadow-sm"
          >
            <div className="mb-3 flex items-center gap-2">
              <span className="flex h-7 w-7 items-center justify-center rounded-full bg-emerald-500 text-xs font-bold text-white">
                {index + 1}
              </span>
              {index < steps.length - 1 && (
                <ArrowRight
                  size={14}
                  className="hidden text-slate-300 xl:absolute xl:right-[-0.65rem] xl:top-6 xl:block"
                />
              )}
            </div>
            <p className="font-semibold text-slate-900">{step.title}</p>
            <p className="mt-1 text-sm leading-6 text-slate-500">{step.detail}</p>
          </li>
        ))}
      </ol>
    </div>
  );
}
