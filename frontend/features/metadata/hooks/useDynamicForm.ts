import {
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";

import {
  FieldSchema,
  FieldValue,
  FormValues,
  ModuleField,
  VisibilityRule,
} from "../types";

import { computeHiddenFields } from "../lib/conditions";
import { validateValues } from "../lib/validation";

function defaultForField(field: ModuleField): FieldValue {
  switch (field.field_type) {
    case "multiselect":
      return [];
    case "boolean":
    case "checkbox":
    case "toggle":
      return field.default_value === "true";
    default:
      return field.default_value ?? "";
  }
}

export function buildInitialValues(
  fields: ModuleField[],
  initial?: FormValues
): FormValues {
  const values: FormValues = {};
  for (const field of fields) {
    values[field.api_name] =
      initial?.[field.api_name] ?? defaultForField(field);
  }
  return values;
}

function fieldsSignature(fields: ModuleField[]): string {
  return fields.map((f) => f.id).join(",");
}

type Options = {
  fields: ModuleField[];
  schema: FieldSchema[];
  initialValues?: FormValues;
  visibilityRules?: VisibilityRule[];
};

// useDynamicForm owns the form's runtime state: values, computed visibility, and
// validation. Keeping this logic out of the presentational components follows
// the single-responsibility principle and makes the renderer reusable.
export function useDynamicForm({
  fields,
  schema,
  initialValues,
  visibilityRules = [],
}: Options) {
  const [values, setValues] = useState<FormValues>(() =>
    buildInitialValues(fields, initialValues)
  );
  const [errors, setErrors] = useState<Record<string, string>>({});

  // Re-seed values when the field set changes (e.g. form layout / new fields load).
  const signature = fieldsSignature(fields);
  useEffect(() => {
    setValues((prev) =>
      buildInitialValues(fields, { ...prev, ...initialValues })
    );
    // eslint-disable-next-line react-hooks/exhaustive-deps -- only when field ids change
  }, [signature]);

  const hidden = useMemo(
    () => computeHiddenFields(visibilityRules, values),
    [visibilityRules, values]
  );

  const setValue = useCallback((name: string, value: FieldValue) => {
    setValues((prev) => ({ ...prev, [name]: value }));
    setErrors((prev) => {
      if (!prev[name]) return prev;
      const next = { ...prev };
      delete next[name];
      return next;
    });
  }, []);

  const validate = useCallback((): boolean => {
    const result = validateValues(schema, values, hidden);
    setErrors(result);
    return Object.keys(result).length === 0;
  }, [schema, values, hidden]);

  // visibleValues strips hidden fields so we never submit conditionally-hidden
  // data.
  const visibleValues = useCallback((): FormValues => {
    const out: FormValues = {};
    for (const [key, value] of Object.entries(values)) {
      if (!hidden.has(key)) out[key] = value;
    }
    return out;
  }, [values, hidden]);

  return {
    values,
    errors,
    hidden,
    setValue,
    setErrors,
    validate,
    visibleValues,
  };
}
