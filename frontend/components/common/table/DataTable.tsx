import { ReactNode } from "react";

type Props = {
  columns: ReactNode;
  children: ReactNode;
  pagination?: ReactNode;
  emptyState?: ReactNode;
  hasData: boolean;
};

export default function DataTable({
  columns,
  children,
  pagination,
  emptyState,
  hasData,
}: Props) {
  return (
    <section className="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-sm">
      <div className="overflow-x-auto">
        <table className="min-w-full">
          <thead className="sticky top-0 bg-slate-50">
            {columns}
          </thead>

          <tbody className="divide-y divide-slate-100 bg-white">
            {hasData
              ? children
              : emptyState}
          </tbody>
        </table>
      </div>

      {pagination}
    </section>
  );
}