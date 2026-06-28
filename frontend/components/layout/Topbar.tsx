"use client";

import { useRouter } from "next/navigation";

import { useAuth } from "@/context/AuthContext";

export default function Topbar(){

    const auth = useAuth();

    const router = useRouter();

    function logout(){

        auth.logout();

        router.replace("/login");

    }

    return(

        <header className="flex h-16 items-center justify-between border-b bg-white px-6">

            <h1 className="text-lg font-semibold">

                CRM Lite

            </h1>

            <button
                onClick={logout}
                className="rounded bg-red-500 px-4 py-2 text-white"
            >

                Logout

            </button>

        </header>

    );

}