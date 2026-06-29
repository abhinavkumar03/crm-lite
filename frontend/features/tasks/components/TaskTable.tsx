import {
    CalendarDays,
    FileText,
    Link2,
} from "lucide-react";

import { Task } from "../types";

import DataTable from "@/components/common/table/DataTable";
import StatusBadge from "@/components/common/table/StatusBadge";
import TableActionMenu from "@/components/common/table/TableActionMenu";
import TablePagination from "@/components/common/table/TablePagination";
import EmptyTable from "@/components/common/table/EmptyTable";
import TaskCardList from "./TaskCardList";

type Props = {
    tasks: Task[];
    page: number;
    setPage: React.Dispatch<React.SetStateAction<number>>;
    onEdit: (task: Task) => void;
    onDelete: (task: Task) => void;
};

export default function TaskTable({
    tasks,
    page,
    setPage,
    onEdit,
    onDelete,
}: Props) {
    return (
        <>

            {/* Desktop */}
            <div className="hidden lg:block">
                <DataTable
                    hasData={tasks.length > 0}
                    columns={
                        <tr className="border-b border-slate-200 text-left">
                            <th className="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500">
                                Task
                            </th>

                            <th className="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500">
                                Status
                            </th>

                            <th className="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500">
                                Due Date
                            </th>

                            <th className="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500">
                                Linked
                            </th>

                            <th className="px-6 py-4 text-right text-xs font-semibold uppercase tracking-wider text-slate-500">
                                Actions
                            </th>
                        </tr>
                    }
                    emptyState={
                        <EmptyTable
                            title="No Tasks Found"
                            description="Create your first task to start tracking customer activities."
                        />
                    }
                    pagination={
                        <TablePagination
                            page={page}
                            onPageChange={setPage}
                        />
                    }
                >
                    {tasks.map((task) => (
                        <tr
                            key={task.id}
                            className="transition-colors hover:bg-slate-50"
                        >
                            {/* Task */}

                            <td className="px-6 py-5">
                                <div className="flex items-start gap-3">
                                    <div className="mt-1 rounded-xl bg-emerald-100 p-2">
                                        <FileText
                                            size={18}
                                            className="text-emerald-600"
                                        />
                                    </div>

                                    <div>
                                        <h3 className="font-semibold text-slate-900">
                                            {task.title}
                                        </h3>

                                        <p className="mt-1 line-clamp-2 text-sm text-slate-500">
                                            {task.description || "No description provided"}
                                        </p>
                                    </div>
                                </div>
                            </td>

                            {/* Status */}

                            <td className="px-6 py-5">
                                <StatusBadge
                                    status={task.status}
                                />
                            </td>

                            {/* Due Date */}

                            <td className="px-6 py-5">
                                <div className="flex items-center gap-2 text-sm text-slate-600">
                                    <CalendarDays
                                        size={16}
                                        className="text-slate-400"
                                    />

                                    {task.due_date
                                        ? new Date(task.due_date).toLocaleDateString()
                                        : "-"}
                                </div>
                            </td>

                            {/* Related Entity */}

                            <td className="px-6 py-5">
                                <div className="flex items-center gap-2 text-sm text-slate-600">
                                    <Link2
                                        size={16}
                                        className="text-slate-400"
                                    />

                                    {task.lead_id
                                        ? `Lead #${task.lead_id}`
                                        : task.contact_id
                                            ? `Contact #${task.contact_id}`
                                            : "-"}
                                </div>
                            </td>

                            {/* Actions */}

                            <td className="px-6 py-5 text-right">
                                <TableActionMenu
                                    onEdit={() => onEdit(task)}
                                    onDelete={() => onDelete(task)}
                                />
                            </td>
                        </tr>
                    ))}
                </DataTable>  </div>


            {/* Mobile */}
            <div className="lg:hidden">
                <TaskCardList
                    tasks={tasks}
                    onEdit={onEdit}
                    onDelete={onDelete}
                />

                {tasks.length > 0 && (
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