import { ReactNode } from "react";

type Props = {
    badge?: string;
    title: string;
    description: string;
    action?: ReactNode;
};

export default function PageHeader({
    badge,
    title,
    description,
    action,
}: Props) {
    return (
        <section className="flex flex-col gap-6 rounded-3xl border border-slate-200 bg-white p-8 shadow-sm lg:flex-row lg:items-center lg:justify-between">
            <div>
                {badge && (
                    <span className="mb-4 inline-flex rounded-full bg-emerald-50 px-4 py-1 text-sm font-medium text-emerald-700">
                        {badge}
                    </span>
                )}

                <h1 className="text-4xl font-bold tracking-tight text-slate-900">
                    {title}
                </h1>
                <p className="mt-2 max-w-2xl text-slate-500">
                    {description}
                </p>
            </div>

            {action && (
                <div className="flex shrink-0">
                    {action}
                </div>
            )}
        </section>
    );
}