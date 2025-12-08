import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import { Domain } from "domain";

export const domainService = {
  async getAll(params?: {
    code?: string;
    name?: string;
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
    
    const endpoint = `/domains?${queryParams.toString()}`;
    return await apiClient.get<ApiResponse<Domain[]>>(endpoint);
  },

  async getById(id: number) {
    const response = await apiClient.get<ApiResponse<Domain>>(`/domains/${id}`);
    return response.data;
  },
};