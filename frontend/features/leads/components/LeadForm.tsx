"use client";

import { useState } from "react";

import { X } from "lucide-react";

import { CreateLeadPayload } from "../api";

import FormCard from "@/components/common/form/FormCard";
import FormSection from "@/components/common/form/FormSection";
import FormInput from "@/components/common/form/FormInput";
import FormTextarea from "@/components/common/form/FormTextarea";
import FormSelect from "@/components/common/form/FormSelect";
import FormActions from "@/components/common/form/FormActions";

type Props = {
  initialValues?: CreateLeadPayload;
  submitText: string;
  onSubmit: (
    values: CreateLeadPayload
  ) => Promise<void>;
  onClose?: () => void;
};

export default function LeadForm({
  initialValues,
  submitText,
  onSubmit,
  onClose,
}: Props) {
  const [loading, setLoading] =
    useState(false);

  const [form, setForm] =
    useState<CreateLeadPayload>(
      initialValues ?? {
        name: "",
        email: "",
        phone: "",
        company: "",
        status: "NEW",
        notes: "",
      }
    );

  function update(
    key: keyof CreateLeadPayload,
    value: string
  ) {
    setForm((prev) => ({
      ...prev,
      [key]: value,
    }));
  }

  async function submit(
    e: React.FormEvent
  ) {
    e.preventDefault();

    try {
      setLoading(true);

      await onSubmit(form);
    } finally {
      setLoading(false);
    }
  }

  return (
    <form
      onSubmit={submit}
      className="space-y-8"
    >
      <FormCard>
        {/* Header */}

        <div className="mb-8 flex items-start justify-between">
          <div>
            <h2 className="text-2xl font-bold text-slate-900">
              Create Lead
            </h2>

            <p className="mt-1 text-sm text-slate-500">
              Enter the lead information below.
            </p>
          </div>

          {onClose && (
            <button
              type="button"
              onClick={onClose}
              className="
              flex
              h-10
              w-10
              items-center
              justify-center
              rounded-xl
              border
              border-slate-200
              bg-white
              text-slate-500
              transition
              hover:border-red-200
              hover:bg-red-50
              hover:text-red-500
              "
            >
              <X size={20} />
            </button>
          )}
        </div>

        {/* Basic Information */}

        <FormSection
          title="Basic Information"
          description="Primary information about the lead."
        >
          <div className="grid gap-5 md:grid-cols-2">
            <FormInput
              label="Full Name"
              requiredMark
              placeholder="John Smith"
              value={form.name}
              onChange={(e) =>
                update(
                  "name",
                  e.target.value
                )
              }
              required
            />

            <FormInput
              label="Email Address"
              requiredMark
              type="email"
              placeholder="john@example.com"
              value={form.email}
              onChange={(e) =>
                update(
                  "email",
                  e.target.value
                )
              }
              required
            />

            <FormInput
              label="Phone Number"
              placeholder="+91 9876543210"
              value={form.phone}
              onChange={(e) =>
                update(
                  "phone",
                  e.target.value
                )
              }
            />

            <FormInput
              label="Company"
              placeholder="Google"
              value={form.company}
              onChange={(e) =>
                update(
                  "company",
                  e.target.value
                )
              }
            />
          </div>
        </FormSection>

        {/* Sales Information */}

        <FormSection
          title="Sales Information"
          description="Track the progress of this lead."
        >
          <div className="grid gap-5 md:grid-cols-2">
            <FormSelect
              label="Lead Status"
              value={form.status}
              onChange={(e) =>
                update(
                  "status",
                  e.target.value
                )
              }
            >
              <option value="NEW">
                New
              </option>

              <option value="CONTACTED">
                Contacted
              </option>

              <option value="QUALIFIED">
                Qualified
              </option>

              <option value="WON">
                Won
              </option>

              <option value="LOST">
                Lost
              </option>
            </FormSelect>
          </div>

          <div className="mt-5">
            <FormTextarea
              label="Notes"
              rows={5}
              placeholder="Add meeting notes, follow-up reminders or important customer information..."
              value={form.notes}
              onChange={(e) =>
                update(
                  "notes",
                  e.target.value
                )
              }
            />
          </div>
        </FormSection>

        <FormActions
          loading={loading}
          submitText={submitText}
        />
      </FormCard>
    </form>
  );
}