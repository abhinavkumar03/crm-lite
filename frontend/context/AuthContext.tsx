"use client";

import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react";

import {
  getProfile,
} from "@/features/auth/api";

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

const AuthContext =
  createContext<AuthContextType>(
    {} as AuthContextType
  );

export function AuthProvider({
  children,
}: {
  children: ReactNode;
}) {
  const [token, setToken] =
    useState<string | null>(null);

  const [user, setUser] =
    useState<User | null>(null);

  const [loading, setLoading] =
    useState(true);

  useEffect(() => {
    loadUser();
  }, []);

  async function loadUser() {
    const stored =
      localStorage.getItem("token");

    if (!stored) {
      setLoading(false);
      return;
    }

    setToken(stored);

    try {
      const profile =
        await getProfile();

      setUser(profile.data);
    } catch {
      logout();
    } finally {
      setLoading(false);
    }
  }

  async function login(jwt: string) {
    localStorage.setItem(
      "token",
      jwt
    );

    setToken(jwt);

    const profile =
      await getProfile();

    setUser(profile.data);
  }

  function logout() {
    localStorage.removeItem(
      "token"
    );

    setToken(null);

    setUser(null);
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