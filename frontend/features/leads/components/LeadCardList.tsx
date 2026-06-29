import { Lead } from "../types";

import LeadCard from "./LeadCard";

type Props = {
  leads: Lead[];
  onEdit: (lead: Lead) => void;
  onDelete: (lead: Lead) => void;
};

export default function LeadCardList({
  leads,
  onEdit,
  onDelete,
}: Props) {
  if (!leads.length) {
    return (
      <div className="rounded-3xl border border-dashed border-slate-300 bg-white p-10 text-center">
        <p className="text-slate-500">
          No leads found.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {leads.map((lead) => (
        <LeadCard
          key={lead.id}
          lead={lead}
          onEdit={onEdit}
          onDelete={onDelete}
        />
      ))}
    </div>
  );
}