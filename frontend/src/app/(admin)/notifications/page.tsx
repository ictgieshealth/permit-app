"use client";

import React, { useEffect, useState } from "react";
import { useNotifications } from "@/context/NotificationContext";
import notificationService, { Notification } from "@/services/notificationService";
import Link from "next/link";
import { formatDistanceToNow } from "date-fns";
import { id } from "date-fns/locale";

export default function NotificationsPage() {
  const { markAsRead, markAllAsRead, deleteNotification } = useNotifications();
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);

  const fetchNotifications = async (pageNum: number) => {
    try {
      setLoading(true);
      const response = await notificationService.getNotifications(pageNum, 20);
      if (response.data) {
        if (pageNum === 1) {
          setNotifications(response.data);
        } else {
          setNotifications((prev) => [...prev, ...response.data]);
        }
        setHasMore(response.data.length === 20);
      }
    } catch (error) {
      console.error("Error fetching notifications:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNotifications(1);
  }, []);

  const handleMarkAsRead = async (id: number) => {
    await markAsRead([id]);
    fetchNotifications(1);
  };

  const handleMarkAllAsRead = async () => {
    await markAllAsRead();
    fetchNotifications(1);
  };

  const handleDelete = async (id: number) => {
    if (confirm("Apakah Anda yakin ingin menghapus notifikasi ini?")) {
      await deleteNotification(id);
      fetchNotifications(1);
    }
  };

  const loadMore = () => {
    const nextPage = page + 1;
    setPage(nextPage);
    fetchNotifications(nextPage);
  };

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case "expired":
        return (
          <div className="flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
            <svg className="h-6 w-6 text-red-600 dark:text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
        );
      case "expiry_warning":
        return (
          <div className="flex h-12 w-12 items-center justify-center rounded-full bg-orange-100 dark:bg-orange-900/30">
            <svg className="h-6 w-6 text-orange-600 dark:text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        );
      default:
        return (
          <div className="flex h-12 w-12 items-center justify-center rounded-full bg-blue-100 dark:bg-blue-900/30">
            <svg className="h-6 w-6 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        );
    }
  };

  const unreadCount = notifications.filter((n) => !n.is_read).length;

  return (
    <div className="mx-auto max-w-4xl">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-gray-800 dark:text-white">
            Notifikasi
          </h1>
          {unreadCount > 0 && (
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
              {unreadCount} notifikasi belum dibaca
            </p>
          )}
        </div>
        {unreadCount > 0 && (
          <button
            onClick={handleMarkAllAsRead}
            className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-opacity-90"
          >
            Tandai Semua Dibaca
          </button>
        )}
      </div>

      <div className="rounded-lg border border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-dark">
        {loading && page === 1 ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-10 w-10 animate-spin rounded-full border-4 border-solid border-primary border-t-transparent"></div>
          </div>
        ) : notifications.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12">
            <svg className="h-20 w-20 text-gray-400 dark:text-gray-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
            </svg>
            <h3 className="text-lg font-medium text-gray-800 dark:text-white mb-2">
              Tidak Ada Notifikasi
            </h3>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Anda tidak memiliki notifikasi saat ini
            </p>
          </div>
        ) : (
          <>
            <ul className="divide-y divide-gray-200 dark:divide-gray-800">
              {notifications.map((notification) => (
                <li
                  key={notification.id}
                  className={`p-4 hover:bg-gray-50 dark:hover:bg-gray-800/50 ${
                    !notification.is_read ? "bg-blue-50 dark:bg-blue-900/10" : ""
                  }`}
                >
                  <div className="flex gap-4">
                    {getNotificationIcon(notification.type)}
                    <div className="flex-1 min-w-0">
                      <Link
                        href={`/permits/${notification.permit_id}`}
                        className="block"
                      >
                        <h3 className="text-base font-medium text-gray-800 dark:text-white hover:text-primary">
                          {notification.title}
                        </h3>
                        <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                          {notification.message}
                        </p>
                        {notification.permit && (
                          <p className="text-xs text-gray-500 dark:text-gray-500 mt-2">
                            Permit: {notification.permit.name} ({notification.permit.permit_no})
                          </p>
                        )}
                        <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                          {formatDistanceToNow(new Date(notification.created_at), {
                            addSuffix: true,
                            locale: id,
                          })}
                        </p>
                      </Link>
                    </div>
                    <div className="flex flex-col gap-2">
                      {!notification.is_read && (
                        <button
                          onClick={() => handleMarkAsRead(notification.id)}
                          className="text-xs text-primary hover:underline whitespace-nowrap"
                        >
                          Tandai dibaca
                        </button>
                      )}
                      <button
                        onClick={() => handleDelete(notification.id)}
                        className="text-xs text-red-600 hover:underline dark:text-red-400"
                      >
                        Hapus
                      </button>
                    </div>
                  </div>
                </li>
              ))}
            </ul>

            {hasMore && (
              <div className="p-4 text-center border-t border-gray-200 dark:border-gray-800">
                <button
                  onClick={loadMore}
                  disabled={loading}
                  className="rounded-lg bg-gray-200 px-6 py-2 text-sm font-medium text-gray-800 hover:bg-gray-300 disabled:opacity-50 dark:bg-gray-800 dark:text-gray-200 dark:hover:bg-gray-700"
                >
                  {loading ? "Memuat..." : "Muat Lebih Banyak"}
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
