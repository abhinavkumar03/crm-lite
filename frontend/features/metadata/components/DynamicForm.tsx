"use client";

import {
  useEffect,
  useMemo,
  useState,
  type ReactNode,
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

import type { FormLayout, FormLayoutSection } from "@/features/workspace/types";

type Props = {
  fields: ModuleField[];
  schema: ValidationSchema;
  submitText?: string;
  initialValues?: FormValues;
  visibilityRules?: VisibilityRule[];
  sectionTitle?: string;
  sectionDescription?: string;
  /** When set, render sections/order/editable from backend form layout metadata. */
  formLayout?: FormLayout | null;
  externalErrors?: Record<string, string>;
  /** When true, fields stay interactive but there is no create/save submit. */
  previewOnly?: boolean;
  /** Extra footer content (e.g. Open module CTA) shown below actions / preview badge. */
  footerSlot?: ReactNode;
  onSubmit?: (values: FormValues) => Promise<void>;
  onChange?: (values: FormValues) => void;
  onCancel?: () => void;
  emphasizeSubmit?: boolean;
};

const FULL_WIDTH_TYPES: ReadonlySet<FieldType> = new Set<FieldType>([
  "textarea",
  "richtext",
  "json",
  "multiselect",
  "radio",
  "address",
]);

/**
 * Apply form-layout order/flags, and append any ModuleFields missing from the
 * layout (e.g. just-created fields before layout sync) into the first section.
 */
function resolveLayoutSections(
  fields: ModuleField[],
  formLayout?: FormLayout | null
): { sections: FormLayoutSection[]; layoutFields: ModuleField[] } | null {
  if (!formLayout?.sections?.length) return null;

  const byApi = new Map(fields.map((f) => [f.api_name, f]));
  const seen = new Set<string>();
  const layoutFields: ModuleField[] = [];

  const sections: FormLayoutSection[] = formLayout.sections.map((sec) => {
    const resolved = [];
    for (const lf of sec.fields) {
      const base = byApi.get(lf.key);
      if (!base || seen.has(lf.key) || !base.is_visible) continue;
      seen.add(lf.key);
      const merged: ModuleField = {
        ...base,
        label: lf.label || base.label,
        is_required: lf.required,
        is_read_only: !lf.editable,
        is_visible: true,
        sort_order: lf.display_order,
      };
      layoutFields.push(merged);
      resolved.push({
        ...lf,
        label: merged.label,
        required: merged.is_required,
        editable: !merged.is_read_only,
      });
    }
    return { ...sec, fields: resolved };
  });

  const orphans = fields.filter(
    (f) => f.is_visible && !seen.has(f.api_name)
  );
  if (orphans.length > 0) {
    const targetIdx = Math.max(
      0,
      sections.findIndex((s) => s.id !== "system")
    );
    const target = sections[targetIdx] ?? sections[0];
    for (const f of orphans) {
      seen.add(f.api_name);
      layoutFields.push(f);
      target.fields.push({
        id: f.id,
        key: f.api_name,
        label: f.label,
        type: f.field_type,
        required: f.is_required,
        editable: !f.is_read_only,
        locked: false,
        display_order: target.fields.length + 1,
      });
    }
  }

  return { sections, layoutFields };
}

export default function DynamicForm({
  fields,
  schema,
  submitText = "Save",
  initialValues,
  visibilityRules = [],
  sectionTitle = "Details",
  sectionDescription,
  formLayout,
  externalErrors,
  previewOnly = false,
  footerSlot,
  onSubmit,
  onChange,
  onCancel,
  emphasizeSubmit = false,
}: Props) {
  const [loading, setLoading] = useState(false);

  const resolved = useMemo(
    () => resolveLayoutSections(fields, formLayout),
    [fields, formLayout]
  );

  const layoutFields = resolved?.layoutFields ?? fields.filter((f) => f.is_visible);
  const sections = resolved?.sections ?? null;

  const { values, errors, hidden, setValue, validate, visibleValues } =
    useDynamicForm({
      fields: layoutFields,
      schema: schema.fields,
      initialValues,
      visibilityRules,
    });

  useEffect(() => {
    onChange?.(values);
  }, [values, onChange]);

  const flatVisibleFields = useMemo(
    () =>
      [...layoutFields]
        .filter((f) => f.is_visible)
        .sort((a, b) => a.sort_order - b.sort_order),
    [layoutFields]
  );

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (previewOnly || !onSubmit) return;
    if (!validate()) return;

    try {
      setLoading(true);
      await onSubmit(visibleValues());
    } finally {
      setLoading(false);
    }
  }

  const byApi = useMemo(
    () => new Map(layoutFields.map((f) => [f.api_name, f])),
    [layoutFields]
  );

  function renderField(field: ModuleField) {
    if (hidden.has(field.api_name)) return null;
    const error =
      errors[field.api_name] ?? externalErrors?.[field.api_name];
    return (
      <div
        key={field.id}
        className={
          FULL_WIDTH_TYPES.has(field.field_type) ? "md:col-span-2" : ""
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
  }

  return (
    <form onSubmit={submit} className="space-y-8">
      <FormCard>
        {sections ? (
          sections.map((sec) => {
            if (sec.fields.length === 0) return null;
            const cols = sec.columns >= 1 && sec.columns <= 3 ? sec.columns : 2;
            const gridClass =
              cols === 1
                ? "grid gap-5"
                : cols === 3
                  ? "grid gap-5 md:grid-cols-3"
                  : "grid gap-5 md:grid-cols-2";
            return (
              <FormSection
                key={sec.id}
                title={sec.title}
                description={sec.description}
              >
                <div className={gridClass}>
                  {sec.fields.map((lf) => {
                    const field = byApi.get(lf.key);
                    if (!field) return null;
                    return renderField(field);
                  })}
                </div>
              </FormSection>
            );
          })
        ) : (
          <FormSection title={sectionTitle} description={sectionDescription}>
            <div className="grid gap-5 md:grid-cols-2">
              {flatVisibleFields.map((field) => renderField(field))}
            </div>
          </FormSection>
        )}

        {previewOnly ? (
          <div className="flex flex-col gap-3 border-t border-slate-200 pt-6 sm:flex-row sm:items-center sm:justify-between">
            <p className="inline-flex items-center rounded-full border border-amber-200 bg-amber-50 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-amber-800">
              Preview only
            </p>
            {footerSlot}
          </div>
        ) : (
          <>
            <FormActions
              loading={loading}
              submitText={submitText}
              onCancel={onCancel}
              emphasizeSubmit={emphasizeSubmit}
            />
            {footerSlot ? <div className="pt-2">{footerSlot}</div> : null}
          </>
        )}
      </FormCard>
    </form>
  );
}
