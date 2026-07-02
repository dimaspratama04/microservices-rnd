export interface WebResponse<T = any> {
    message: string;
    data?: T;
}

export function successResponse<T>(message: string, data?: T): WebResponse<T> {
    if (data !== undefined) {
        return { message, data };
    }
    return { message };
}

export function errorResponse(message: string): WebResponse {
    return { message };
}
