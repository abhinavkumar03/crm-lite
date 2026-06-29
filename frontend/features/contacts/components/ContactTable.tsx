import {
    Building2,
    Mail,
    Phone,
} from "lucide-react";

import { Contact } from "../types";

import DataTable from "@/components/common/table/DataTable";
import AvatarCell from "@/components/common/table/AvatarCell";
import TableActionMenu from "@/components/common/table/TableActionMenu";
import TablePagination from "@/components/common/table/TablePagination";
import EmptyTable from "@/components/common/table/EmptyTable";
import ContactCardList from "./ContactCardList";

type Props = {
    contacts: Contact[];
    page: number;
    setPage: React.Dispatch<React.SetStateAction<number>>;
    onEdit: (contact: Contact) => void;
    onDelete: (contact: Contact) => void;
};

export default function ContactTable({
    contacts,
    page,
    setPage,
    onEdit,
    onDelete,
}: Props) {
    return (
        <>
            <div className="hidden lg:block">
                <DataTable
                    hasData={contacts.length > 0}
                    columns={
                        <tr className="border-b border-slate-200 text-left">
                            <th className="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500">
                                Contact
                            </th>

                            <th className="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500">
                                Contact Details
                            </th>

                            <th className="px-6 py-4 text-right text-xs font-semibold uppercase tracking-wider text-slate-500">
                                Actions
                            </th>
                        </tr>
                    }
                    emptyState={
                        <EmptyTable
                            title="No Contacts Found"
                            description="Create your first contact to begin managing customer relationships."
                        />
                    }
                    pagination={
                        <TablePagination
                            page={page}
                            onPageChange={setPage}
                        />
                    }
                >
                    {contacts.map((contact) => (
                        <tr
                            key={contact.id}
                            className="transition-colors hover:bg-slate-50"
                        >
                            {/* Contact */}

                            <td className="px-6 py-5">
                                <AvatarCell
                                    name={`${contact.first_name} ${contact.last_name}`}
                                    subtitle={contact.company}
                                />
                            </td>

                            {/* Details */}

                            <td className="px-6 py-5">
                                <div className="space-y-2">
                                    <div className="flex items-center gap-2 text-sm text-slate-700">
                                        <Mail
                                            size={15}
                                            className="text-slate-400"
                                        />

                                        <span className="truncate">
                                            {contact.email}
                                        </span>
                                    </div>

                                    <div className="flex items-center gap-2 text-sm text-slate-500">
                                        <Phone
                                            size={15}
                                            className="text-slate-400"
                                        />

                                        {contact.phone}
                                    </div>

                                    <div className="flex items-center gap-2 text-sm text-slate-500">
                                        <Building2
                                            size={15}
                                            className="text-slate-400"
                                        />

                                        {contact.company}
                                    </div>
                                </div>
                            </td>

                            {/* Actions */}

                            <td className="px-6 py-5 text-right">
                                <TableActionMenu
                                    onEdit={() => onEdit(contact)}
                                    onDelete={() => onDelete(contact)}
                                />
                            </td>
                        </tr>
                    ))}
                </DataTable>
            </div>

            <div className="lg:hidden">

                <ContactCardList
                    contacts={contacts}
                    onEdit={onEdit}
                    onDelete={onDelete}
                />

                {contacts.length > 0 && (
                    <div className="mt-6">
                        <TablePagination
                            page={page}
                            onPageChange={setPage}
                        />
                    </div>
                )}
            </div>
        </>

    );
}