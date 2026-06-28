import { Task } from "../types";

import StatusBadge from "./StatusBadge";

export default function TaskTable({
    tasks,
    page,
    setPage,
    onEdit,
    onDelete,
}: {
    tasks: Task[];
    page: number;
    setPage: React.Dispatch<React.SetStateAction<number>>;
    onEdit: (task: Task) => void;
    onDelete: (task: Task) => void;
}) {

    return (
        <>
            <table className="w-full border border-gray-200">

                <thead>

                    <tr className="border-b bg-gray-100">

                        <th className="p-3 text-left">

                            Title

                        </th>

                        <th className="p-3 text-left">

                            Status

                        </th>

                        <th className="p-3 text-left">

                            Due Date

                        </th>

                        <th className="p-3 text-left">

                            Actions

                        </th>

                    </tr>

                </thead>

                <tbody>

                    {tasks.map((task) => (

                        <tr
                            key={task.id}
                            className="border-b"
                        >

                            <td className="p-3">

                                {task.title}

                            </td>

                            <td className="p-3">

                                <StatusBadge
                                    status={task.status}
                                />

                            </td>

                            <td className="p-3">

                                {task.due_date
                                    ? new Date(task.due_date).toLocaleString()
                                    : "-"}

                            </td>

                            <td className="p-3">

                                <div className="flex gap-3">

                                    <button
                                        onClick={() => onEdit(task)}
                                        className="text-blue-600 hover:underline"
                                    >

                                        Edit

                                    </button>

                                    <button
                                        onClick={() => onDelete(task)}
                                        className="text-red-600 hover:underline"
                                    >

                                        Delete

                                    </button>

                                </div>

                            </td>

                        </tr>

                    ))}

                    {tasks.length === 0 && (

                        <tr>

                            <td
                                colSpan={4}
                                className="p-6 text-center text-gray-500"
                            >

                                No tasks found.

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