import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import { PermitType } from "@/types/permitType";

export const permitTypeService = {
  async getAll(params?: {
    code?: string;
    name?: string;
    page?: number;
    limit?: number;
  }): Promise<
    ApiResponse<PermitType[]> & {
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

    const response = await apiClient.get<ApiResponse<PermitType[]>>(
      `/permit-types?${queryParams.toString()}`
    );
    return response as any;
  },

  async getById(id: number): Promise<PermitType> {
    const response = await apiClient.get<ApiResponse<PermitType>>(
      `/permit-types/${id}`
    );
    return response.data;
  },

  async create(data: { code: string; name: string; description?: string }) {
    const response = await apiClient.post<ApiResponse<PermitType>>('/permit-types', data);
    return response.data;
  },

  async update(id: number, data: { code: string; name: string; description?: string }) {
    const response = await apiClient.put<ApiResponse<PermitType>>(`/permit-types/${id}`, data);
    return response.data;
  },

  async delete(id: number) {
    const response = await apiClient.delete<ApiResponse<void>>(`/permit-types/${id}`);
    return response.data;
  },
};
