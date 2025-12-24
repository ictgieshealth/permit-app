import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import {
  ReferenceCategory,
  ReferenceCategoryRequest,
  ReferenceCategoryListRequest,
} from "@/types/reference";

export const referenceCategoryService = {
  async getAll(
    params: ReferenceCategoryListRequest = {}
  ): Promise<
    ApiResponse<ReferenceCategory[]> & {
      meta?: { page: number; limit: number; total: number };
    }
  > {
    const queryParams = new URLSearchParams();

    if (params.module_id)
      queryParams.append("module_id", params.module_id.toString());
    if (params.name) queryParams.append("name", params.name);
    if (params.is_active !== undefined)
      queryParams.append("is_active", params.is_active.toString());
    if (params.page) queryParams.append("page", params.page.toString());
    if (params.limit) queryParams.append("limit", params.limit.toString());

    const response = await apiClient.get<ApiResponse<ReferenceCategory[]>>(
      `/reference-categories?${queryParams.toString()}`
    );
    return response as any;
  },

  async getById(id: number): Promise<ReferenceCategory> {
    const response = await apiClient.get<ApiResponse<ReferenceCategory>>(
      `/reference-categories/${id}`
    );
    return response.data;
  },

  async getByModuleId(moduleId: number): Promise<ReferenceCategory[]> {
    const response = await apiClient.get<ApiResponse<ReferenceCategory[]>>(
      `/modules/${moduleId}/categories`
    );
    return response.data;
  },

  async create(
    data: ReferenceCategoryRequest
  ): Promise<ReferenceCategory> {
    const response = await apiClient.post<ApiResponse<ReferenceCategory>>(
      "/reference-categories",
      data
    );
    return response.data;
  },

  async update(
    id: number,
    data: ReferenceCategoryRequest
  ): Promise<ReferenceCategory> {
    const response = await apiClient.put<ApiResponse<ReferenceCategory>>(
      `/reference-categories/${id}`,
      data
    );
    return response.data;
  },

  async delete(id: number): Promise<void> {
    await apiClient.delete(`/reference-categories/${id}`);
  },
};
