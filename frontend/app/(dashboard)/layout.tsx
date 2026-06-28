"use client";

import { useEffect } from "react";

import { useRouter } from "next/navigation";

import { useAuth } from "@/context/AuthContext";

import Sidebar from "@/components/layout/Sidebar";

import Topbar from "@/components/layout/Topbar";

export default function DashboardLayout({
    children,
}:{
    children:React.ReactNode
}){

    const auth = useAuth();

    const router = useRouter();

    useEffect(()=>{

        if(!auth.token){

            router.replace("/login");

        }

    },[
        auth.token,
        router,
    ]);

    if(!auth.token){

        return null;

    }

    return(

        <div className="flex h-screen">

            <Sidebar/>

            <div className="flex flex-1 flex-col">

                <Topbar/>

                <main className="flex-1 overflow-auto p-6">

                    {children}

                </main>

            </div>

        </div>

    );

}