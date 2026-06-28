export interface User {
    id: string;
    name: string;
    email: string;
}

export interface LoginResponse {
    token: string;
}

export interface ProfileResponse {
    id: string;
    name: string;
    email: string;
}