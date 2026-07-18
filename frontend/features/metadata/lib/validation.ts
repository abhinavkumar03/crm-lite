import {
  FieldError,
  FieldSchema,
  FormValues,
} from "../types";

import { isEmptyValue } from "./conditions";

const EMAIL_REGEX = /^[^@\s]+@[^@\s]+\.[^@\s]+$/;

function isValidURL(value: string): boolean {
  try {
    const url = new URL(value.trim());
    return !!url.protocol && !!url.host;
  } catch {
    return false;
  }
}

function message(
  field: FieldSchema,
  ruleType: string,
  fallback: string
): string {
  return field.messages?.[ruleType] ?? fallback;
}

// validateValues mirrors the backend engine's field-level checks so the client
// can validate before submitting. Hidden fields are skipped. Returns a map of
// api_name -> first error message.
export function validateValues(
  schema: FieldSchema[],
  values: FormValues,
  hidden: Set<string>
): Record<string, string> {
  const errors: Record<string, string> = {};

  for (const field of schema) {
    if (hidden.has(field.api_name)) continue;

    const value = values[field.api_name];

    if (field.required && isEmptyValue(value)) {
      errors[field.api_name] = message(
        field,
        "required",
        "This field is required"
      );
      continue;
    }

    if (isEmptyValue(value)) continue;

    // Multiselect option membership.
    if (field.multiple) {
      if (Array.isArray(value) && field.options?.length) {
        const invalid = value.some(
          (item) => !field.options!.includes(String(item))
        );
        if (invalid) {
          errors[field.api_name] = message(
            field,
            "in",
            "Invalid option selected"
          );
        }
      }
      continue;
    }

    const text = String(value);
    const length = [...text].length;

    if (field.min_length != null && length < field.min_length) {
      errors[field.api_name] = message(
        field,
        "min_length",
        "Value is too short"
      );
      continue;
    }
    if (field.max_length != null && length > field.max_length) {
      errors[field.api_name] = message(
        field,
        "max_length",
        "Value is too long"
      );
      continue;
    }

    if (field.min != null && Number(value) < field.min) {
      errors[field.api_name] = message(field, "min", "Value is too small");
      continue;
    }
    if (field.max != null && Number(value) > field.max) {
      errors[field.api_name] = message(field, "max", "Value is too large");
      continue;
    }

    if (field.pattern) {
      try {
        if (!new RegExp(field.pattern).test(text)) {
          errors[field.api_name] = message(
            field,
            "pattern",
            "Invalid format"
          );
          continue;
        }
      } catch {
        // Invalid regex on the server side; skip client enforcement.
      }
    }

    if (field.format === "email" && !EMAIL_REGEX.test(text)) {
      errors[field.api_name] = message(
        field,
        "email",
        "Must be a valid email"
      );
      continue;
    }
    if (field.format === "url" && !isValidURL(text)) {
      errors[field.api_name] = message(field, "url", "Must be a valid URL");
      continue;
    }

    if (field.options?.length && !field.options.includes(text)) {
      errors[field.api_name] = message(field, "in", "Invalid option selected");
    }
  }

  return errors;
}

export function errorListToMap(
  errors: FieldError[]
): Record<string, string> {
  const map: Record<string, string> = {};
  for (const e of errors) {
    if (!map[e.field]) map[e.field] = e.message;
  }
  return map;
}
