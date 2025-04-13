import { ApiError } from './types';

const CUSTOM_SLUG_ERROR = 'rpc error: code = Unknown desc = custom slug already exists';
const UNKNOWN_ERROR_PREFIX = 'rpc error: code = Unknown desc =';
const UNAUTHORIZED_ERROR = 'Authorization header is required';

export class ApiErrorHandler {
    static handleError(error: unknown): ApiError {
        if (typeof error === 'object' && error !== null) {
            const errorObj = error as Record<string, any>;

            if ('error' in errorObj) {
                const errorMessage = errorObj.error;
                if (typeof errorMessage === 'string') {
                    return this.processError(errorMessage);
                }
            }

            return this.processError(JSON.stringify(error));
        }

        if (error instanceof Error) {
            return this.processError(error.message);
        }

        return this.processError(String(error));
    }

    private static processError(errorMessage: string): ApiError {
        if (errorMessage.includes(CUSTOM_SLUG_ERROR)) {
            return {
                message: 'This custom URL is already taken. Please choose a different one.',
                code: 'CUSTOM_SLUG_EXISTS'
            };
        }

        if (errorMessage.includes(UNAUTHORIZED_ERROR)) {
            return {
                message: 'You need to be logged in to perform this action. Please log in and try again.',
                code: 'UNAUTHORIZED'
            };
        }

        if (errorMessage.includes(UNKNOWN_ERROR_PREFIX)) {
            const match = errorMessage.match(/desc = (.+)$/);
            if (match && match[1]) {
                return { message: match[1] };
            }
        }

        return { message: errorMessage };
    }
} 