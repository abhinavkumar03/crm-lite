import { useMemo, useState, useEffect } from "react";

import { filterRows, paginate, sortRows } from "../lib/table";
import {
  ModuleField,
  SavedView,
  TableRow,
  ViewFilter,
  ViewSort,
} from "../types";

export interface TableConfig {
  columns: string[];
  filters: ViewFilter[];
  sort: ViewSort;
}

interface UseDynamicTableArgs {
  rows: TableRow[];
  fields: ModuleField[];
  pageSize?: number;
  /** Org default list columns from GET /layouts/list (excludes _actions). */
  initialColumns?: string[];
}

const EMPTY_SORT: ViewSort = { field: "", order: "" };

// useDynamicTable owns the client-side table state (visible columns, filters,
// sort and pagination) and derives the rows that should be rendered. It is the
// single source of truth consumed by DynamicTable and the saved-views bar.
export function useDynamicTable({
  rows,
  fields,
  pageSize: initialPageSize = 10,
  initialColumns,
}: UseDynamicTableArgs) {
  const [columns, setColumns] = useState<string[]>(() =>
    initialColumns?.length
      ? initialColumns
      : fields.map((f) => f.api_name)
  );
  const [filters, setFilters] = useState<ViewFilter[]>([]);
  const [sort, setSort] = useState<ViewSort>(EMPTY_SORT);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(initialPageSize);

  // Apply org list layout when it arrives after first paint.
  useEffect(() => {
    if (initialColumns?.length) {
      setColumns(initialColumns);
    }
  }, [initialColumns]);

  const processed = useMemo(() => {
    const filtered = filterRows(rows, filters);
    const sorted = sortRows(filtered, sort, fields);
    return paginate(sorted, page, pageSize);
  }, [rows, filters, sort, fields, page, pageSize]);

  // Toggle sort direction for a column: asc -> desc -> off.
  function toggleSort(apiName: string) {
    setPage(1);
    setSort((prev) => {
      if (prev.field !== apiName) return { field: apiName, order: "asc" };
      if (prev.order === "asc") return { field: apiName, order: "desc" };
      return EMPTY_SORT;
    });
  }

  function setFilter(field: string, filter: Partial<ViewFilter>) {
    setPage(1);
    setFilters((prev) => {
      const existing = prev.find((f) => f.field === field);
      if (!existing) {
        return [
          ...prev,
          { field, operator: "contains", value: "", ...filter },
        ];
      }
      return prev.map((f) => (f.field === field ? { ...f, ...filter } : f));
    });
  }

  function clearFilters() {
    setPage(1);
    setFilters([]);
  }

  // Apply a saved view's configuration to the live table state.
  function applyView(view: SavedView) {
    setColumns(view.columns.length ? view.columns : fields.map((f) => f.api_name));
    setFilters(view.filters ?? []);
    setSort(view.sort?.field ? view.sort : EMPTY_SORT);
    setPage(1);
  }

  function toggleColumn(apiName: string) {
    setColumns((prev) =>
      prev.includes(apiName)
        ? prev.filter((c) => c !== apiName)
        : [...prev, apiName]
    );
  }

  const currentConfig: TableConfig = { columns, filters, sort };

  return {
    columns,
    filters,
    sort,
    page,
    pageSize,
    setPage,
    setPageSize,
    toggleSort,
    setFilter,
    clearFilters,
    toggleColumn,
    applyView,
    currentConfig,
    result: processed,
  };
}
