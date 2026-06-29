"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

import { login } from "@/features/auth/api";
import { useAuth } from "@/context/AuthContext";

export default function LoginPage() {

    const router = useRouter();

    const auth = useAuth();

    const [email, setEmail] = useState("");

    const [password, setPassword] =
        useState("");

    const [loading, setLoading] =
        useState(false);

    async function handleLogin(
        e: React.FormEvent,
    ) {

        e.preventDefault();

        try {
            setLoading(true);

            const res = await login(email, password);
            console.log("Login response:", res);

            console.log("Before auth.login");
            await auth.login(res.data.access_token);
            console.log("After auth.login");

            console.log("Before router.push");
            router.push("/dashboard");

        } catch (err) {
            console.error(err);
            alert("Invalid credentials");
        } finally {
            setLoading(false);
        }
    }

    return (

        <div className="flex min-h-screen items-center justify-center">

            <form
                onSubmit={handleLogin}
                className="w-full max-w-md rounded border bg-white p-8 shadow"
            >

                <h1 className="mb-6 text-2xl font-bold">

                    CRM Lite

                </h1>

                <input
                    className="mb-4 w-full rounded border p-3"
                    placeholder="Email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                />

                <input
                    type="password"
                    className="mb-6 w-full rounded border p-3"
                    placeholder="Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />

                <button
                    disabled={loading}
                    className="w-full rounded bg-blue-600 py-3 text-white"
                >

                    {loading
                        ? "Logging in..."
                        : "Login"}

                </button>

            </form>

        </div>

    );
}