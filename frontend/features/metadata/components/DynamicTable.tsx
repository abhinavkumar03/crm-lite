import {
  ArrowDown,
  ArrowUp,
  ChevronLeft,
  ChevronRight,
  ChevronsUpDown,
  Trash2,
} from "lucide-react";

import TableCell from "./TableCell";
import {
  FilterOperator,
  ModuleField,
  TableRow,
  ViewFilter,
  ViewSort,
} from "../types";

interface Props {
  fields: ModuleField[];
  columns: string[];
  rows: TableRow[];
  sort: ViewSort;
  onToggleSort: (apiName: string) => void;
  filters: ViewFilter[];
  onFilter: (field: string, filter: Partial<ViewFilter>) => void;
  showFilters?: boolean;
  page: number;
  totalPages: number;
  total: number;
  pageSize: number;
  onPage: (page: number) => void;
  onPageSize: (size: number) => void;
  onDeleteRow?: (row: TableRow) => void;
  onRowClick?: (row: TableRow) => void;
}

const PAGE_SIZES = [5, 10, 25, 50];

function SortIcon({ order }: { order: ViewSort["order"] | undefined }) {
  if (order === "asc") return <ArrowUp className="h-3.5 w-3.5" />;
  if (order === "desc") return <ArrowDown className="h-3.5 w-3.5" />;
  return <ChevronsUpDown className="h-3.5 w-3.5 text-slate-300" />;
}

function FilterInput({
  field,
  filter,
  onFilter,
}: {
  field: ModuleField;
  filter?: ViewFilter;
  onFilter: (field: string, filter: Partial<ViewFilter>) => void;
}) {
  const value = filter?.value ?? "";

  const isChoice =
    field.field_type === "dropdown" ||
    field.field_type === "radio" ||
    field.field_type === "user" ||
    field.field_type === "lookup";

  if (isChoice) {
    return (
      <select
        value={String(value)}
        onChange={(e) =>
          onFilter(field.api_name, {
            operator: "equals" as FilterOperator,
            value: e.target.value,
          })
        }
        className="w-full rounded-lg border border-slate-200 bg-white px-2 py-1 text-xs text-slate-700 focus:border-emerald-400 focus:outline-none"
      >
        <option value="">All</option>
        {field.options.map((opt) => (
          <option key={opt.value} value={opt.value}>
            {opt.label}
          </option>
        ))}
      </select>
    );
  }

  return (
    <input
      value={String(value)}
      placeholder="Filter…"
      onChange={(e) =>
        onFilter(field.api_name, {
          operator: "contains" as FilterOperator,
          value: e.target.value,
        })
      }
      className="w-full rounded-lg border border-slate-200 bg-white px-2 py-1 text-xs text-slate-700 focus:border-emerald-400 focus:outline-none"
    />
  );
}

// DynamicTable is a presentational, metadata-driven table. It renders whatever
// columns it is handed (in order), delegates cell rendering to TableCell and
// surfaces sorting/filtering/pagination through callbacks. All state lives in
// the useDynamicTable hook.
export default function DynamicTable({
  fields,
  columns,
  rows,
  sort,
  onToggleSort,
  filters,
  onFilter,
  showFilters = true,
  page,
  totalPages,
  total,
  pageSize,
  onPage,
  onPageSize,
  onDeleteRow,
  onRowClick,
}: Props) {
  const fieldByName = new Map(fields.map((f) => [f.api_name, f]));
  const visible = columns
    .map((name) => fieldByName.get(name))
    .filter((f): f is ModuleField => Boolean(f));

  return (
    <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
      <div className="overflow-x-auto">
        <table className="w-full border-collapse text-sm">
          <thead>
            <tr className="border-b border-slate-200 bg-slate-50">
              {visible.map((field) => {
                const active = sort.field === field.api_name;
                return (
                  <th
                    key={field.api_name}
                    className="px-4 py-3 text-left font-semibold text-slate-600"
                  >
                    <button
                      type="button"
                      onClick={() => onToggleSort(field.api_name)}
                      className="inline-flex items-center gap-1.5 hover:text-slate-900"
                    >
                      {field.label}
                      <SortIcon order={active ? sort.order : undefined} />
                    </button>
                  </th>
                );
              })}
              {onDeleteRow && <th className="w-12 px-4 py-3" />}
            </tr>
            {showFilters && (
              <tr className="border-b border-slate-200 bg-white">
                {visible.map((field) => (
                  <th key={field.api_name} className="px-3 py-2">
                    <FilterInput
                      field={field}
                      filter={filters.find((f) => f.field === field.api_name)}
                      onFilter={onFilter}
                    />
                  </th>
                ))}
                {onDeleteRow && <th className="px-3 py-2" />}
              </tr>
            )}
          </thead>
          <tbody>
            {rows.length === 0 ? (
              <tr>
                <td
                  colSpan={Math.max(1, visible.length + (onDeleteRow ? 1 : 0))}
                  className="px-4 py-10 text-center text-sm text-slate-400"
                >
                  No records to display.
                </td>
              </tr>
            ) : (
              rows.map((row, idx) => (
                <tr
                  key={(row.id as string) ?? idx}
                  onClick={() => onRowClick?.(row)}
                  className={`border-b border-slate-100 last:border-0 hover:bg-slate-50/60 ${
                    onRowClick ? "cursor-pointer" : ""
                  }`}
                >
                  {visible.map((field) => (
                    <td key={field.api_name} className="px-4 py-3 text-slate-700">
                      <TableCell field={field} value={row[field.api_name]} />
                    </td>
                  ))}
                  {onDeleteRow && (
                    <td className="px-4 py-3 text-right">
                      <button
                        type="button"
                        onClick={(e) => {
                          e.stopPropagation();
                          onDeleteRow(row);
                        }}
                        className="text-slate-300 transition hover:text-red-500"
                        aria-label="Delete record"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </td>
                  )}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <div className="flex flex-wrap items-center justify-between gap-3 border-t border-slate-200 px-4 py-3 text-xs text-slate-500">
        <div className="flex items-center gap-2">
          <span>Rows per page</span>
          <select
            value={pageSize}
            onChange={(e) => onPageSize(Number(e.target.value))}
            className="rounded-lg border border-slate-200 bg-white px-2 py-1 text-slate-700 focus:border-emerald-400 focus:outline-none"
          >
            {PAGE_SIZES.map((size) => (
              <option key={size} value={size}>
                {size}
              </option>
            ))}
          </select>
        </div>

        <div className="flex items-center gap-4">
          <span>
            {total === 0
              ? "0 results"
              : `Page ${page} of ${totalPages} · ${total} total`}
          </span>
          <div className="flex items-center gap-1">
            <button
              type="button"
              disabled={page <= 1}
              onClick={() => onPage(page - 1)}
              className="rounded-lg border border-slate-200 p-1.5 text-slate-600 disabled:cursor-not-allowed disabled:opacity-40 hover:bg-slate-100"
            >
              <ChevronLeft className="h-4 w-4" />
            </button>
            <button
              type="button"
              disabled={page >= totalPages}
              onClick={() => onPage(page + 1)}
              className="rounded-lg border border-slate-200 p-1.5 text-slate-600 disabled:cursor-not-allowed disabled:opacity-40 hover:bg-slate-100"
            >
              <ChevronRight className="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
