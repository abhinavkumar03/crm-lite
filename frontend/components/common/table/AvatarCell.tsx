type Props = {
  name: string;
  subtitle?: string;
  color?: string;
};

export default function AvatarCell({
  name,
  subtitle,
  color = "from-emerald-500 to-teal-500",
}: Props) {
  const initials = name
    .split(" ")
    .filter(Boolean)
    .map((word) => word[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();

  return (
    <div className="flex items-center gap-4">
      <div
        className={`
          flex
          h-11
          w-11
          shrink-0
          items-center
          justify-center
          rounded-2xl
          bg-gradient-to-br
          ${color}
          text-sm
          font-bold
          text-white
          shadow-sm
        `}
      >
        {initials}
      </div>

      <div className="min-w-0">
        <h4 className="truncate font-semibold text-slate-900">
          {name}
        </h4>

        {subtitle && (
          <p className="truncate text-sm text-slate-500">
            {subtitle}
          </p>
        )}
      </div>
    </div>
  );
}