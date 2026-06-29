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
    <div
      className="
      rounded-3xl
      bg-white
      "
    >
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
  );
}