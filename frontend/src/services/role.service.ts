import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import { Role } from "@/types/role";

export const roleService = {
  async getAll(params?: { page?: number; limit?: number; is_active?: boolean }) {
    const queryParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          queryParams.append(key, String(value));
        }
      });
    }
    
    const endpoint = `/roles?${queryParams.toString()}`;
    return await apiClient.get<ApiResponse<Role[]>>(endpoint);
  },

  async getById(id: number) {
    const response = await apiClient.get<ApiResponse<Role>>(`/roles/${id}`);
    return response.data;
  },
};