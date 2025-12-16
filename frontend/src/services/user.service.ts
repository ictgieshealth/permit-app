import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import { ChangePasswordRequest, UpdateProfileRequest, User, UserRequest, UserUpdateRequest } from "@/types/user";

export const userService = {
  async getAll(params?: {
    domain_id?: number;
    role_id?: number;
    username?: string;
    email?: string;
    full_name?: string;
    is_active?: boolean;
    page?: number;
    limit?: number;
  }) {
    const queryParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          queryParams.append(key, String(value));
        }
      });
    }
    
    const endpoint = `/users?${queryParams.toString()}`;
    return await apiClient.get<ApiResponse<User[]>>(endpoint);
  },

  async getById(id: number) {
    const response = await apiClient.get<ApiResponse<User>>(`/users/${id}`);
    return response.data;
  },

  async create(data: UserRequest) {
    const response = await apiClient.post<ApiResponse<User>>('/users', data);
    return response.data;
  },

  async update(id: number, data: UserUpdateRequest) {
    const response = await apiClient.put<ApiResponse<User>>(`/users/${id}`, data);
    return response.data;
  },

  async delete(id: number) {
    return await apiClient.delete<ApiResponse<null>>(`/users/${id}`);
  },

  async changePassword(id: number, data: ChangePasswordRequest) {
    return await apiClient.put<ApiResponse<null>>(
      `/users/${id}/change-password`,
      data
    );
  },

  async addDomainRole(id: number, domainId: number, roleId: number, isDefault: boolean = false) {
    return await apiClient.post<ApiResponse<null>>(
      `/users/${id}/domain-roles`,
      { domain_id: domainId, role_id: roleId, is_default: isDefault }
    );
  },

  async removeDomainRole(id: number, domainId: number, roleId: number) {
    return await apiClient.delete<ApiResponse<null>>(
      `/users/${id}/domain-roles/${domainId}/${roleId}`
    );
  },

  async setDefaultDomainRole(id: number, domainId: number, roleId: number) {
    return await apiClient.put<ApiResponse<null>>(
      `/users/${id}/domain-roles/${domainId}/${roleId}/set-default`,
      {}
    );
  },

  async updateProfile(data: UpdateProfileRequest) {
    const response = await apiClient.put<ApiResponse<User>>('/auth/profile', data);
    return response.data;
  },
};