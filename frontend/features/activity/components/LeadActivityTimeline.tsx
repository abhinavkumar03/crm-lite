"use client";

import {
    useEffect,
    useState,
} from "react";

import {
    History,
} from "lucide-react";

import {
    getLeadActivities,
} from "../api";

import {
    Activity,
} from "../types";

import ActivityItem from "./ActivityItem";

type Props = {
    leadId: string;
};

export default function LeadActivityTimeline({
    leadId,
}: Props) {
    const [
        activities,
        setActivities,
    ] = useState<Activity[]>([]);

    const [
        loading,
        setLoading,
    ] = useState(true);

    async function load() {
        try {
            const data =
                await getLeadActivities(
                    leadId
                );

            setActivities(data);
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        load();
    }, []);

    if (loading) {
        return (
            <div className="rounded-3xl border border-slate-200 bg-white p-10 text-center">
                Loading activity...
            </div>
        );
    }

    return (
        <div className="space-y-6">

            {/* Header */}

            <div
                className="
                rounded-3xl
                border
                border-slate-200
                bg-white
                p-6
                shadow-sm
                "
            >
                <h2 className="text-xl font-semibold">
                    Activity Timeline
                    </h2>

                    <p className="mt-1 text-sm text-slate-500">
                    {activities.length} recorded activities
                    </p>

                <p className="mt-1 text-sm text-slate-500">
                    Complete history of changes and interactions for this lead.
                </p>
            </div>

            {/* Empty */}

            {activities.length === 0 && (
                <div
                    className="
                    rounded-3xl
                    border
                    border-dashed
                    border-slate-300
                    bg-white
                    p-14
                    text-center
                    "
                >
                    <div
                        className="
                        mx-auto
                        flex
                        h-20
                        w-20
                        items-center
                        justify-center
                        rounded-full
                        bg-emerald-50
                        "
                    >
                        <History
                            size={36}
                            className="text-emerald-500"
                        />
                    </div>

                    <h3 className="mt-6 text-xl font-semibold">
                        No Activity Yet
                    </h3>

                    <p className="mt-3 text-slate-500">
                        Activity will appear automatically as users interact with this lead.
                    </p>
                </div>
            )}

            {/* Timeline */}

            <div className="space-y-4">
                {activities.map(
                        (
                        activity,
                        index
                        )=>(
                        <ActivityItem
                            key={activity.id}
                            activity={activity}
                            isLast={
                                index ===
                                activities.length - 1
                            }
                        />
                    )
                )}
            </div>

        </div>
    );
}