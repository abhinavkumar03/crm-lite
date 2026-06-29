import { ReactNode } from "react";

type Props = {
  title: string;
  description?: string;
  children: ReactNode;
};

export default function FormSection({
  title,
  description,
  children,
}: Props) {
  return (
    <section className="rounded-3xl border border-slate-200 bg-slate-50/60 p-6">
      <div className="mb-6">
        <h3 className="text-lg font-semibold text-slate-900">
          {title}
        </h3>

        {description && (
          <p className="mt-1 text-sm text-slate-500">
            {description}
          </p>
        )}
      </div>

      <div className="space-y-5">
        {children}
      </div>
    </section>
  );
}