import { Lead } from "../types";

import StatusBadge from "./StatusBadge";

export default function LeadTable({
    leads,
    page,
    setPage,
    onEdit,
    onDelete,
}: {
    leads: Lead[];
    page: number;
    setPage: React.Dispatch<React.SetStateAction<number>>;
    onEdit: (lead: Lead) => void;
    onDelete: (lead: Lead) => void;
}) {

    return (
        <>
            <table className="w-full border border-gray-200">

                <thead>

                    <tr className="border-b bg-gray-100">

                        <th className="p-3 text-left">

                            Name

                        </th>

                        <th className="p-3 text-left">

                            Company

                        </th>

                        <th className="p-3 text-left">

                            Status

                        </th>

                        <th className="p-3 text-left">

                            Phone

                        </th>

                        <th className="p-3 text-left">

                            Actions

                        </th>

                    </tr>

                </thead>

                <tbody>

                    {leads.map((lead) => (

                        <tr
                            key={lead.id}
                            className="border-b"
                        >

                            <td className="p-3">

                                {lead.name}

                            </td>

                            <td className="p-3">

                                {lead.company}

                            </td>

                            <td className="p-3">

                                <StatusBadge
                                    status={lead.status}
                                />

                            </td>

                            <td className="p-3">

                                {lead.phone}

                            </td>

                            <td className="p-3">

                                <div className="flex gap-3">

                                    <button
                                        onClick={() => onEdit(lead)}
                                        className="text-blue-600 hover:underline"
                                    >
                                        Edit
                                    </button>

                                    <button
                                        onClick={() => onDelete(lead)}
                                        className="text-red-600 hover:underline"
                                    >
                                        Delete
                                    </button>

                                </div>

                            </td>

                        </tr>

                    ))}

                    {leads.length === 0 && (

                        <tr>

                            <td
                                colSpan={5}
                                className="p-6 text-center text-gray-500"
                            >

                                No leads found.

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