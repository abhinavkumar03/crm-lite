export interface SearchLead {
    id: string;
    name: string;
    company: string;
    status: string;
}

export interface SearchContact {
    id: string;
    name: string;
    email: string;
}

export interface SearchTask {
    id: string;
    title: string;
    status: string;
}

export interface SearchResponse {
    leads: SearchLead[];
    contacts: SearchContact[];
    tasks: SearchTask[];
}