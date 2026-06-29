"use client";

import { useState } from "react";

import { CreateContactPayload } from "../types";

import FormActions from "@/components/common/form/FormActions";
import FormCard from "@/components/common/form/FormCard";
import FormInput from "@/components/common/form/FormInput";
import FormSection from "@/components/common/form/FormSection";

type Props = {
  initialValues?: CreateContactPayload;
  submitText: string;
  onSubmit: (
    values: CreateContactPayload
  ) => Promise<void>;
};

export default function ContactForm({
  initialValues,
  submitText,
  onSubmit,
}: Props) {
  const [loading, setLoading] =
    useState(false);

  const [form, setForm] =
    useState<CreateContactPayload>(
      initialValues ?? {
        first_name: "",
        last_name: "",
        email: "",
        phone: "",
        company: "",
      }
    );

  function update(
    key: keyof CreateContactPayload,
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
        <FormSection
          title="Contact Information"
          description="Store customer information for future communication."
        >
          <div className="grid gap-5 md:grid-cols-2">
            <FormInput
              label="First Name"
              requiredMark
              placeholder="John"
              value={form.first_name}
              onChange={(e) =>
                update(
                  "first_name",
                  e.target.value
                )
              }
              required
            />

            <FormInput
              label="Last Name"
              requiredMark
              placeholder="Smith"
              value={form.last_name}
              onChange={(e) =>
                update(
                  "last_name",
                  e.target.value
                )
              }
              required
            />
          </div>

          <div className="grid gap-5 md:grid-cols-2">
            <FormInput
              label="Email Address"
              type="email"
              requiredMark
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
          </div>

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
        </FormSection>

        <FormActions
          loading={loading}
          submitText={submitText}
        />
      </FormCard>
    </form>
  );
}