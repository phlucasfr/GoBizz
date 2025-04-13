export interface ApiResponse<T = void> {
    data?: T;
    success: boolean;
    message: string;
}

export interface ApiRequest {
    body?: unknown;
    method: 'GET' | 'POST' | 'PUT' | 'DELETE';
    endpoint: string;
    headers?: Record<string, string>;
    isSecure?: boolean;
    credentials?: RequestCredentials;
}

export interface ApiError {
    code?: string;
    message: string;
    details?: unknown;
} 