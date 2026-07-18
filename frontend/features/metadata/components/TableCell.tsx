import { Check, Minus } from "lucide-react";

import { FieldValue, ModuleField } from "../types";

function optionLabel(field: ModuleField, value: string): string {
  return field.options.find((o) => o.value === value)?.label ?? value;
}

function Badge({ children }: { children: React.ReactNode }) {
  return (
    <span className="inline-flex items-center rounded-full bg-slate-100 px-2.5 py-0.5 text-xs font-medium text-slate-700">
      {children}
    </span>
  );
}

function formatDate(value: string, withTime: boolean): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return withTime ? date.toLocaleString() : date.toLocaleDateString();
}

// TableCell renders a single value according to its field metadata, so each
// column looks right (badges for choices, links for urls, formatted dates, etc.)
// without the table having to know about field types.
export default function TableCell({
  field,
  value,
}: {
  field: ModuleField;
  value: FieldValue;
}) {
  if (value === null || value === undefined || value === "") {
    return <span className="text-slate-300">—</span>;
  }

  switch (field.field_type) {
    case "boolean":
    case "checkbox":
      return value ? (
        <Check className="h-4 w-4 text-emerald-500" />
      ) : (
        <Minus className="h-4 w-4 text-slate-300" />
      );

    case "dropdown":
    case "radio":
    case "user":
    case "lookup":
      return <Badge>{optionLabel(field, String(value))}</Badge>;

    case "multiselect": {
      const values = Array.isArray(value) ? value : [String(value)];
      return (
        <div className="flex flex-wrap gap-1">
          {values.map((v) => (
            <Badge key={String(v)}>{optionLabel(field, String(v))}</Badge>
          ))}
        </div>
      );
    }

    case "currency":
      return <span>{Number(value).toLocaleString(undefined, {
        style: "currency",
        currency: "USD",
      })}</span>;

    case "number":
      return <span className="tabular-nums">{Number(value).toLocaleString()}</span>;

    case "date":
      return <span>{formatDate(String(value), false)}</span>;

    case "datetime":
      return <span>{formatDate(String(value), true)}</span>;

    case "email":
      return (
        <a className="text-emerald-600 hover:underline" href={`mailto:${value}`}>
          {String(value)}
        </a>
      );

    case "url":
      return (
        <a
          className="text-emerald-600 hover:underline"
          href={String(value)}
          target="_blank"
          rel="noreferrer"
        >
          {String(value)}
        </a>
      );

    case "phone":
      return (
        <a className="text-emerald-600 hover:underline" href={`tel:${value}`}>
          {String(value)}
        </a>
      );

    default:
      return <span>{String(value)}</span>;
  }
}
