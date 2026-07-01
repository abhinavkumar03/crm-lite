"use client";

import { useState } from "react";

import { CreateTaskPayload } from "../types";
import { Lead } from "@/features/leads/types";
import { Contact } from "@/features/contacts/types";

import FormActions from "@/components/common/form/FormActions";
import FormCard from "@/components/common/form/FormCard";
import FormDateTime from "@/components/common/form/FormDateTime";
import FormInput from "@/components/common/form/FormInput";
import FormSection from "@/components/common/form/FormSection";
import FormSelect from "@/components/common/form/FormSelect";
import FormTextarea from "@/components/common/form/FormTextarea";

type Props = {
  initialValues?: CreateTaskPayload;
  submitText: string;
  leads: Lead[];
  contacts: Contact[];
  onSubmit: (
    values: CreateTaskPayload
  ) => Promise<void>;
  onClose?: () => void;
};

export default function TaskForm({
  initialValues,
  submitText,
  leads,
  contacts,
  onSubmit,
  onClose,
}: Props) {
  const [loading, setLoading] =
    useState(false);

  const [form, setForm] =
    useState<CreateTaskPayload>(
      initialValues ?? {
        title: "",
        description: "",
        status: "PENDING",
        due_date: "",
        lead_id: "",
        contact_id: "",
      }
    );

  function update(
    key: keyof CreateTaskPayload,
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

      await onSubmit({
        ...form,
        due_date: form.due_date
          ? new Date(
              form.due_date
            ).toISOString()
          : undefined,
        lead_id:
          form.lead_id || undefined,
        contact_id:
          form.contact_id || undefined,
      });
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
        {/* Task Details */}

        <FormSection
          title="Task Details"
          description="Define the work that needs to be completed."
        >
          <FormInput
            label="Task Title"
            requiredMark
            placeholder="Follow up with Google"
            value={form.title}
            onChange={(e) =>
              update(
                "title",
                e.target.value
              )
            }
            required
          />

          <FormTextarea
            label="Description"
            rows={5}
            placeholder="Add notes, meeting agenda or task description..."
            value={form.description}
            onChange={(e) =>
              update(
                "description",
                e.target.value
              )
            }
          />
        </FormSection>

        {/* Schedule */}

        <FormSection
          title="Schedule"
          description="Control task progress and due date."
        >
          <div className="grid gap-5 md:grid-cols-2">
            <FormSelect
              label="Status"
              value={form.status}
              onChange={(e) =>
                update(
                  "status",
                  e.target.value
                )
              }
            >
              <option value="PENDING">
                Pending
              </option>

              <option value="IN_PROGRESS">
                In Progress
              </option>

              <option value="COMPLETED">
                Completed
              </option>
            </FormSelect>

            <FormDateTime
              label="Due Date"
              value={form.due_date ?? ""}
              onChange={(e) =>
                update(
                  "due_date",
                  e.target.value
                )
              }
            />
          </div>
        </FormSection>

        {/* Relationships */}

        <FormSection
          title="Relationships"
          description="Optionally link this task to a lead or contact."
        >
          <div className="grid gap-5 md:grid-cols-2">
            <FormSelect
              label="Lead"
              value={form.lead_id ?? ""}
              onChange={(e) =>
                update(
                  "lead_id",
                  e.target.value
                )
              }
            >
              <option value="">
                None
              </option>

              {leads.map((lead) => (
                <option
                  key={lead.id}
                  value={lead.id}
                >
                  {lead.name}
                </option>
              ))}
            </FormSelect>

            <FormSelect
              label="Contact"
              value={form.contact_id ?? ""}
              onChange={(e) =>
                update(
                  "contact_id",
                  e.target.value
                )
              }
            >
              <option value="">
                None
              </option>

              {contacts.map(
                (contact) => (
                  <option
                    key={contact.id}
                    value={contact.id}
                  >
                    {contact.first_name}{" "}
                    {contact.last_name}
                  </option>
                )
              )}
            </FormSelect>
          </div>
        </FormSection>

        <FormActions
          loading={loading}
          submitText={submitText}
          onCancel={onClose}
        />
      </FormCard>
    </form>
  );
}