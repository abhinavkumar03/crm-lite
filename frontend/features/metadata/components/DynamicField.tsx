import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";
import FormTextarea from "@/components/common/form/FormTextarea";
import FormDateTime from "@/components/common/form/FormDateTime";

import {
  FieldValue,
  ModuleField,
} from "../types";

type Props = {
  field: ModuleField;
  value: FieldValue;
  error?: string;
  disabled?: boolean;
  onChange: (value: FieldValue) => void;
};

function FieldError({ error }: { error?: string }) {
  if (!error) return null;
  return <p className="mt-1 text-xs text-red-500">{error}</p>;
}

function FieldLabel({
  field,
}: {
  field: ModuleField;
}) {
  return (
    <label className="block text-sm font-semibold text-slate-700">
      {field.label}
      {field.is_required && <span className="ml-1 text-red-500">*</span>}
    </label>
  );
}

// DynamicField renders a single metadata field using the shared form primitives.
// It is the only place that knows how each field_type maps to an input, keeping
// the form renderer itself type-agnostic.
export default function DynamicField({
  field,
  value,
  error,
  disabled,
  onChange,
}: Props) {
  const common = {
    label: field.label,
    requiredMark: field.is_required,
    helperText: error ? undefined : field.help_text ?? undefined,
    disabled: disabled || field.is_read_only,
    placeholder: field.placeholder ?? undefined,
  };

  switch (field.field_type) {
    case "textarea":
    case "richtext":
    case "json":
    case "address":
      return (
        <div>
          <FormTextarea
            {...common}
            value={String(value ?? "")}
            onChange={(e) => onChange(e.target.value)}
          />
          <FieldError error={error} />
        </div>
      );

    case "number":
    case "currency":
    case "percentage":
      return (
        <div>
          <FormInput
            {...common}
            type="number"
            value={value === null || value === undefined ? "" : String(value)}
            onChange={(e) => onChange(e.target.value)}
          />
          <FieldError error={error} />
        </div>
      );

    case "date":
      return (
        <div>
          <FormInput
            {...common}
            type="date"
            value={String(value ?? "")}
            onChange={(e) => onChange(e.target.value)}
          />
          <FieldError error={error} />
        </div>
      );

    case "time":
      return (
        <div>
          <FormInput
            {...common}
            type="time"
            value={String(value ?? "")}
            onChange={(e) => onChange(e.target.value)}
          />
          <FieldError error={error} />
        </div>
      );

    case "datetime":
      return (
        <div>
          <FormDateTime
            {...common}
            value={String(value ?? "")}
            onChange={(e) => onChange(e.target.value)}
          />
          <FieldError error={error} />
        </div>
      );

    case "email":
    case "phone":
    case "url":
    case "gst":
    case "pan":
    case "barcode":
    case "serial_number": {
      const type =
        field.field_type === "email"
          ? "email"
          : field.field_type === "phone"
          ? "tel"
          : field.field_type === "url"
          ? "url"
          : "text";
      return (
        <div>
          <FormInput
            {...common}
            type={type}
            value={String(value ?? "")}
            onChange={(e) => onChange(e.target.value)}
          />
          <FieldError error={error} />
        </div>
      );
    }

    case "dropdown":
    case "user":
    case "lookup":
      return (
        <div>
          <FormSelect
            {...common}
            value={String(value ?? "")}
            onChange={(e) => onChange(e.target.value)}
          >
            <option value="">Select...</option>
            {field.options.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </FormSelect>
          <FieldError error={error} />
        </div>
      );

    case "radio":
      return (
        <div className="space-y-1">
          <FieldLabel field={field} />
          <div className="flex flex-wrap gap-4 pt-1">
            {field.options.map((opt) => (
              <label
                key={opt.value}
                className="flex items-center gap-2 text-sm text-slate-700"
              >
                <input
                  type="radio"
                  name={field.api_name}
                  value={opt.value}
                  checked={String(value ?? "") === opt.value}
                  disabled={common.disabled}
                  onChange={() => onChange(opt.value)}
                  className="h-4 w-4 accent-emerald-500"
                />
                {opt.label}
              </label>
            ))}
          </div>
          <FieldError error={error} />
        </div>
      );

    case "multiselect": {
      const selected = Array.isArray(value) ? value : [];
      const toggle = (optValue: string) => {
        const next = selected.includes(optValue)
          ? selected.filter((v) => v !== optValue)
          : [...selected, optValue];
        onChange(next);
      };
      return (
        <div className="space-y-1">
          <FieldLabel field={field} />
          <div className="flex flex-wrap gap-3 pt-1">
            {field.options.map((opt) => (
              <label
                key={opt.value}
                className="flex items-center gap-2 rounded-2xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-700"
              >
                <input
                  type="checkbox"
                  checked={selected.includes(opt.value)}
                  disabled={common.disabled}
                  onChange={() => toggle(opt.value)}
                  className="h-4 w-4 accent-emerald-500"
                />
                {opt.label}
              </label>
            ))}
          </div>
          <FieldError error={error} />
        </div>
      );
    }

    case "boolean":
    case "checkbox":
    case "toggle":
      return (
        <div className="space-y-1">
          <label className="flex items-center gap-3 text-sm font-semibold text-slate-700">
            <input
              type="checkbox"
              checked={Boolean(value)}
              disabled={common.disabled}
              onChange={(e) => onChange(e.target.checked)}
              className="h-4 w-4 accent-emerald-500"
            />
            {field.label}
            {field.is_required && <span className="text-red-500">*</span>}
          </label>
          {field.help_text && !error && (
            <p className="text-xs text-slate-500">{field.help_text}</p>
          )}
          <FieldError error={error} />
        </div>
      );

    case "formula":
    case "auto_number":
      return (
        <div>
          <FormInput
            {...common}
            value={String(value ?? "")}
            disabled
            helperText={
              field.field_type === "auto_number"
                ? "Auto-generated"
                : "Calculated field"
            }
          />
          <FieldError error={error} />
        </div>
      );

    case "file":
    case "image":
      return (
        <div>
          <FormInput
            {...common}
            type="file"
            value={undefined}
            onChange={(e) => onChange(e.target.value)}
          />
          <FieldError error={error} />
        </div>
      );

    default:
      return (
        <div>
          <FormInput
            {...common}
            value={String(value ?? "")}
            onChange={(e) => onChange(e.target.value)}
          />
          <FieldError error={error} />
        </div>
      );
  }
}
