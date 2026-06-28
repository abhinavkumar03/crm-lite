import { DashboardResponse } from "../types";

export default function LeadStatusCard({
    dashboard,
}: {
    dashboard: DashboardResponse;
}) {

    return (

        <div className="rounded-lg border bg-white p-6">

            <h2 className="mb-4 text-lg font-semibold">

                Lead Status

            </h2>

            <div className="space-y-3">

                <p>New: {dashboard.new_leads}</p>

                <p>Contacted: {dashboard.contacted_leads}</p>

                <p>Qualified: {dashboard.qualified_leads}</p>

                <p>Won: {dashboard.won_leads}</p>

                <p>Lost: {dashboard.lost_leads}</p>

            </div>

        </div>

    );

}