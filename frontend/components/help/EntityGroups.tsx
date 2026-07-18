type Group = { name: string; items: string[] };

type Props = {
  title: string;
  groups: Group[];
};

export default function EntityGroups({ title, groups }: Props) {
  return (
    <div className="space-y-4">
      <h4 className="text-sm font-semibold uppercase tracking-wider text-emerald-700">
        {title}
      </h4>

      <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
        {groups.map((group) => (
          <div
            key={group.name}
            className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm"
          >
            <p className="text-sm font-bold text-slate-900">{group.name}</p>
            <ul className="mt-3 flex flex-wrap gap-1.5">
              {group.items.map((item) => (
                <li
                  key={item}
                  className="rounded-md bg-slate-100 px-2 py-1 font-mono text-xs text-slate-700"
                >
                  {item}
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>
    </div>
  );
}
