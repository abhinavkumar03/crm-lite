import {
  useCallback,
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
