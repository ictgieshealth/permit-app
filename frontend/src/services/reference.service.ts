import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import {
  Reference,
  ReferenceRequest,
  ReferenceListRequest,
} from "@/types/reference";

export const referenceService = {
  async getAll(
    params: ReferenceListRequest = {}
  ): Promise<
    ApiResponse<Reference[]> & {
      meta?: { page: number; limit: number; total: number };
    }
  > {
    const queryParams = new URLSearchParams();

    if (params.reference_category_id)
      queryParams.append(
        "reference_category_id",
        params.reference_category_id.toString()
      );
    if (params.module_id)
      queryParams.append("module_id", params.module_id.toString());
    if (params.name) queryParams.append("name", params.name);
    if (params.is_active !== undefined)
      queryParams.append("is_active", params.is_active.toString());
    if (params.page) queryParams.append("page", params.page.toString());
    if (params.limit) queryParams.append("limit", params.limit.toString());

    const response = await apiClient.get<ApiResponse<Reference[]>>(
      `/references?${queryParams.toString()}`
    );
    return response as any;
  },

  async getById(id: number): Promise<Reference> {
    const response = await apiClient.get<ApiResponse<Reference>>(
      `/references/${id}`
    );
    return response.data;
  },

  async getByCategoryId(categoryId: number): Promise<Reference[]> {
    const response = await apiClient.get<ApiResponse<Reference[]>>(
      `/reference-categories/${categoryId}/references`
    );
    return response.data;
  },

  async getByModuleId(moduleId: number): Promise<Reference[]> {
    const response = await apiClient.get<ApiResponse<Reference[]>>(
      `/modules/${moduleId}/references`
    );
    return response.data;
  },

  async create(data: ReferenceRequest): Promise<Reference> {
    const response = await apiClient.post<ApiResponse<Reference>>(
      "/references",
      data
    );
    return response.data;
  },

  async update(id: number, data: ReferenceRequest): Promise<Reference> {
    const response = await apiClient.put<ApiResponse<Reference>>(
      `/references/${id}`,
      data
    );
    return response.data;
  },

  async delete(id: number): Promise<void> {
    await apiClient.delete(`/references/${id}`);
  },
};

