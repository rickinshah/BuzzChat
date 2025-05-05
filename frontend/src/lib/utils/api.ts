import { writable } from "svelte/store";
import { goto } from "../../stores/navigation";
import { triggerError } from "../../stores/modal";


type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

interface ApiCallOptions<T = unknown> {
    endpoint: string;
    method?: HttpMethod;
    data?: any;
    headers?: Record<string, string>;
    credentials?: RequestCredentials;
    isMultipart?: boolean;
    onSuccess: (response: T) => void;
    onError?: (error: any) => void;
};

type ApiConfig = {
    protocol: "http" | "https";
    host: string;
    port: number;
};

const config: ApiConfig = {
    protocol: "http",
    host: "localhost",
    port: 4000
};

export const BASE_URL = `${config.protocol}://${config.host}:${config.port}`

export async function apiCall<T = unknown>({
    endpoint,
    method = 'GET',
    data = null,
    headers = {},
    credentials,
    isMultipart = false,
    onSuccess,
    onError
}: ApiCallOptions<T>): Promise<void> {
    try {
        const fetchOptions: RequestInit = {
            method,
            headers: isMultipart ? headers : {
                'Content-Type': 'application/json',
                ...headers
            },
            body: data ? (isMultipart ? data : JSON.stringify(data)) : undefined
        };

        if (credentials) {
            fetchOptions.credentials = credentials;
        }

        const response = await fetch(`${BASE_URL}${endpoint}`, fetchOptions);
        const responseData = await response.json();

        if (response.ok) {
            onSuccess(responseData as T)
        } else if (response.status === 401) {
            localStorage.clear();
            goto.set("/");
            triggerError(responseData.error)
        } else if (response.status === 422) {
            triggerError(responseData.error)
        }
        else {
            if (onError) {
                onError(responseData)
            } else {
                triggerError(responseData.error);
            }
        }

    } catch (err) {
        console.error(err)
    }
}
