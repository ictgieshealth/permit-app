import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import {
  Permit,
  PermitCreateRequest,
  PermitUpdateRequest,
  PermitListRequest,
} from "@/types/permit";

export const permitService = {
  async getAll(
    params: PermitListRequest = {}
  ): Promise<
    ApiResponse<Permit[]> & {
      meta?: { page: number; limit: number; total: number };
    }
  > {
    const queryParams = new URLSearchParams();

    if (params.domain_id)
      queryParams.append("domain_id", params.domain_id.toString());
    if (params.permit_type_id)
      queryParams.append("permit_type_id", params.permit_type_id.toString());
    if (params.application_type)
      queryParams.append("application_type", params.application_type);
    if (params.permit_no) queryParams.append("permit_no", params.permit_no);
    if (params.responsible_person)
      queryParams.append("responsible_person", params.responsible_person);
    if (params.status) queryParams.append("status", params.status);
    if (params.page) queryParams.append("page", params.page.toString());
    if (params.limit) queryParams.append("limit", params.limit.toString());

    const response = await apiClient.get<ApiResponse<Permit[]>>(
      `/permits?${queryParams.toString()}`
    );
    return response as any;
  },

  async getById(id: number): Promise<Permit> {
    const response = await apiClient.get<ApiResponse<Permit>>(`/permits/${id}`);
    return response.data;
  },

  async create(data: PermitCreateRequest): Promise<Permit> {
    const response = await apiClient.post<ApiResponse<Permit>>("/permits", data);
    return response.data;
  },

  async update(id: number, data: PermitUpdateRequest): Promise<Permit> {
    const response = await apiClient.put<ApiResponse<Permit>>(
      `/permits/${id}`,
      data
    );
    return response.data;
  },

  async delete(id: number): Promise<void> {
    await apiClient.delete(`/permits/${id}`);
  },
};
