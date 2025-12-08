"use client";
import { useEffect, useState } from "react";
import { permitService } from "@/services/permit.service";
import { divisionService } from "@/services/division.service";
import { permitTypeService } from "@/services/permitType.service";
import { domainService } from "@/services/domain.service";

export default function DashboardMetrics() {
  const [stats, setStats] = useState({
    totalPermits: 0,
    activePermits: 0,
    expiredPermits: 0,
    totalDivisions: 0,
    totalPermitTypes: 0,
    totalDomains: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      const [permits, divisions, permitTypes, domains] = await Promise.all([
        permitService.getAll({ limit: 1000 }),
        divisionService.getAll({ limit: 1000 }),
        permitTypeService.getAll({ limit: 1000 }),
        domainService.getAll({ limit: 1000 }),
      ]);

      const activePermits = permits.data.filter((p) => p.status === "active").length;
      const expiredPermits = permits.data.filter((p) => p.status === "expired").length;

      setStats({
        totalPermits: permits.data.length,
        activePermits,
        expiredPermits,
        totalDivisions: divisions.data.length,
        totalPermitTypes: permitTypes.data.length,
        totalDomains: domains.data.length,
      });
    } catch (err) {
      console.error("Failed to load stats:", err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        {[...Array(6)].map((_, i) => (
          <div
            key={i}
            className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800 animate-pulse"
          >
            <div className="h-4 bg-gray-200 rounded dark:bg-gray-700"></div>
            <div className="w-20 h-8 mt-4 bg-gray-200 rounded dark:bg-gray-700"></div>
          </div>
        ))}
      </div>
    );
  }

  const metrics = [
    {
      title: "Total Permits",
      value: stats.totalPermits,
      icon: "📄",
      color: "brand",
    },
    {
      title: "Active Permits",
      value: stats.activePermits,
      icon: "✅",
      color: "success",
    },
    {
      title: "Expired Permits",
      value: stats.expiredPermits,
      icon: "⚠️",
      color: "error",
    },
    {
      title: "Total Divisions",
      value: stats.totalDivisions,
      icon: "🏢",
      color: "blue",
    },
    {
      title: "Permit Types",
      value: stats.totalPermitTypes,
      icon: "📋",
      color: "purple",
    },
    {
      title: "Total Domains",
      value: stats.totalDomains,
      icon: "🌐",
      color: "indigo",
    },
  ];

  const colorClasses = {
    brand: "text-brand-500 bg-brand-100 dark:bg-brand-900/20",
    success: "text-success-500 bg-success-100 dark:bg-success-900/20",
    error: "text-error-500 bg-error-100 dark:bg-error-900/20",
    blue: "text-blue-500 bg-blue-100 dark:bg-blue-900/20",
    purple: "text-purple-500 bg-purple-100 dark:bg-purple-900/20",
    indigo: "text-indigo-500 bg-indigo-100 dark:bg-indigo-900/20",
  };

  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
      {metrics.map((metric, index) => (
        <div
          key={index}
          className="p-6 transition-shadow bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800 hover:shadow-lg"
        >
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400">
                {metric.title}
              </p>
              <p className="mt-2 text-3xl font-bold text-gray-800 dark:text-white">
                {metric.value}
              </p>
            </div>
            <div
              className={`flex items-center justify-center w-12 h-12 rounded-lg ${
                colorClasses[metric.color as keyof typeof colorClasses]
              }`}
            >
              <span className="text-2xl">{metric.icon}</span>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
