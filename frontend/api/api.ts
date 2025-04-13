import { getCookie } from '@/context/auth-context';
import { apiConfig } from './config';
import { ApiErrorHandler } from './error-handler';
import { EncryptionService } from './encryption-service';
import { ApiRequest, ApiResponse } from './types';

export async function apiRequest<T>({
    method,
    endpoint,
    body,
    headers = {},
    isSecure = true,
}: ApiRequest): Promise<ApiResponse<T>> {
    try {
        let token: string | undefined;
        if(isSecure) {
            token = getCookie("auth-token");
        }

        const requestBody = prepareRequestBody(method, body);
        const requestHeaders = prepareHeaders(headers, token);

        const response = await makeRequest(endpoint, method, requestHeaders, requestBody);
        const responseData = await response.json();

        if (!response.ok) return handleErrorResponse<T>(responseData);

        return handleSuccessResponse<T>(responseData);
    } catch (error) {
        console.error('API request error:', error);
        return handleError<T>(error);
    }
}

function prepareHeaders(headers: Record<string, string>, token?: string): Record<string, string> {
    const defaultHeaders: Record<string, string> = {
        'Content-Type': 'application/json',
    };

    if (token) defaultHeaders['Authorization'] = `Bearer ${token}`;

    return { ...defaultHeaders, ...headers };
}

function prepareRequestBody(method: string, body?: unknown): string | undefined {
    if (method === 'GET' || !body) {
        return undefined;
    }

    return JSON.stringify(EncryptionService.encryptData(body));
}

async function makeRequest(
    endpoint: string,
    method: string,
    headers: Record<string, string>,
    body?: string
): Promise<Response> {    
    return fetch(`${apiConfig.baseUrl}${endpoint}`, {
        method,
        headers,
        body,
        credentials: 'include',
    });
}

function handleErrorResponse<T>(errorData: unknown): ApiResponse<T> {
    const error = ApiErrorHandler.handleError(errorData);
    return {
        success: false,
        message: error.message,
    } as ApiResponse<T>;
}

function handleSuccessResponse<T>(data: unknown): ApiResponse<T> {
    try {
        const parsedData = EncryptionService.decryptAndParse<T>(data);
        return {
            success: true,
            message: 'Operation completed successfully',
            data: parsedData,
        };
    } catch (error) {
        return {
            success: false,
            message: 'Failed to process response data',
        } as ApiResponse<T>;
    }
}

function handleError<T>(error: unknown): ApiResponse<T> {
    const apiError = ApiErrorHandler.handleError(error);
    return {
        success: false,
        message: apiError.message,
    } as ApiResponse<T>;
}