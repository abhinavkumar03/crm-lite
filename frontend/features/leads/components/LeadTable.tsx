import {
  Building2,
  Mail,
  Phone,
} from "lucide-react";

import { Lead } from "../types";

import DataTable from "@/components/common/table/DataTable";
import AvatarCell from "@/components/common/table/AvatarCell";
import StatusBadge from "@/components/common/table/StatusBadge";
import TableActionMenu from "@/components/common/table/TableActionMenu";
import TablePagination from "@/components/common/table/TablePagination";
import EmptyTable from "@/components/common/table/EmptyTable";

type Props = {
  leads: Lead[];
  page: number;
  setPage: React.Dispatch<React.SetStateAction<number>>;
  onEdit: (lead: Lead) => void;
  onDelete: (lead: Lead) => void;
};

export default function LeadTable({
  leads,
  page,
  setPage,
  onEdit,
  onDelete,
}: Props) {
  return (
    <DataTable
      hasData={leads.length > 0}
      columns={
        <tr className="border-b border-slate-200 text-left">
          <th className="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500">
            Lead
          </th>

          <th className="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500">
            Contact
          </th>

          <th className="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500">
            Status
          </th>

          <th className="px-6 py-4 text-right text-xs font-semibold uppercase tracking-wider text-slate-500">
            Actions
          </th>
        </tr>
      }
      emptyState={
        <EmptyTable
          title="No Leads Found"
          description="Start building your sales pipeline by creating your first lead."
        />
      }
      pagination={
        <TablePagination
          page={page}
          onPageChange={setPage}
        />
      }
    >
      {leads.map((lead) => (
        <tr
          key={lead.id}
          className="transition-colors hover:bg-slate-50"
        >
          {/* Lead */}

          <td className="px-6 py-5">
            <AvatarCell
              name={lead.name}
              subtitle={lead.company}
            />
          </td>

          {/* Contact */}

          <td className="px-6 py-5">
            <div className="space-y-2">
              <div className="flex items-center gap-2 text-sm text-slate-700">
                <Mail
                  size={15}
                  className="text-slate-400"
                />

                <span className="truncate">
                  {lead.email}
                </span>
              </div>

              <div className="flex items-center gap-2 text-sm text-slate-500">
                <Phone
                  size={15}
                  className="text-slate-400"
                />

                {lead.phone}
              </div>

              <div className="flex items-center gap-2 text-sm text-slate-500">
                <Building2
                  size={15}
                  className="text-slate-400"
                />

                {lead.company}
              </div>
            </div>
          </td>

          {/* Status */}

          <td className="px-6 py-5">
            <StatusBadge
              status={lead.status}
            />
          </td>

          {/* Actions */}

          <td className="px-6 py-5 text-right">
            <TableActionMenu
              onEdit={() =>
                onEdit(lead)
              }
              onDelete={() =>
                onDelete(lead)
              }
            />
          </td>
        </tr>
      ))}
    </DataTable>
  );
}