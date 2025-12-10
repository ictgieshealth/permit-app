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

  async create(data: PermitCreateRequest, file?: File): Promise<Permit> {
    // If file is provided, send as multipart/form-data
    if (file) {
      const formData = new FormData();
      
      // Append all fields
      Object.entries(data).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          formData.append(key, value.toString());
        }
      });
      
      // Append file
      formData.append("file", file);
      
      const response = await apiClient.post<ApiResponse<Permit>>("/permits", formData);
      return response.data;
    }
    
    // Otherwise send as JSON
    const response = await apiClient.post<ApiResponse<Permit>>("/permits", data);
    return response.data;
  },

  async update(id: number, data: PermitUpdateRequest, file?: File): Promise<Permit> {
    // If file is provided, send as multipart/form-data
    if (file) {
      const formData = new FormData();
      
      // Append all fields
      Object.entries(data).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          formData.append(key, value.toString());
        }
      });
      
      // Append file
      formData.append("file", file);
      
      const response = await apiClient.put<ApiResponse<Permit>>(
        `/permits/${id}`,
        formData
      );
      return response.data;
    }
    
    // Otherwise send as JSON
    const response = await apiClient.put<ApiResponse<Permit>>(
      `/permits/${id}`,
      data
    );
    return response.data;
  },

  async delete(id: number): Promise<void> {
    await apiClient.delete(`/permits/${id}`);
  },

  async uploadDocument(id: number, file: File): Promise<Permit> {
    const formData = new FormData();
    formData.append("file", file);

    const response = await apiClient.post<ApiResponse<Permit>>(
      `/permits/${id}/upload`,
      formData
    );
    return response.data;
  },

  async downloadDocument(id: number): Promise<void> {
    const blob = await apiClient.get<Blob>(`/permits/${id}/download`, {
      responseType: "blob",
    });

    // Create blob link to download
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.setAttribute("download", "document");
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
  },

  getPreviewUrl(id: number): string {
    return `${apiClient.getBaseUrl()}/permits/${id}/preview`;
  },

  async getPreviewBlob(id: number): Promise<Blob> {
    const blob = await apiClient.get<Blob>(`/permits/${id}/preview`, {
      responseType: "blob",
    });
    return blob;
  },

  async search(
    query: string,
    params: { page?: number; limit?: number } = {}
  ): Promise<
    ApiResponse<Permit[]> & {
      meta?: { page: number; limit: number; total: number };
    }
  > {
    const queryParams = new URLSearchParams();
    queryParams.append("q", query);
    
    if (params.page) queryParams.append("page", params.page.toString());
    if (params.limit) queryParams.append("limit", params.limit.toString());

    const response = await apiClient.get<ApiResponse<Permit[]>>(
      `/permits/search?${queryParams.toString()}`
    );
    return response as any;
  },
};
