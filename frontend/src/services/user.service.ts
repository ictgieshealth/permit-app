import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import { ChangePasswordRequest, User, UserRequest, UserUpdateRequest } from "@/types/user";

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

  async addDomain(id: number, domainId: number, isDefault: boolean = false) {
    return await apiClient.post<ApiResponse<null>>(
      `/users/${id}/domains`,
      { domain_id: domainId, is_default: isDefault }
    );
  },

  async removeDomain(id: number, domainId: number) {
    return await apiClient.delete<ApiResponse<null>>(
      `/users/${id}/domains/${domainId}`
    );
  },

  async setDefaultDomain(id: number, domainId: number) {
    return await apiClient.put<ApiResponse<null>>(
      `/users/${id}/domains/${domainId}/set-default`,
      {}
    );
  },
};