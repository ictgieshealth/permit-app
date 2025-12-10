import { apiClient } from "@/lib/api-client";
import { ApiResponse } from "@/types/api";

export interface Notification {
  id: number;
  user_id: number;
  permit_id: number;
  type: "expiry_reminder" | "expiry_warning" | "expired";
  title: string;
  message: string;
  is_read: boolean;
  read_at: string | null;
  created_at: string;
  permit?: {
    id: number;
    name: string;
    permit_no: string;
    expiry_date: string;
  };
}

export interface NotificationResponse {
  status: string;
  message: string;
  data: Notification[];
}

export interface UnreadCountResponse {
  status: string;
  message: string;
  data: {
    count: number;
  };
}

class NotificationService {
  async getNotifications(page: number = 1, limit: number = 10) {
    return apiClient.get<ApiResponse<Notification[]>>(
      `/notifications?page=${page}&limit=${limit}`
    );
  }

  async getUnreadNotifications() {
    return apiClient.get<ApiResponse<Notification[]>>("/notifications/unread");
  }

  async getUnreadCount() {
    return apiClient.get<ApiResponse<{ count: number }>>("/notifications/unread/count");
  }

  async markAsRead(notificationIds: number[]) {
    return apiClient.post<ApiResponse<{ success: boolean }>>("/notifications/read", {
      notification_ids: notificationIds,
    });
  }

  async markAllAsRead() {
    return apiClient.post<ApiResponse<{ success: boolean }>>("/notifications/read/all");
  }

  async deleteNotification(id: number) {
    return apiClient.delete<ApiResponse<{ success: boolean }>>(`/notifications/${id}`);
  }
}

export default new NotificationService();
