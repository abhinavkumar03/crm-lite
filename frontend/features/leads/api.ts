import api from "@/services/api";

export async function getLeads(
    page = 1,
    search = ""
) {
    const response = await api.get("/leads", {
        params: {
            page,
            search,
        },
    });

    return response.data;
}