"use client";

import {
  useEffect,
  useMemo,
  useState,
} from "react";

import FormCard from "@/components/common/form/FormCard";
import FormSection from "@/components/common/form/FormSection";
import FormActions from "@/components/common/form/FormActions";

import DynamicField from "./DynamicField";

import { useDynamicForm } from "../hooks/useDynamicForm";

import {
  FieldType,
  FormValues,
  ModuleField,
  ValidationSchema,
  VisibilityRule,
} from "../types";

type Props = {
  fields: ModuleField[];
  schema: ValidationSchema;
  submitText: string;
  initialValues?: FormValues;
  visibilityRules?: VisibilityRule[];
  sectionTitle?: string;
  sectionDescription?: string;
  externalErrors?: Record<string, string>;
  onSubmit: (values: FormValues) => Promise<void>;
  onChange?: (values: FormValues) => void;
  onCancel?: () => void;
};

const FULL_WIDTH_TYPES: ReadonlySet<FieldType> = new Set<FieldType>([
  "textarea",
  "richtext",
  "json",
  "multiselect",
  "radio",
]);

// DynamicForm renders a complete form from module field metadata. It is fully
// module-agnostic: give it fields + a validation schema and it handles layout,
// conditional visibility, client-side validation, and submission.
export default function DynamicForm({
  fields,
  schema,
  submitText,
  initialValues,
  visibilityRules = [],
  sectionTitle = "Details",
  sectionDescription,
  externalErrors,
  onSubmit,
  onChange,
  onCancel,
}: Props) {
  const [loading, setLoading] = useState(false);

  const { values, errors, hidden, setValue, validate, visibleValues } =
    useDynamicForm({
      fields,
      schema: schema.fields,
      initialValues,
      visibilityRules,
    });

  useEffect(() => {
    onChange?.(values);
  }, [values, onChange]);

  const visibleFields = useMemo(
    () =>
      [...fields]
        .filter((f) => f.is_visible)
        .sort((a, b) => a.sort_order - b.sort_order),
    [fields]
  );

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!validate()) return;

    try {
      setLoading(true);
      await onSubmit(visibleValues());
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={submit} className="space-y-8">
      <FormCard>
        <FormSection title={sectionTitle} description={sectionDescription}>
          <div className="grid gap-5 md:grid-cols-2">
            {visibleFields.map((field) => {
              if (hidden.has(field.api_name)) return null;

              const error =
                errors[field.api_name] ??
                externalErrors?.[field.api_name];

              return (
                <div
                  key={field.id}
                  className={
                    FULL_WIDTH_TYPES.has(field.field_type)
                      ? "md:col-span-2"
                      : ""
                  }
                >
                  <DynamicField
                    field={field}
                    value={values[field.api_name] ?? ""}
                    error={error}
                    onChange={(value) => setValue(field.api_name, value)}
                  />
                </div>
              );
            })}
          </div>
        </FormSection>

        <FormActions
          loading={loading}
          submitText={submitText}
          onCancel={onCancel}
        />
      </FormCard>
    </form>
  );
}
