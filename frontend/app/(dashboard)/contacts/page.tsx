"use client";

import { useEffect, useState } from "react";

import {
  createContact,
  deleteContact,
  getContacts,
  updateContact,
} from "@/features/contacts/api";

import { Contact } from "@/features/contacts/types";

import ContactTable from "@/features/contacts/components/ContactTable";
import ContactForm from "@/features/contacts/components/ContactForm";

import Modal from "@/components/common/Modal";

import PageHeader from "@/components/common/PageHeader";
import Toolbar from "@/components/common/Toolbar";
import SearchInput from "@/components/common/SearchInput";
import CreateButton from "@/components/common/CreateButton";

export default function ContactsPage() {
  const [contacts, setContacts] =
    useState<Contact[]>([]);

  const [page, setPage] =
    useState(1);

  const [search, setSearch] =
    useState("");

  const [open, setOpen] =
    useState(false);

  const [editingContact, setEditingContact] =
    useState<Contact | null>(null);

  useEffect(() => {
    loadContacts();
  }, [page, search]);

  async function loadContacts() {
  const res = await getContacts({
    page,
    limit: 10,
    search,
  });

  setContacts(res.data.data);
}

  async function handleDelete(
    contact: Contact
  ) {
    const ok =
      window.confirm(
        "Delete this contact?"
      );

    if (!ok) return;

    await deleteContact(contact.id);

    loadContacts();
  }

  return (
    <div className="space-y-8">
      <PageHeader
        // badge="CRM Module"
        title="Contacts"
        description="Manage people, customer information and business relationships from a single place."
        action={
          <CreateButton
            text="New Contact"
            onClick={() =>
              setOpen(true)
            }
          />
        }
      />

      <Toolbar
        search={
          <SearchInput
            value={search}
            onChange={(value) => {
                setSearch(value);
                setPage(1);
            }}
            placeholder="Search contacts..."
            />
        }
      />

      <ContactTable
        contacts={contacts}
        page={page}
        setPage={setPage}
        onEdit={setEditingContact}
        onDelete={handleDelete}
      />

      <Modal
        open={open}
        title="Create Contact"
        onClose={() =>
          setOpen(false)
        }
      >
        <ContactForm
          submitText="Create Contact"
          onSubmit={async (
            values
          ) => {
            await createContact(
              values
            );

            setOpen(false);

            loadContacts();
          }}
        />
      </Modal>

      <Modal
        open={!!editingContact}
        title="Edit Contact"
        onClose={() =>
          setEditingContact(null)
        }
      >
        {editingContact && (
          <ContactForm
            initialValues={{
              first_name:
                editingContact.first_name,
              last_name:
                editingContact.last_name,
              email:
                editingContact.email,
              phone:
                editingContact.phone,
              company:
                editingContact.company,
            }}
            submitText="Update Contact"
            onSubmit={async (
              values
            ) => {
              await updateContact(
                editingContact.id,
                values
              );

              setEditingContact(
                null
              );

              loadContacts();
            }}
          />
        )}
      </Modal>
    </div>
  );
}