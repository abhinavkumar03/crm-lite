export interface Contact {
    id: string;
    first_name: string;
    last_name: string;
    email: string;
    phone: string;
    company: string;
}

export interface CreateContactPayload {
    first_name: string;
    last_name: string;
    email: string;
    phone: string;
    company: string;
}