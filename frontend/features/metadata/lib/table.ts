import {
  FieldType,
  ModuleField,
  SortDirection,
  TableRow,
  ViewFilter,
  ViewSort,
} from "../types";

const NUMERIC_TYPES: ReadonlySet<FieldType> = new Set<FieldType>([
  "number",
  "currency",
]);

const DATE_TYPES: ReadonlySet<FieldType> = new Set<FieldType>([
  "date",
  "datetime",
]);

function compareValues(
  a: unknown,
  b: unknown,
  type: FieldType
): number {
  const aEmpty = a === null || a === undefined || a === "";
  const bEmpty = b === null || b === undefined || b === "";
  if (aEmpty && bEmpty) return 0;
  if (aEmpty) return -1;
  if (bEmpty) return 1;

  if (NUMERIC_TYPES.has(type)) {
    return Number(a) - Number(b);
  }
  if (DATE_TYPES.has(type)) {
    return new Date(String(a)).getTime() - new Date(String(b)).getTime();
  }
  return String(a).localeCompare(String(b), undefined, {
    sensitivity: "base",
    numeric: true,
  });
}

export function sortRows(
  rows: TableRow[],
  sort: ViewSort,
  fields: ModuleField[]
): TableRow[] {
  if (!sort.field || !sort.order) return rows;

  const field = fields.find((f) => f.api_name === sort.field);
  const type = field?.field_type ?? "text";
  const direction: SortDirection = sort.order === "desc" ? "desc" : "asc";

  return [...rows].sort((a, b) => {
    const result = compareValues(a[sort.field], b[sort.field], type);
    return direction === "asc" ? result : -result;
  });
}

function matchesFilter(row: TableRow, filter: ViewFilter): boolean {
  const raw = row[filter.field];

  switch (filter.operator) {
    case "contains":
      return String(raw ?? "")
        .toLowerCase()
        .includes(String(filter.value ?? "").toLowerCase());
    case "equals":
      return String(raw ?? "") === String(filter.value ?? "");
    case "not_equals":
      return String(raw ?? "") !== String(filter.value ?? "");
    case "gt":
      return Number(raw) > Number(filter.value);
    case "lt":
      return Number(raw) < Number(filter.value);
    case "in":
      return (
        Array.isArray(filter.value) &&
        filter.value.map(String).includes(String(raw ?? ""))
      );
    default:
      return true;
  }
}

export function filterRows(
  rows: TableRow[],
  filters: ViewFilter[]
): TableRow[] {
  const active = filters.filter(
    (f) => f.value !== "" && f.value !== null && f.value !== undefined
  );
  if (active.length === 0) return rows;

  return rows.filter((row) => active.every((f) => matchesFilter(row, f)));
}

export interface Paged<T> {
  rows: T[];
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

export function paginate<T>(
  rows: T[],
  page: number,
  pageSize: number
): Paged<T> {
  const total = rows.length;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const safePage = Math.min(Math.max(1, page), totalPages);
  const start = (safePage - 1) * pageSize;

  return {
    rows: rows.slice(start, start + pageSize),
    page: safePage,
    pageSize,
    total,
    totalPages,
  };
}
