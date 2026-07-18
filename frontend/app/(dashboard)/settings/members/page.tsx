"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { UserPlus } from "lucide-react";

import Modal from "@/components/common/Modal";
import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";

import { inviteMember, listMembers } from "@/features/organization/api";
import { OrgMember } from "@/features/organization/types";
import { listRoles } from "@/features/roles/api";
import { RoleSummary } from "@/features/roles/types";
import { apiErrorMessage } from "@/features/settings/errors";

export default function MembersSettingsPage() {
  const [members, setMembers] = useState<OrgMember[]>([]);
  const [roles, setRoles] = useState<RoleSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [inviteOpen, setInviteOpen] = useState(false);
  const [email, setEmail] = useState("");
  const [roleId, setRoleId] = useState("");
  const [inviting, setInviting] = useState(false);

  async function reload() {
    const [m, r] = await Promise.all([listMembers(), listRoles()]);
    setMembers(m);
    setRoles(r);
    if (!roleId && r.length) setRoleId(r[0].id);
  }

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        await reload();
      } catch (err) {
        if (active) toast.error(apiErrorMessage(err, "Failed to load members"));
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function handleInvite() {
    if (!email.trim() || !roleId) {
      toast.error("Email and role are required");
      return;
    }
    try {
      setInviting(true);
      const invite = await inviteMember({
        email: email.trim(),
        role_id: roleId,
      });
      toast.success(`Invite created for ${invite.email}`);
      if (invite.token) {
        toast.message(`Accept token (demo): ${invite.token}`);
      }
      setInviteOpen(false);
      setEmail("");
      await reload();
    } catch (err) {
      toast.error(apiErrorMessage(err, "Failed to invite"));
    } finally {
      setInviting(false);
    }
  }

  if (loading) {
    return (
      <div className="rounded-3xl border border-slate-200 bg-white p-8 text-slate-400 shadow-sm">
        Loading members...
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div className="mb-5 flex flex-wrap items-center justify-between gap-3">
          <div>
            <h2 className="text-lg font-semibold text-slate-900">Members</h2>
            <p className="text-sm text-slate-500">
              People in this organization, roles, and reporting links.
            </p>
          </div>
          <button
            type="button"
            onClick={() => setInviteOpen(true)}
            className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-4 py-2 text-sm font-semibold text-white transition hover:bg-emerald-600"
          >
            <UserPlus className="h-4 w-4" />
            Invite
          </button>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full min-w-[640px] text-left text-sm">
            <thead>
              <tr className="border-b border-slate-100 text-xs uppercase tracking-wide text-slate-400">
                <th className="px-2 py-2 font-semibold">Name</th>
                <th className="px-2 py-2 font-semibold">Email</th>
                <th className="px-2 py-2 font-semibold">Role</th>
                <th className="px-2 py-2 font-semibold">Level</th>
                <th className="px-2 py-2 font-semibold">Designation</th>
              </tr>
            </thead>
            <tbody>
              {members.map((m) => (
                <tr key={m.user_id} className="border-b border-slate-50">
                  <td className="px-2 py-3 font-medium text-slate-800">
                    {m.name}
                  </td>
                  <td className="px-2 py-3 text-slate-600">{m.email}</td>
                  <td className="px-2 py-3 text-slate-600">{m.role_slug}</td>
                  <td className="px-2 py-3 text-slate-600">
                    {m.hierarchy_level}
                  </td>
                  <td className="px-2 py-3 text-slate-600">
                    {m.designation || "—"}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <Modal
        open={inviteOpen}
        onClose={() => setInviteOpen(false)}
        title="Invite member"
      >
        <div className="space-y-4">
          <FormInput
            label="Email"
            type="email"
            value={email}
            requiredMark
            onChange={(e) => setEmail(e.target.value)}
          />
          <FormSelect
            label="Role"
            value={roleId}
            onChange={(e) => setRoleId(e.target.value)}
          >
            {roles.map((r) => (
              <option key={r.id} value={r.id}>
                {r.name} (L{r.hierarchy_level})
              </option>
            ))}
          </FormSelect>
          <div className="flex justify-end gap-2">
            <button
              type="button"
              onClick={() => setInviteOpen(false)}
              className="rounded-full border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-600"
            >
              Cancel
            </button>
            <button
              type="button"
              disabled={inviting}
              onClick={handleInvite}
              className="rounded-full bg-emerald-500 px-4 py-2 text-sm font-semibold text-white disabled:opacity-50"
            >
              {inviting ? "Sending…" : "Send invite"}
            </button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
