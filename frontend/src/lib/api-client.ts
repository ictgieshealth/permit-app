import { API_BASE_URL, STORAGE_KEYS } from '@/config/api';

interface RequestOptions extends RequestInit {
  requireAuth?: boolean;
  responseType?: 'json' | 'blob';
}

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private getHeaders(requireAuth: boolean = true, isFormData: boolean = false): HeadersInit {
    const headers: HeadersInit = {};

    // Only set Content-Type for JSON requests, not for FormData
    if (!isFormData) {
      headers['Content-Type'] = 'application/json';
    }

    if (requireAuth && typeof window !== 'undefined') {
      const token = localStorage.getItem(STORAGE_KEYS.TOKEN);
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }

    return headers;
  }

  private async request<T>(
    endpoint: string,
    options: RequestOptions = {}
  ): Promise<T> {
    const { requireAuth = true, responseType = 'json', ...fetchOptions } = options;
    
    // Check if body is FormData
    const isFormData = fetchOptions.body instanceof FormData;

    const response = await fetch(`${this.baseURL}${endpoint}`, {
      ...fetchOptions,
      headers: {
        ...this.getHeaders(requireAuth, isFormData),
        ...fetchOptions.headers,
      },
    });

    if (!response.ok) {
      // If unauthorized and it's an authentication error, clear stored credentials
      if (response.status === 401 && typeof window !== 'undefined') {
        localStorage.removeItem(STORAGE_KEYS.TOKEN);
        localStorage.removeItem(STORAGE_KEYS.USER);
        localStorage.removeItem(STORAGE_KEYS.DOMAIN);
        // Redirect to login if not already on login page
        // if (!window.location.pathname.includes('/signin')) {
        //   window.location.href = '/signin';
        // }
      }
      
      const data = await response.json().catch(() => ({}));
      throw new Error(data.message || `HTTP error! status: ${response.status}`);
    }

    if (responseType === 'blob') {
      return response.blob() as T;
    }

    return response.json();
  }

  async get<T>(endpoint: string, options?: { requireAuth?: boolean; responseType?: 'json' | 'blob' }): Promise<T> {
    return this.request<T>(endpoint, { 
      method: 'GET', 
      requireAuth: options?.requireAuth ?? true,
      responseType: options?.responseType ?? 'json'
    });
  }

  async post<T>(
    endpoint: string,
    body?: unknown,
    requireAuth: boolean = true
  ): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: body instanceof FormData ? body : (body ? JSON.stringify(body) : undefined),
      requireAuth,
    });
  }

  async put<T>(
    endpoint: string,
    body?: unknown,
    requireAuth: boolean = true
  ): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: body instanceof FormData ? body : (body ? JSON.stringify(body) : undefined),
      requireAuth,
    });
  }

  async delete<T>(endpoint: string, requireAuth: boolean = true): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE', requireAuth });
  }

  getBaseUrl(): string {
    return this.baseURL;
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
