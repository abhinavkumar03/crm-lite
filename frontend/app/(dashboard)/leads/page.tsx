"use client";

import { useEffect, useState } from "react";

import {
  createLead,
  deleteLead,
  getLeads,
  updateLead,
} from "@/features/leads/api";

import { Lead } from "@/features/leads/types";

import LeadTable from "@/features/leads/components/LeadTable";
import LeadForm from "@/features/leads/components/LeadForm";

import Modal from "@/components/common/Modal";

import PageHeader from "@/components/common/PageHeader";
import Toolbar from "@/components/common/Toolbar";
import SearchInput from "@/components/common/SearchInput";
import CreateButton from "@/components/common/CreateButton";
import { toast } from "sonner";

export default function LeadsPage() {
  const [leads, setLeads] = useState<Lead[]>([]);

  const [page, setPage] = useState(1);

  const [search, setSearch] = useState("");

  const [editingLead, setEditingLead] =
    useState<Lead | null>(null);

  const [open, setOpen] =
    useState(false);

  useEffect(() => {
    loadLeads();
  }, [page]);

  async function loadLeads() {
    const res = await getLeads(
      page,
      search
    );

    setLeads(res.data);
  }

  async function handleDelete(
    lead: Lead
  ) {
    const ok = window.confirm(
      "Delete this lead?"
    );

    if (!ok) return;

    await deleteLead(lead.id);
    toast.success("Lead deleted.");

    await loadLeads();
  }

  return (
    <div className="space-y-8">
      <PageHeader
        // badge="CRM Module"
        title="Leads"
        description="Manage your sales pipeline, qualify prospects and track lead progress throughout the customer journey."
        action={
          <CreateButton
            text="New Lead"
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
            onChange={setSearch}
            placeholder="Search leads..."
          />
        }
        actions={
          <button
            onClick={() => {
              setPage(1);
              loadLeads();
            }}
            className="
              rounded-2xl
              bg-emerald-500
              px-6
              py-3
              font-medium
              text-white
              transition
              hover:bg-emerald-600
            "
          >
            Search
          </button>
        }
      />

      <LeadTable
        leads={leads}
        page={page}
        setPage={setPage}
        onEdit={setEditingLead}
        onDelete={handleDelete}
      />

      {/* Create */}

      <Modal
        open={open}
        title="Create Lead"
        onClose={() =>
          setOpen(false)
        }
      >
        <LeadForm
          submitText="Create Lead"
          onSubmit={async (
            values
          ) => {
            await createLead(values);
            toast.success("Lead created.");

            setOpen(false);

            loadLeads();
          }}
          onClose={() => setLeadOpen(false)}
        />
      </Modal>

      {/* Edit */}

      <Modal
        open={!!editingLead}
        title="Edit Lead"
        onClose={() =>
          setEditingLead(null)
        }
      >
        {editingLead && (
          <LeadForm
            initialValues={{
              name: editingLead.name,
              email:
                editingLead.email,
              phone:
                editingLead.phone,
              company:
                editingLead.company,
              status:
                editingLead.status,
              notes:
                editingLead.notes,
            }}
            submitText="Update Lead"
            onSubmit={async (
              values
            ) => {
              await updateLead(editingLead.id,values);
              toast.success("Lead updated.");

              setEditingLead(
                null
              );

              loadLeads();
            }}
            onClose={() => setEditingLead(null)}
          />
        )}
      </Modal>
    </div>
  );
}