import api from "@/services/api";

export async function login(
    email: string,
    password: string,
) {

    const response = await api.post(
        "/auth/login",
        {
            email,
            password,
        },
    );

    return response.data;
}

export async function getProfile() {
    const response = await api.get("/auth/profile");
    return response.data;
}
