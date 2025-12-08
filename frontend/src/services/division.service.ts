import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import { Division } from "@/types/division";

export const divisionService = {
  async getAll(params: { limit?: number } = {}): Promise<
    ApiResponse<Division[]> & {
      meta?: { page: number; limit: number; total: number };
    }
  > {
    const queryParams = new URLSearchParams();
    if (params.limit) queryParams.append("limit", params.limit.toString());

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
};
