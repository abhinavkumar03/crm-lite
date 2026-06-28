import { DashboardResponse } from "../types";

type Props = {
    dashboard: DashboardResponse;
};

export default function MetricsGrid({
    dashboard,
}: Props) {

    const cards = [

        {
            title: "Total Leads",
            value: dashboard.total_leads,
        },

        {
            title: "Contacts",
            value: dashboard.total_contacts,
        },

        {
            title: "Tasks",
            value: dashboard.total_tasks,
        },

        {
            title: "Won Leads",
            value: dashboard.won_leads,
        },

    ];

    return (

        <div className="grid grid-cols-4 gap-6">

            {cards.map(card => (

                <div
                    key={card.title}
                    className="rounded-lg border bg-white p-6 shadow-sm"
                >

                    <p className="text-sm text-gray-500">

                        {card.title}

                    </p>

                    <h2 className="mt-2 text-3xl font-bold">

                        {card.value}

                    </h2>

                </div>

            ))}

        </div>

    );

}