import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import { Domain } from "@/types/domain";

export const domainService = {
  async getAll(params?: {
    code?: string;
    name?: string;
    is_active?: boolean;
    page?: number;
    limit?: number;
  }): Promise<ApiResponse<Domain[]> & { meta?: { page: number; limit: number; total: number } }> {
    const queryParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          queryParams.append(key, String(value));
        }
      });
    }
    
    const endpoint = `/domains?${queryParams.toString()}`;
    return await apiClient.get<ApiResponse<Domain[]>>(endpoint) as any;
  },

  async getById(id: number) {
    const response = await apiClient.get<ApiResponse<Domain>>(`/domains/${id}`);
    return response.data;
  },

  async create(data: { code: string; name: string; description?: string }) {
    const response = await apiClient.post<ApiResponse<Domain>>('/domains', data);
    return response.data;
  },

  async update(id: number, data: { code: string; name: string; description?: string }) {
    const response = await apiClient.put<ApiResponse<Domain>>(`/domains/${id}`, data);
    return response.data;
  },

  async delete(id: number) {
    const response = await apiClient.delete<ApiResponse<void>>(`/domains/${id}`);
    return response.data;
  },
};