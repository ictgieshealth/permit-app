"use client";
import { useEffect, useState } from "react";
import { Permit } from "@/types/permit";
import { permitService } from "@/services/permit.service";
import Link from "next/link";

export default function RecentPermits() {
  const [permits, setPermits] = useState<Permit[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadRecentPermits();
  }, []);

  const loadRecentPermits = async () => {
    try {
      const response = await permitService.getAll({ limit: 5, page: 1 });
      setPermits(response.data);
    } catch (err) {
      console.error("Failed to load recent permits:", err);
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const getStatusBadge = (status: string) => {
    const statusClasses = {
      active:
        "bg-success-100 text-success-700 dark:bg-success-900/20 dark:text-success-400",
      expired:
        "bg-error-100 text-error-700 dark:bg-error-900/20 dark:text-error-400",
      inactive:
        "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400",
    };

    return (
      <span
        className={`px-2 py-1 text-xs font-medium rounded-full ${
          statusClasses[status as keyof typeof statusClasses] ||
          statusClasses.inactive
        }`}
      >
        {status}
      </span>
    );
  };

  if (loading) {
    return (
      <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <h3 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">
          Recent Permits
        </h3>
        <div className="space-y-3">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="p-3 bg-gray-100 rounded-lg animate-pulse dark:bg-gray-800">
              <div className="h-4 bg-gray-200 rounded dark:bg-gray-700"></div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-gray-800 dark:text-white">
          Recent Permits
        </h3>
        <Link
          href="/permits"
          className="text-sm text-brand-500 hover:text-brand-600"
        >
          View All →
        </Link>
      </div>

      <div className="space-y-3">
        {permits.length === 0 ? (
          <p className="text-sm text-center text-gray-500 dark:text-gray-400">
            No permits found
          </p>
        ) : (
          permits.map((permit) => (
            <Link
              key={permit.id}
              href={`/permits/${permit.id}`}
              className="block p-3 transition-colors border border-gray-200 rounded-lg hover:bg-gray-50 dark:border-gray-700 dark:hover:bg-gray-800/50"
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <p className="font-medium text-gray-800 dark:text-white">
                    {permit.permit_no}
                  </p>
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    {permit.permit_type?.name || "-"}
                  </p>
                  <p className="text-xs text-gray-400 dark:text-gray-500">
                    Expires: {formatDate(permit.expiry_date)}
                  </p>
                </div>
                <div>{getStatusBadge(permit.status)}</div>
              </div>
            </Link>
          ))
        )}
      </div>
    </div>
  );
}
