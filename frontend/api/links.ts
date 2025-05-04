import { apiRequest } from './api';
import { apiConfig } from './config';

export interface Link {
    id: string;
    clicks: number;
    short_url: string;
    custom_slug: string;
    created_at: string;
    updated_at: string;
    customer_id: string;
    original_url: string;
    expiration_date?: string;
}

export interface CreateLinkRequest {
    customer_id: string;
    original_url: string;
    custom_slug?: string;
    expiration_date?: string;
}

export interface UpdateLinkRequest {
    id: string;
    customer_id: string;
    custom_slug?: string;
    original_url?: string;
    expiration_date?: string;
}

export interface PaginatedResponse<T> {
    data: T[];
    page: number;
    total: number;
    per_page: number;
    total_pages: number;
}

export interface GetCustomerLinksParams {
    limit?: number;
    offset?: number;
    sort_by?: keyof Link;
    sort_direction?: 'asc' | 'desc';
    search?: string;
    status?: 'all' | 'active' | 'expired';
    customerId: string;
    slug_type?: 'all' | 'custom' | 'auto';
}

export interface GetCustomerLinksResponse {
    links: Link[];
    total: number;
}

export interface GetCustomerLinksApiResponse {
    data: {
        data: Link[];
        total: number;
        page: number;
        per_page: number;
        total_pages: number;
    };
    success: boolean;
}

export const linksApi = {
    createLink: async (data: CreateLinkRequest) => {
        return apiRequest<Link>({
            method: 'POST',
            endpoint: apiConfig.endpoints.links.createLink,
            body: data,
        });
    },

    getLink: async (shortUrl: string) => {
        return apiRequest<Link>({
            method: 'GET',
            endpoint: apiConfig.endpoints.links.getLink.replace(':shortUrl', shortUrl),
        });
    },

    getCustomerLinks: async ({
        customerId,
        limit,
        offset,
        sort_by,
        sort_direction,
        search,
        status,
        slug_type
    }: GetCustomerLinksParams) => {
        const queryParams = new URLSearchParams();

        // Always include search parameter if it exists
        if (search !== undefined) {
            queryParams.append('search', search);
        }

        if (limit) queryParams.append('limit', limit.toString());
        if (offset) queryParams.append('offset', offset.toString());
        if (status) queryParams.append('status', status);
        if (sort_by) queryParams.append('sort_by', sort_by);
        if (slug_type) queryParams.append('slug_type', slug_type);
        if (sort_direction) queryParams.append('sort_direction', sort_direction);

        const response = await apiRequest<GetCustomerLinksResponse>({
            method: 'GET',
            endpoint: `${apiConfig.endpoints.links.getCustomerLinks.replace(':customerId', customerId)}?${queryParams.toString()}`,
        });

        if (response.success && response.data) {
            const apiResponse: GetCustomerLinksApiResponse = {
                success: true,
                data: {
                    data: Array.isArray(response.data.links) ? response.data.links : [],
                    total: response.data.total || 0,
                    page: 1,
                    per_page: limit || 0,
                    total_pages: limit ? Math.ceil((response.data.total || 0) / limit) : 1
                }
            };
            return apiResponse;
        }

        return {
            success: false,
            data: {
                data: [],
                total: 0,
                page: 1,
                per_page: limit || 0,
                total_pages: 0
            }
        };
    },

    deleteLink: async (id: string) => {
        return apiRequest<{ success: boolean }>({
            method: 'DELETE',
            endpoint: `${apiConfig.endpoints.links.deleteLink}/${id}`,
        });
    },

    updateLink: async (data: UpdateLinkRequest) => {
        return apiRequest<Link>({
            method: 'PUT',
            endpoint: apiConfig.endpoints.links.updateLink.replace(':id', data.id),
            body: data,
        });
    },

    updateLinkClicks: async (id: string) => {
        return apiRequest<Link>({
            method: 'PUT',
            endpoint: apiConfig.endpoints.links.updateLinkClicks.replace(':id', id),
        });
    },
}; 