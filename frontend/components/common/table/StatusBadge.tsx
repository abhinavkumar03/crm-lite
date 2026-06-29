type Props = {
  status: string;
};

const variants: Record<
  string,
  string
> = {
  new: "bg-sky-100 text-sky-700",

  contacted:
    "bg-amber-100 text-amber-700",

  qualified:
    "bg-indigo-100 text-indigo-700",

  won: "bg-emerald-100 text-emerald-700",

  lost: "bg-red-100 text-red-700",

  pending:
    "bg-yellow-100 text-yellow-700",

  "in progress":
    "bg-blue-100 text-blue-700",

  completed:
    "bg-emerald-100 text-emerald-700",

  active:
    "bg-emerald-100 text-emerald-700",

  inactive:
    "bg-slate-200 text-slate-700",
};

export default function StatusBadge({
  status,
}: Props) {
  const value = status.toLowerCase();

  const color =
    variants[value] ??
    "bg-slate-100 text-slate-700";

  return (
    <span
      className={`
      inline-flex
      items-center
      rounded-full
      px-3
      py-1
      text-xs
      font-semibold
      capitalize
      ${color}
      `}
    >
      {status}
    </span>
  );
}