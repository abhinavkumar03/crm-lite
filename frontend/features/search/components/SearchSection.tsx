import { ReactNode } from "react";

type Props = {
  title: string;
  children: ReactNode;
};

export default function SearchSection({
  title,
  children,
}: Props) {
  return (
    <div>
      <h3
        className="
        px-4
        py-2
        text-xs
        font-semibold
        uppercase
        tracking-wider
        text-slate-400
        "
      >
        {title}
      </h3>

      <div className="space-y-1">
        {children}
      </div>
    </div>
  );
}