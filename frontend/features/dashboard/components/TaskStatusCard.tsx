import { DashboardResponse } from "../types";

export default function TaskStatusCard({
    dashboard,
}: {
    dashboard: DashboardResponse;
}) {

    return (

        <div className="rounded-lg border bg-white p-6">

            <h2 className="mb-4 text-lg font-semibold">

                Task Status

            </h2>

            <div className="space-y-3">

                <p>Pending: {dashboard.pending_tasks}</p>

                <p>In Progress: {dashboard.in_progress_tasks}</p>

                <p>Completed: {dashboard.completed_tasks}</p>

            </div>

        </div>

    );

}