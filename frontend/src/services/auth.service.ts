import { apiClient } from '@/lib/api-client';
import { STORAGE_KEYS } from '@/config/api';
import { LoginRequest, LoginResponse, SwitchDomainRequest, SwitchDomainResponse } from '@/types/loginUser';
import { ApiResponse } from '@/types/api';
import { User, UserDomainRole } from '@/types/user';
import { Domain } from '@/types/domain';
import { Role } from '@/types/role';

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
      // Store token, user data, current domain/role, and all domains
      localStorage.setItem(STORAGE_KEYS.TOKEN, loginData.token);
      localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(loginData.user));
      localStorage.setItem(STORAGE_KEYS.DOMAIN, JSON.stringify(loginData.current_domain));
      localStorage.setItem('current_role', JSON.stringify(loginData.current_role));
      localStorage.setItem('user_domains', JSON.stringify(loginData.domains));
    }
    
    return loginData;
  },

  async switchDomain(domainId: number): Promise<SwitchDomainResponse> {
    const request: SwitchDomainRequest = { domain_id: domainId };
    const response = await apiClient.post<ApiResponse<SwitchDomainResponse>>(
      '/auth/switch-domain',
      request
    );

    const switchData = response.data;
    
    if (switchData) {
      // Update token and current domain/role
      localStorage.setItem(STORAGE_KEYS.TOKEN, switchData.token);
      localStorage.setItem(STORAGE_KEYS.DOMAIN, JSON.stringify(switchData.current_domain));
      localStorage.setItem('current_role', JSON.stringify(switchData.current_role));
    }
    
    return switchData;
  },

  async logout(): Promise<void> {
    // Clear localStorage
    localStorage.removeItem(STORAGE_KEYS.TOKEN);
    localStorage.removeItem(STORAGE_KEYS.USER);
    localStorage.removeItem(STORAGE_KEYS.DOMAIN);
    localStorage.removeItem('current_role');
    localStorage.removeItem('user_domains');
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

  getStoredRole(): Role | null {
    if (typeof window === 'undefined') return null;
    const roleStr = localStorage.getItem('current_role');
    return roleStr ? JSON.parse(roleStr) : null;
  },

  getStoredDomains(): UserDomainRole[] {
    if (typeof window === 'undefined') return [];
    const domainsStr = localStorage.getItem('user_domains');
    return domainsStr ? JSON.parse(domainsStr) : [];
  },

  isAuthenticated(): boolean {
    if (typeof window === 'undefined') return false;
    return !!localStorage.getItem(STORAGE_KEYS.TOKEN);
  },
  
  hasRole(allowedRoleCodes: string[]): boolean {
    if (typeof window === 'undefined') return false;
    const role = this.getStoredRole();
    if (!role) return false;
    return allowedRoleCodes.includes(role.code);
  },

};
