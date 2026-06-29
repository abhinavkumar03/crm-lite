import { ReactNode } from "react";

type Props = {
    title?: string;
    description?: string;
    children: ReactNode;
};

export default function FormCard({
    title,
    description,
    children,
}: Props) {
    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
            <div className="flex max-h-[90vh] w-full max-w-4xl flex-col overflow-hidden rounded-3xl bg-white">
                <div className="rounded-3xl bg-white flex-1 overflow-y-auto p-6">
                    {(title || description) && (
                        <div className="mb-8">
                            {title && (
                                <h2 className="text-2xl font-bold text-slate-900">
                                    {title}
                                </h2>
                            )}

                            {description && (
                                <p className="mt-2 text-slate-500">
                                    {description}
                                </p>
                            )}
                        </div>
                    )}

                    <div className="space-y-8">
                        {children}
                    </div>
                </div>
            </div>
        </div>
    );
}