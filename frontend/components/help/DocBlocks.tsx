import { AlertTriangle, Info, Lightbulb } from "lucide-react";

import type { DocBlock } from "@/features/docs/content";

import EntityGroups from "./EntityGroups";
import FlowSteps from "./FlowSteps";
import SystemMap from "./SystemMap";

type Props = {
  blocks: DocBlock[];
};

export default function DocBlocks({ blocks }: Props) {
  return (
    <div className="space-y-8">
      {blocks.map((block, index) => {
        const key = `${block.type}-${index}`;

        switch (block.type) {
          case "intro":
            return (
              <p key={key} className="text-base leading-7 text-slate-600">
                {block.text}
              </p>
            );

          case "bullets":
            return (
              <div key={key} className="space-y-3">
                {block.title && (
                  <h4 className="text-sm font-semibold uppercase tracking-wider text-emerald-700">
                    {block.title}
                  </h4>
                )}
                <ul className="space-y-2">
                  {block.items.map((item) => (
                    <li
                      key={item}
                      className="flex gap-3 text-sm leading-6 text-slate-600"
                    >
                      <span className="mt-2 h-1.5 w-1.5 shrink-0 rounded-full bg-emerald-500" />
                      <span>{item}</span>
                    </li>
                  ))}
                </ul>
              </div>
            );

          case "table":
            return (
              <div key={key} className="space-y-3">
                {block.title && (
                  <h4 className="text-sm font-semibold uppercase tracking-wider text-emerald-700">
                    {block.title}
                  </h4>
                )}
                <div className="overflow-x-auto rounded-2xl border border-slate-200">
                  <table className="min-w-full text-left text-sm">
                    <thead className="bg-slate-50 text-slate-500">
                      <tr>
                        {block.headers.map((header) => (
                          <th
                            key={header}
                            className="px-4 py-3 font-semibold"
                          >
                            {header}
                          </th>
                        ))}
                      </tr>
                    </thead>
                    <tbody>
                      {block.rows.map((row, rowIndex) => (
                        <tr
                          key={`${key}-row-${rowIndex}`}
                          className="border-t border-slate-100"
                        >
                          {row.map((cell, cellIndex) => (
                            <td
                              key={`${key}-cell-${rowIndex}-${cellIndex}`}
                              className="px-4 py-3 text-slate-700"
                            >
                              <code className="rounded bg-slate-50 px-1.5 py-0.5 text-[13px]">
                                {cell}
                              </code>
                            </td>
                          ))}
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            );

          case "flow":
            return (
              <FlowSteps key={key} title={block.title} steps={block.steps} />
            );

          case "map":
            return (
              <SystemMap
                key={key}
                title={block.title}
                nodes={block.nodes}
                edges={block.edges}
              />
            );

          case "groups":
            return (
              <EntityGroups
                key={key}
                title={block.title}
                groups={block.groups}
              />
            );

          case "callout": {
            const Icon =
              block.tone === "warn"
                ? AlertTriangle
                : block.tone === "tip"
                  ? Lightbulb
                  : Info;
            const styles =
              block.tone === "warn"
                ? "border-amber-200 bg-amber-50 text-amber-900"
                : block.tone === "tip"
                  ? "border-emerald-200 bg-emerald-50 text-emerald-900"
                  : "border-sky-200 bg-sky-50 text-sky-900";

            return (
              <div
                key={key}
                className={`flex gap-3 rounded-2xl border px-4 py-3 text-sm leading-6 ${styles}`}
              >
                <Icon size={18} className="mt-0.5 shrink-0" />
                <p>{block.text}</p>
              </div>
            );
          }

          default:
            return null;
        }
      })}
    </div>
  );
}
