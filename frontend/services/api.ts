import axios from "axios";

const api = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_URL,
});

api.interceptors.request.use(config => {
    console.log("Base URL:", config.baseURL);
    console.log("Request URL:", config.url);
    console.log("Final URL:", `${config.baseURL}${config.url}`);

    const token = localStorage.getItem(
        "token",
    );

    if (token) {

        config.headers.Authorization =
            `Bearer ${token}`;
    }

    return config;
});

export default api;