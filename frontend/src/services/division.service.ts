import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import { Division } from "@/types/division";

export const divisionService = {
  async getAll(params?: {
    domain_id?: number;
    code?: string;
    name?: string;
    is_active?: boolean;
    page?: number;
    limit?: number;
  }): Promise<
    ApiResponse<Division[]> & {
      meta?: { page: number; limit: number; total: number };
    }
  > {
    const queryParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          queryParams.append(key, String(value));
        }
      });
    }

    const response = await apiClient.get<ApiResponse<Division[]>>(
      `/divisions?${queryParams.toString()}`
    );
    return response as any;
  },

  async getById(id: number): Promise<Division> {
    const response = await apiClient.get<ApiResponse<Division>>(
      `/divisions/${id}`
    );
    return response.data;
  },

  async create(data: { domain_id: number; code: string; name: string }) {
    const response = await apiClient.post<ApiResponse<Division>>('/divisions', data);
    return response.data;
  },

  async update(id: number, data: { domain_id: number; code: string; name: string }) {
    const response = await apiClient.put<ApiResponse<Division>>(`/divisions/${id}`, data);
    return response.data;
  },

  async delete(id: number) {
    const response = await apiClient.delete<ApiResponse<void>>(`/divisions/${id}`);
    return response.data;
  },
};
