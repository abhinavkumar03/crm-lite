type Node = { id: string; label: string; hint?: string };

type Props = {
  title: string;
  nodes: Node[];
  edges: string[];
};

export default function SystemMap({ title, nodes, edges }: Props) {
  return (
    <div className="space-y-4">
      <h4 className="text-sm font-semibold uppercase tracking-wider text-emerald-700">
        {title}
      </h4>

      <div className="overflow-hidden rounded-2xl border border-slate-200 bg-gradient-to-br from-slate-50 to-emerald-50/40 p-5">
        <div className="flex flex-wrap items-center justify-center gap-3">
          {nodes.map((node, i) => (
            <div key={node.id} className="flex items-center gap-3">
              <div className="min-w-[7.5rem] rounded-xl border border-slate-200 bg-white px-3 py-3 text-center shadow-sm">
                <p className="text-sm font-bold text-slate-900">{node.label}</p>
                {node.hint && (
                  <p className="mt-0.5 text-xs text-slate-500">{node.hint}</p>
                )}
              </div>
              {i < nodes.length - 1 && (
                <span className="hidden text-slate-300 sm:inline" aria-hidden>
                  →
                </span>
              )}
            </div>
          ))}
        </div>

        <ul className="mt-5 grid gap-2 sm:grid-cols-2">
          {edges.map((edge) => (
            <li
              key={edge}
              className="rounded-lg border border-slate-200/80 bg-white/80 px-3 py-2 text-sm text-slate-600"
            >
              {edge}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
