export interface Lead {
    id: string;
    name: string;
    email: string;
    phone: string;
    company: string;
    status: string;
    notes: string;
}

export interface LeadListResponse {
    data: Lead[];
}