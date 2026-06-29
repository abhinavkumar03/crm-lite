"use client";

import { getProfile } from "@/features/auth/api";
import {
    createContext,
    useContext,
    useEffect,
    useState,
} from "react";

type User = {
    id: string;
    name: string;
    email: string;
};

type AuthContextType = {
    token: string | null;
    user: User | null;
    loading: boolean;
    login: (token: string) => Promise<void>;
    logout: () => void;
};


const AuthContext = createContext<AuthContextType>(
    {} as AuthContextType
);


export function AuthProvider({
    children,
}: {
    children: React.ReactNode;
}) {

    const [token, setToken] = useState<string | null>(null);

    const [user, setUser] = useState<User | null>(null);

    const [loading, setLoading] = useState(true);

    useEffect(() => {

        async function loadUser() {

            const stored = localStorage.getItem("token");

            if (!stored) {
                setLoading(false);
                return;
            }

            setToken(stored);

            try {

                const profile = await getProfile();

                setUser(profile.data);

            } catch {

                logout();

            } finally {

                setLoading(false);

            }

        }

        loadUser();

    }, []);

    async function login(jwt: string) {
        console.log("Received JWT:", jwt);

        localStorage.setItem("token", jwt);

        console.log(
            "Stored JWT:",
            localStorage.getItem("token")
        );

        setToken(jwt);

        try {
            console.log("Calling profile...");

            const profile = await getProfile();

            console.log("Profile response:", profile);

            setUser(profile.data);

        } catch (err) {
            console.error("Profile failed", err);
            throw err;
        }
    }

    function logout() {

        localStorage.removeItem(
            "token",
        );

        setUser(null);

        setToken(null);
    }

    return (
        <AuthContext.Provider
            value={{
                token,
                user,
                loading,
                login,
                logout,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    return useContext(AuthContext);
}