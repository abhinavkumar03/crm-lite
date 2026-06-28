import Link from "next/link";

export default function Sidebar() {
    return (
        <aside className="w-64 border-r bg-white">

            <div className="p-6 text-xl font-bold">

                CRM Lite

            </div>

            <nav className="flex flex-col">

                <Link href="/dashboard" className="p-4 hover:bg-gray-100">
                    Dashboard
                </Link>

                <Link href="/leads" className="p-4 hover:bg-gray-100">
                    Leads
                </Link>

                <Link href="/contacts" className="p-4 hover:bg-gray-100">
                    Contacts
                </Link>

                <Link href="/tasks" className="p-4 hover:bg-gray-100">
                    Tasks
                </Link>

            </nav>

        </aside>
    );
}