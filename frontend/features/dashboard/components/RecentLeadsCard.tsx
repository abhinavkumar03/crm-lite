import { Lead } from "../types";

export default function RecentLeadsCard({
    leads,
}: {
    leads: Lead[];
}) {

    return (

        <div className="rounded-lg border bg-white p-6">

            <h2 className="mb-4 text-lg font-semibold">

                Recent Leads

            </h2>

            <div className="space-y-3">

                {leads.map(lead => (

                    <div key={lead.id}>

                        <p className="font-medium">

                            {lead.name}

                        </p>

                        <p className="text-sm text-gray-500">

                            {lead.company}

                        </p>

                    </div>

                ))}

            </div>

        </div>

    );

}