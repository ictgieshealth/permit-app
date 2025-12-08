import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import { PermitType } from "@/types/permitType";

export const permitTypeService = {
  async getAll(params: { limit?: number } = {}): Promise<
    ApiResponse<PermitType[]> & {
      meta?: { page: number; limit: number; total: number };
    }
  > {
    const queryParams = new URLSearchParams();
    if (params.limit) queryParams.append("limit", params.limit.toString());

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
};
