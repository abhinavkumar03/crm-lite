"use client";

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
    login: (token: string) => void;
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

    useEffect(() => {

        const stored = localStorage.getItem("token");

        if (stored) {
            setToken(stored);
        }

    }, []);

    function login(jwt: string) {

        localStorage.setItem(
            "token",
            jwt,
        );

        setToken(jwt);
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