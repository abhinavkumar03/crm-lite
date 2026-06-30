"use client";

import { useEffect, useState } from "react";

import { toast } from "sonner";

import Modal from "@/components/common/Modal";

import LeadForm from "@/features/leads/components/LeadForm";
import ContactForm from "@/features/contacts/components/ContactForm";
import TaskForm from "@/features/tasks/components/TaskForm";

import { createLead, getLeads } from "@/features/leads/api";
import { createContact, getContacts } from "@/features/contacts/api";
import { createTask } from "@/features/tasks/api";

import { Lead } from "@/features/leads/types";
import { Contact } from "@/features/contacts/types";

type Props = {
  leadOpen: boolean;
  contactOpen: boolean;
  taskOpen: boolean;

  onCloseLead: () => void;
  onCloseContact: () => void;
  onCloseTask: () => void;

  onSuccess: () => void;
};

export default function DashboardQuickActionModals({
  leadOpen,
  contactOpen,
  taskOpen,
  onCloseLead,
  onCloseContact,
  onCloseTask,
  onSuccess,
}: Props) {
  const [leads, setLeads] =
    useState<Lead[]>([]);

  const [contacts, setContacts] =
    useState<Contact[]>([]);

  useEffect(() => {
    if (taskOpen) {
      loadRelations();
    }
  }, [taskOpen]);

  async function loadRelations() {
  try {
    const [leadRes, contactRes] =
      await Promise.all([
        getLeads({
          page: 1,
          limit: 100,
        }),
        getContacts({
          page: 1,
          limit: 100,
        }),
      ]);

    setLeads(leadRes.data.data);

    setContacts(contactRes.data.data);
  } catch (err) {
    console.error(err);

    toast.error(
      "Unable to load leads and contacts."
    );
  }
}

  return (
    <>
      {/* Lead */}

      <Modal
        open={leadOpen}
        title="Create Lead"
        onClose={onCloseLead}
      >
        <LeadForm
          submitText="Create Lead"
          onSubmit={async (
            values
          ) => {
            await createLead(values);

            toast.success(
              "Lead created successfully."
            );

            onCloseLead();

            onSuccess();
          }}
        />
      </Modal>

      {/* Contact */}

      <Modal
        open={contactOpen}
        title="Create Contact"
        onClose={onCloseContact}
      >
        <ContactForm
          submitText="Create Contact"
          onSubmit={async (
            values
          ) => {
            await createContact(values);

            toast.success(
              "Contact created successfully."
            );

            onCloseContact();

            onSuccess();
          }}
        />
      </Modal>

      {/* Task */}

      <Modal
        open={taskOpen}
        title="Create Task"
        onClose={onCloseTask}
      >
        <TaskForm
          submitText="Create Task"
          leads={leads}
          contacts={contacts}
          onSubmit={async (
            values
          ) => {
            await createTask(values);

            toast.success(
              "Task created successfully."
            );

            onCloseTask();

            onSuccess();
          }}
        />
      </Modal>
    </>
  );
}