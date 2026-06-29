import { Contact } from "../types";

export default function ContactTable({
    contacts,
    page,
    setPage,
    onEdit,
    onDelete,
}: {
    contacts: Contact[];
    page: number;
    setPage: React.Dispatch<React.SetStateAction<number>>;
    onEdit: (contact: Contact) => void;
    onDelete: (contact: Contact) => void;
}) {

    return (
        <>
            <table className="w-full border border-gray-200">

                <thead>

                    <tr className="border-b bg-gray-100">

                        <th className="p-3 text-left">

                            First Name

                        </th>

                        <th className="p-3 text-left">

                            Last Name

                        </th>

                        <th className="p-3 text-left">

                            Email

                        </th>

                        <th className="p-3 text-left">

                            Phone

                        </th>

                        <th className="p-3 text-left">

                            Company

                        </th>

                        <th className="p-3 text-left">

                            Actions

                        </th>

                    </tr>

                </thead>

                <tbody>

                    {contacts?.map((contact) => (

                        <tr
                            key={contact.id}
                            className="border-b"
                        >

                            <td className="p-3">

                                {contact.first_name}

                            </td>

                            <td className="p-3">

                                {contact.last_name}

                            </td>

                            <td className="p-3">

                                {contact.email}

                            </td>

                            <td className="p-3">

                                {contact.phone}

                            </td>

                            <td className="p-3">

                                {contact.company}

                            </td>

                            <td className="p-3">

                                <div className="flex gap-3">

                                    <button
                                        onClick={() => onEdit(contact)}
                                        className="text-blue-600 hover:underline"
                                    >

                                        Edit

                                    </button>

                                    <button
                                        onClick={() => onDelete(contact)}
                                        className="text-red-600 hover:underline"
                                    >

                                        Delete

                                    </button>

                                </div>

                            </td>

                        </tr>

                    ))}

                    {(contacts ?? []).length === 0 && (

                        <tr>

                            <td
                                colSpan={6}
                                className="p-6 text-center text-gray-500"
                            >

                                No contacts found.

                            </td>

                        </tr>

                    )}

                </tbody>

            </table>

            <div className="mt-4 flex items-center justify-end gap-3">

                <button
                    onClick={() =>
                        setPage((p) => Math.max(1, p - 1))
                    }
                    disabled={page === 1}
                    className="rounded border px-4 py-2 disabled:cursor-not-allowed disabled:opacity-50"
                >

                    Previous

                </button>

                <span className="font-medium">

                    {page}

                </span>

                <button
                    onClick={() =>
                        setPage((p) => p + 1)
                    }
                    className="rounded border px-4 py-2"
                >

                    Next

                </button>

            </div>
        </>
    );
}