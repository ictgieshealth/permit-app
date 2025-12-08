import { API_BASE_URL, STORAGE_KEYS } from '@/config/api';

interface RequestOptions extends RequestInit {
  requireAuth?: boolean;
}

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private getHeaders(requireAuth: boolean = true): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

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
    const { requireAuth = true, ...fetchOptions } = options;

    const response = await fetch(`${this.baseURL}${endpoint}`, {
      ...fetchOptions,
      headers: {
        ...this.getHeaders(requireAuth),
        ...fetchOptions.headers,
      },
    });

    const data = await response.json();

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
      
      throw new Error(data.message || `HTTP error! status: ${response.status}`);
    }

    return data;
  }

  async get<T>(endpoint: string, requireAuth: boolean = true): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET', requireAuth });
  }

  async post<T>(
    endpoint: string,
    body?: unknown,
    requireAuth: boolean = true
  ): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
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
      body: body ? JSON.stringify(body) : undefined,
      requireAuth,
    });
  }

  async delete<T>(endpoint: string, requireAuth: boolean = true): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE', requireAuth });
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
