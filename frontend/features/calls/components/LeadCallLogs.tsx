"use client";

import {
  useEffect,
  useState,
} from "react";

import {
  Phone,
  Plus,
} from "lucide-react";

import { CallLog } from "../types";
import Modal from "@/components/common/Modal";
import CallForm from "./CallForm";

import { getLeadCalls } from "../api";

import CallLogCard from "./CallLogCard";

import {
  createLeadCall,
    updateLeadCall,
  deleteLeadCall,
} from "../api";

type Props = {
  leadId: string;
};

export default function LeadCallLogs({
  leadId,
}: Props) {
  const [calls, setCalls] =
    useState<CallLog[]>([]);

  const [loading, setLoading] =
    useState(true);
const [open, setOpen] =
  useState(false);
  const [editingCall, setEditingCall] =
  useState<CallLog | null>(
    null
  );
  async function loadCalls() {
    try {
      const data =
        await getLeadCalls(
          leadId
        );

      setCalls(data);
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete(
    call: CallLog
) {
    const ok =
        window.confirm(
    `Delete this ${call.direction.toLowerCase()} call log?`
        );

    if (!ok) {
        return;
    }

    await deleteLeadCall(
        call.id
    );

    await loadCalls();
}

  useEffect(() => {
    loadCalls();
  }, []);

  if (loading) {
    return (
      <div className="rounded-3xl border border-slate-200 bg-white p-10 text-center">
        Loading call logs...
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}

      <div
        className="
        flex
        flex-col
        gap-4
        rounded-3xl
        border
        border-slate-200
        bg-white
        p-6
        shadow-sm
        md:flex-row
        md:items-center
        md:justify-between
        "
      >
        <div>
          <h2 className="text-xl font-semibold">
            Call Logs
          </h2>

          <p className="mt-1 text-sm text-slate-500">
            Record customer conversations and follow-up calls.
          </p>
        </div>

        <button
         onClick={() =>
            setOpen(true)
        }
          className="
          inline-flex
          items-center
          gap-2
          rounded-2xl
          bg-emerald-500
          px-5
          py-3
          font-medium
          text-white
          "
        >
          <Plus size={18} />

          Log Call
        </button>
      </div>

      {calls.length === 0 && (
        <div
          className="
          rounded-3xl
          border
          border-dashed
          border-slate-300
          bg-white
          p-14
          text-center
          "
        >
          <div
            className="
            mx-auto
            flex
            h-20
            w-20
            items-center
            justify-center
            rounded-full
            bg-emerald-50
            "
          >
            <Phone
              size={36}
              className="text-emerald-500"
            />
          </div>

          <h3 className="mt-6 text-xl font-semibold">
            No Calls Logged
          </h3>

          <p className="mx-auto mt-3 max-w-md text-slate-500">
            Record every customer conversation, discovery call,
            demo, pricing discussion and follow-up to keep the
            entire sales team informed.
        </p>

        <button
    onClick={() =>
        setOpen(true)
    }
    className="
    mt-8
    inline-flex
    items-center
    gap-2
    rounded-2xl
    bg-emerald-500
    px-5
    py-3
    font-medium
    text-white
    transition
    hover:bg-emerald-600
    "
>
    <Plus size={18} />

    Log First Call
</button>
        </div>
      )}

      <div className="space-y-5">
        {calls.map((call) => (
          <CallLogCard
            key={call.id}
            call={call}
            onEdit={setEditingCall}
            onDelete={handleDelete}
        />
        ))}
      </div>
      <Modal
  open={open}
  title="Log Call"
  onClose={() =>
    setOpen(false)
  }
>
  <CallForm
    submitText="Save Call"
    onClose={() =>
      setOpen(false)
    }
    onSubmit={async (
      values
    ) => {
      await createLeadCall(
        leadId,
        values
      );

      setOpen(false);

      loadCalls();
    }}
  />
</Modal>
<Modal
    open={!!editingCall}
    title="Edit Call"
    onClose={() =>
        setEditingCall(null)
    }
>
    {editingCall && (
        <CallForm
            initialValues={{
                direction:
                    editingCall.direction,

                status:
                    editingCall.status,

                duration_seconds:
                    editingCall.duration_seconds,

                summary:
                    editingCall.summary,

                follow_up_at:
                    editingCall.follow_up_at
                        ?.slice(0, 16),
            }}
            submitText="Update Call"
            onClose={() =>
                setEditingCall(null)
            }
            onSubmit={async (
                values
            ) => {
                await updateLeadCall(
                    editingCall.id,
                    values
                );

                setEditingCall(
                    null
                );

                loadCalls();
            }}
        />
    )}
</Modal>
    </div>
  );
}