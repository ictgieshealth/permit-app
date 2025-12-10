"use client";
import { useEffect, useState } from "react";
import { Permit } from "@/types/permit";
import { permitService } from "@/services/permit.service";

export default function PermitStatusChart() {
  const [permits, setPermits] = useState<Permit[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPermits();
  }, []);

  const loadPermits = async () => {
    try {
      const response = await permitService.getAll({ limit: 1000 });
      setPermits(response.data);
    } catch (err) {
      console.error("Failed to load permits:", err);
    } finally {
      setLoading(false);
    }
  };

  const now = new Date();
  const statusCounts = {
    active: permits.filter((p) => p.status === "active").length,
    // Count expired based on expiry_date, not status field
    expired: permits.filter((p) => {
      const expiryDate = new Date(p.expiry_date);
      return expiryDate < now && p.status !== "inactive";
    }).length,
    inactive: permits.filter((p) => p.status === "inactive").length,
  };

  const total = permits.length || 1;
  const statusData = [
    {
      status: "Active",
      count: statusCounts.active,
      percentage: ((statusCounts.active / total) * 100).toFixed(1),
      color: "bg-success-500",
      lightColor: "bg-success-100 dark:bg-success-900/20",
      textColor: "text-success-700 dark:text-success-400",
    },
    {
      status: "Expired",
      count: statusCounts.expired,
      percentage: ((statusCounts.expired / total) * 100).toFixed(1),
      color: "bg-error-500",
      lightColor: "bg-error-100 dark:bg-error-900/20",
      textColor: "text-error-700 dark:text-error-400",
    },
    {
      status: "Inactive",
      count: statusCounts.inactive,
      percentage: ((statusCounts.inactive / total) * 100).toFixed(1),
      color: "bg-gray-500",
      lightColor: "bg-gray-100 dark:bg-gray-800",
      textColor: "text-gray-700 dark:text-gray-400",
    },
  ];

  if (loading) {
    return (
      <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <h3 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">
          Permit Status Distribution
        </h3>
        <div className="space-y-4 animate-pulse">
          {[...Array(3)].map((_, i) => (
            <div key={i}>
              <div className="h-4 mb-2 bg-gray-200 rounded dark:bg-gray-700"></div>
              <div className="h-8 bg-gray-200 rounded dark:bg-gray-700"></div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
      <h3 className="mb-6 text-lg font-semibold text-gray-800 dark:text-white">
        Permit Status Distribution
      </h3>

      <div className="space-y-4">
        {statusData.map((item, index) => (
          <div key={index}>
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                {item.status}
              </span>
              <span className="text-sm font-semibold text-gray-800 dark:text-white">
                {item.count} ({item.percentage}%)
              </span>
            </div>
            <div className="w-full h-3 overflow-hidden bg-gray-200 rounded-full dark:bg-gray-700">
              <div
                className={`h-full ${item.color} rounded-full transition-all duration-500`}
                style={{ width: `${item.percentage}%` }}
              ></div>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-3 gap-3 pt-6 mt-6 border-t border-gray-200 dark:border-gray-800">
        {statusData.map((item, index) => (
          <div
            key={index}
            className={`p-3 rounded-lg ${item.lightColor} text-center`}
          >
            <p className={`text-2xl font-bold ${item.textColor}`}>
              {item.count}
            </p>
            <p className="text-xs text-gray-600 dark:text-gray-400">
              {item.status}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
