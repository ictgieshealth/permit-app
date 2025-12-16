import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import {
  Project,
  ProjectCreateRequest,
  ProjectUpdateRequest,
  ProjectListRequest,
  ProjectStatusChangeRequest,
} from "@/types/project";

export const projectService = {
  async getAll(
    params: ProjectListRequest = {}
  ): Promise<
    ApiResponse<Project[]> & {
      meta?: { page: number; limit: number; total: number };
    }
  > {
    const queryParams = new URLSearchParams();

    if (params.domain_id)
      queryParams.append("domain_id", params.domain_id.toString());
    if (params.project_status_id)
      queryParams.append("project_status_id", params.project_status_id.toString());
    if (params.name) queryParams.append("name", params.name);
    if (params.code) queryParams.append("code", params.code);
    if (params.status !== undefined)
      queryParams.append("status", params.status.toString());
    if (params.page) queryParams.append("page", params.page.toString());
    if (params.limit) queryParams.append("limit", params.limit.toString());

    const response = await apiClient.get<ApiResponse<Project[]>>(
      `/projects?${queryParams.toString()}`
    );
    return response as any;
  },

  async getById(id: number): Promise<Project> {
    const response = await apiClient.get<ApiResponse<Project>>(
      `/projects/${id}`
    );
    return response.data;
  },

  async getByDomainId(domainId: number): Promise<Project[]> {
    const response = await apiClient.get<ApiResponse<Project[]>>(
      `/domains/${domainId}/projects`
    );
    return response.data;
  },

  async getByUserId(userId: number): Promise<Project[]> {
    const response = await apiClient.get<ApiResponse<Project[]>>(
      `/users/${userId}/projects`
    );
    return response.data;
  },

  async create(data: ProjectCreateRequest): Promise<Project> {
    const response = await apiClient.post<ApiResponse<Project>>(
      "/projects",
      data
    );
    return response.data;
  },

  async update(id: number, data: ProjectUpdateRequest): Promise<Project> {
    const response = await apiClient.put<ApiResponse<Project>>(
      `/projects/${id}`,
      data
    );
    return response.data;
  },

  async delete(id: number): Promise<void> {
    await apiClient.delete(`/projects/${id}`);
  },

  async changeStatus(
    id: number,
    data: ProjectStatusChangeRequest
  ): Promise<Project> {
    const response = await apiClient.post<ApiResponse<Project>>(
      `/projects/${id}/change-status`,
      data
    );
    return response.data;
  },
};
