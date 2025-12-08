import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import { Menu } from "@/types/menu";

export const menuService = {
  async getAll(params?: {
    name?: string;
    path?: string;
    is_active?: boolean;
    page?: number;
    limit?: number;
  }): Promise<{ data: Menu[]; meta?: { page: number; limit: number; total: number } }> {
    const queryParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          queryParams.append(key, String(value));
        }
      });
    }
    
    const endpoint = `/menus?${queryParams.toString()}`;
    return await apiClient.get<ApiResponse<Menu[]>>(endpoint) as any;
  },

  async getById(id: number) {
    const response = await apiClient.get<ApiResponse<Menu>>(`/menus/${id}`);
    return response.data;
  },

  async getUserMenus() {
    const response = await apiClient.get<ApiResponse<Menu[]>>('/menus/user');
    return response.data;
  },

  async create(data: {
    name: string;
    path: string;
    icon?: string;
    parent_id?: number | null;
    order_index: number;
    role_ids?: number[];
  }) {
    const response = await apiClient.post<ApiResponse<Menu>>('/menus', data);
    return response.data;
  },

  async update(id: number, data: {
    name: string;
    path: string;
    icon?: string;
    parent_id?: number | null;
    order_index: number;
    role_ids?: number[];
  }) {
    const response = await apiClient.put<ApiResponse<Menu>>(`/menus/${id}`, data);
    return response.data;
  },

  async delete(id: number) {
    const response = await apiClient.delete<ApiResponse<void>>(`/menus/${id}`);
    return response.data;
  },

  async assignRoles(menuId: number, roleIds: number[]) {
    const response = await apiClient.post<ApiResponse<void>>(`/menus/${menuId}/roles`, {
      role_ids: roleIds,
    });
    return response.data;
  },
};
