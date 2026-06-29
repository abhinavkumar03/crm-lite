import { ReactNode } from "react";

type Props = {
  search: ReactNode;
  actions?: ReactNode;
};

export default function Toolbar({
  search,
  actions,
}: Props) {
  return (
    <section className="rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-center">
        <div className="flex-1">
          {search}
        </div>

        {actions && (
          <div className="flex gap-3">
            {actions}
          </div>
        )}
      </div>
    </section>
  );
}