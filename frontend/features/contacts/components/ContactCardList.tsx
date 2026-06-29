import { Contact } from "../types";

import ContactCard from "./ContactCard";

type Props = {
  contacts: Contact[];
  onEdit: (contact: Contact) => void;
  onDelete: (contact: Contact) => void;
};

export default function ContactCardList({
  contacts,
  onEdit,
  onDelete,
}: Props) {
  if (!contacts.length) {
    return (
      <div className="rounded-3xl border border-dashed border-slate-300 bg-white p-10 text-center">
        <p className="text-slate-500">
          No contacts found.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {contacts.map((contact) => (
        <ContactCard
          key={contact.id}
          contact={contact}
          onEdit={onEdit}
          onDelete={onDelete}
        />
      ))}
    </div>
  );
}