import axios from "axios";

const api = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_URL,
});

api.interceptors.request.use((config) => {
    // Token lives in localStorage, which is only available in the browser.
    // Guard the access so the client is safe to import in server components.
    if (typeof window !== "undefined") {
        const token = localStorage.getItem("token");
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
            // Only attach workspace scope on authenticated requests so login/register
            // preflights are not blocked by CORS when a stale org id is in storage.
            const orgId = localStorage.getItem("active_organization_id");
            if (orgId) {
                config.headers["X-Organization-Id"] = orgId;
            }
        }
    }

    return config;
});

api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (typeof window !== "undefined") {
            const message = error?.response?.data?.message as string | undefined;
            const status = error?.response?.status;
            const path = window.location.pathname;
            if (
                status === 403 &&
                message === "User does not belong to an organization" &&
                !path.startsWith("/onboarding") &&
                !path.startsWith("/login") &&
                !path.startsWith("/register")
            ) {
                window.location.replace("/onboarding/organization");
            }
        }
        return Promise.reject(error);
    }
);

export default api;
