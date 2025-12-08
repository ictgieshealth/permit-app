import type { Metadata } from "next";
import React from "react";
import DashboardMetrics from "@/components/dashboard/DashboardMetrics";
import RecentPermits from "@/components/dashboard/RecentPermits";
import PermitStatusChart from "@/components/dashboard/PermitStatusChart";

export const metadata: Metadata = {
  title: "Dashboard | Permit Management System",
  description: "Permit Management System Dashboard",
};

export default function Dashboard() {
  return (
    <div className="space-y-6">
      {/* Metrics Overview */}
      <div>
        <h2 className="mb-4 text-xl font-bold text-gray-800 dark:text-white">
          Overview
        </h2>
        <DashboardMetrics />
      </div>

      {/* Charts and Recent Data */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <PermitStatusChart />
        <RecentPermits />
      </div>
    </div>
  );
}
