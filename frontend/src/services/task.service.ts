import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";
import {
  Task,
  TaskCreateRequest,
  TaskUpdateRequest,
  TaskListRequest,
  TaskChangeStatusRequest,
  TaskChangeTypeRequest,
  TaskInReviewRequest,
  TaskReasonRequest,
  TaskRevisionRequest,
  ApprovalRequest,
} from "@/types/task";

export const taskService = {
  // Get all approved tasks (approval_status_id = 22)
  // Matches PHP TaskController::index() - only shows approved tasks
  async getAll(
    params: TaskListRequest = {}
  ): Promise<
    ApiResponse<Task[]> & {
      meta?: { page: number; limit: number; total: number };
    }
  > {
    const queryParams = new URLSearchParams();

    if (params.search) queryParams.append("search", params.search);
    if (params.project_id)
      queryParams.append("project_id", params.project_id.toString());
    if (params.status_id)
      queryParams.append("status_id", params.status_id.toString());
    if (params.approval_status_id)
      queryParams.append("approval_status_id", params.approval_status_id.toString());
    if (params.assigned_id)
      queryParams.append("assigned_id", params.assigned_id.toString());
    if (params.start_date) queryParams.append("start_date", params.start_date);
    if (params.end_date) queryParams.append("end_date", params.end_date);
    if (params.page) queryParams.append("page", params.page.toString());
    if (params.limit) queryParams.append("limit", params.limit.toString());

    const response = await apiClient.get<ApiResponse<Task[]>>(
      `/tasks?${queryParams.toString()}`
    );
    return response as any;
  },

  // Get all task approval requests (can filter by approval_status_id)
  // Matches PHP TaskRequestController::index() - shows pending/rejected/approved tasks
  async getAllRequests(
    params: TaskListRequest = {}
  ): Promise<
    ApiResponse<Task[]> & {
      meta?: { page: number; limit: number; total: number };
    }
  > {
    const queryParams = new URLSearchParams();

    if (params.search) queryParams.append("search", params.search);
    if (params.project_id)
      queryParams.append("project_id", params.project_id.toString());
    if (params.status_id)
      queryParams.append("status_id", params.status_id.toString());
    if (params.approval_status_id)
      queryParams.append("approval_status_id", params.approval_status_id.toString());
    if (params.assigned_id)
      queryParams.append("assigned_id", params.assigned_id.toString());
    if (params.start_date) queryParams.append("start_date", params.start_date);
    if (params.end_date) queryParams.append("end_date", params.end_date);
    if (params.page) queryParams.append("page", params.page.toString());
    if (params.limit) queryParams.append("limit", params.limit.toString());

    const response = await apiClient.get<ApiResponse<Task[]>>(
      `/task-requests?${queryParams.toString()}`
    );
    return response as any;
  },

  async getById(id: number): Promise<Task> {
    const response = await apiClient.get<ApiResponse<Task>>(`/tasks/${id}`);
    return response.data;
  },

  async getByCode(code: string): Promise<Task> {
    const response = await apiClient.get<ApiResponse<Task>>(`/tasks/code/${code}`);
    return response.data;
  },

  async create(data: TaskCreateRequest, files?: File[]): Promise<Task> {
    const formData = new FormData();

    // Append all fields
    Object.entries(data).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        formData.append(key, value.toString());
      }
    });

    // Append files
    if (files && files.length > 0) {
      files.forEach((file) => {
        formData.append("files", file);
      });
    }

    const response = await apiClient.post<ApiResponse<Task>>("/tasks", formData);
    return response.data;
  },

  async update(id: number, data: TaskUpdateRequest, files?: File[], deletedFileIds?: number[]): Promise<Task> {
    const formData = new FormData();

    // Append all fields
    Object.entries(data).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        formData.append(key, value.toString());
      }
    });

    // Append files
    if (files && files.length > 0) {
      files.forEach((file) => {
        formData.append("files", file);
      });
    }

    // Append deleted file IDs
    if (deletedFileIds && deletedFileIds.length > 0) {
      deletedFileIds.forEach((id) => {
        formData.append("deleted_file_ids[]", id.toString());
      });
    }

    const response = await apiClient.put<ApiResponse<Task>>(
      `/tasks/${id}`,
      formData
    );
    return response.data;
  },

  async delete(id: number): Promise<void> {
    await apiClient.delete(`/tasks/${id}`);
  },

  async changeStatus(id: number, data: TaskChangeStatusRequest): Promise<void> {
    await apiClient.post(`/tasks/${id}/change-status`, data);
  },

  async changeType(id: number, data: TaskChangeTypeRequest): Promise<void> {
    await apiClient.post(`/tasks/${id}/change-type`, data);
  },

  async inReview(id: number, data: TaskInReviewRequest, files?: File[]): Promise<void> {
    const formData = new FormData();

    formData.append("description_before", data.description_before);
    formData.append("description_after", data.description_after);

    // Append files
    if (files && files.length > 0) {
      files.forEach((file) => {
        formData.append("files", file);
      });
    }

    await apiClient.post(`/tasks/${id}/in-review`, formData);
  },

  async setReason(id: number, data: TaskReasonRequest): Promise<void> {
    await apiClient.post(`/tasks/${id}/set-reason`, data);
  },

  async setRevision(id: number, data: TaskRevisionRequest, files?: File[]): Promise<void> {
    const formData = new FormData();

    formData.append("revision", data.revision);

    // Append files
    if (files && files.length > 0) {
      files.forEach((file) => {
        formData.append("files", file);
      });
    }

    await apiClient.post(`/tasks/${id}/set-revision`, formData);
  },

  async approveTask(
    taskId: number,
    approvalTaskId: number,
    data: ApprovalRequest
  ): Promise<void> {
    await apiClient.post(
      `/tasks/${taskId}/approvals/${approvalTaskId}/approve`,
      data
    );
  },

  async rejectTask(
    taskId: number,
    approvalTaskId: number,
    data: ApprovalRequest
  ): Promise<void> {
    await apiClient.post(
      `/tasks/${taskId}/approvals/${approvalTaskId}/reject`,
      data
    );
  },
};
