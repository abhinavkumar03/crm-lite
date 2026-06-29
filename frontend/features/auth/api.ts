import api from "@/services/api";

import {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
  ProfileResponse,
} from "./types";

/* ---------------- Login ---------------- */

export async function login(
  payload: LoginRequest
): Promise<LoginResponse> {
  const { data } =
    await api.post<LoginResponse>(
      "/auth/login",
      payload
    );

  return data;
}

/* ---------------- Register ---------------- */

export async function register(
  payload: RegisterRequest
): Promise<RegisterResponse> {
  const { data } =
    await api.post<RegisterResponse>(
      "/auth/register",
      payload
    );

  return data;
}

/* ---------------- Profile ---------------- */

export async function getProfile(): Promise<ProfileResponse> {
  const { data } =
    await api.get<ProfileResponse>(
      "/auth/profile"
    );

  return data;
}