import { apiClient } from '@/lib/api-client';
import { STORAGE_KEYS } from '@/config/api';
import { LoginRequest, LoginResponse } from '@/types/loginUser';
import { ApiResponse } from '@/types/api';
import { User } from '@/types/user';
import { Domain } from '@/types/domain';

export const authService = {
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    // Call backend API directly
    const response = await apiClient.post<ApiResponse<LoginResponse>>(
      '/auth/login',
      credentials,
      false // No auth required for login
    );

    // Backend returns data wrapped in ApiResponse, so access response.data
    const loginData = response.data;
    
    if (loginData) {
      // Store token, user data and domain in localStorage
      localStorage.setItem(STORAGE_KEYS.TOKEN, loginData.token);
      localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(loginData.user));
      
      if (loginData.default_domain) {
        localStorage.setItem(
          STORAGE_KEYS.DOMAIN,
          JSON.stringify(loginData.default_domain)
        );
      }
    }
    
    return loginData;
  },

  async logout(): Promise<void> {
    // Clear localStorage
    localStorage.removeItem(STORAGE_KEYS.TOKEN);
    localStorage.removeItem(STORAGE_KEYS.USER);
    localStorage.removeItem(STORAGE_KEYS.DOMAIN);
  },  

  async getProfile(): Promise<User> {
    const response = await apiClient.get<ApiResponse<User>>('/auth/profile');
    return response.data;
  },

  getStoredUser(): User | null {
    if (typeof window === 'undefined') return null;
    const userStr = localStorage.getItem(STORAGE_KEYS.USER);
    return userStr ? JSON.parse(userStr) : null;
  },

  getStoredDomain(): Domain | null {
    if (typeof window === 'undefined') return null;
    const domainStr = localStorage.getItem(STORAGE_KEYS.DOMAIN);
    return domainStr ? JSON.parse(domainStr) : null;
  },

  isAuthenticated(): boolean {
    if (typeof window === 'undefined') return false;
    return !!localStorage.getItem(STORAGE_KEYS.TOKEN);
  },
  
  hasRole(allowedRoles: string[]): boolean {
    if (typeof window === 'undefined') return false;
    const user = this.getStoredUser();
    if (!user || !user.role) return false;
    return allowedRoles.includes(user.role.name);
  },

};
