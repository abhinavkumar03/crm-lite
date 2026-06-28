import { Task } from "../types";

export default function UpcomingTasksCard({
    tasks,
}: {
    tasks: Task[];
}) {

    return (

        <div className="rounded-lg border bg-white p-6">

            <h2 className="mb-4 text-lg font-semibold">

                Upcoming Tasks

            </h2>

            <div className="space-y-3">

                {tasks.map(task => (

                    <div key={task.id}>

                        <p className="font-medium">

                            {task.title}

                        </p>

                        <p className="text-sm text-gray-500">

                            {task.status}

                        </p>

                    </div>

                ))}

            </div>

        </div>

    );

}