"use client";
import React from "react";
import { useAuth } from "@/context/AuthContext";
import Image from "next/image";


export default function UserMetaCard() {
  const { user } = useAuth();

  if (!user) {
    return (
      <div className="p-5 border border-gray-200 rounded-2xl dark:border-gray-800 lg:p-6">
        <p className="text-gray-500">Loading user information...</p>
      </div>
    );
  }

  return (
    <>
      <div className="p-5 border border-gray-200 rounded-2xl dark:border-gray-800 lg:p-6">
        <div className="flex flex-col gap-5 xl:flex-row xl:items-center xl:justify-between">
          <div className="flex flex-col items-center w-full gap-6 xl:flex-row">
            <div className="w-20 h-20 overflow-hidden border border-gray-200 rounded-full dark:border-gray-800 bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
              {user.full_name ? (
                <span className="text-2xl font-semibold text-gray-600 dark:text-gray-300">
                  {user.full_name.charAt(0).toUpperCase()}
                </span>
              ) : (
                <span className="text-2xl font-semibold text-gray-600 dark:text-gray-300">
                  {user.username.charAt(0).toUpperCase()}
                </span>
              )}
            </div>
            <div className="order-3 xl:order-2">
              <h4 className="mb-2 text-lg font-semibold text-center text-gray-800 dark:text-white/90 xl:text-left">
                {user.full_name || user.username}
              </h4>
              <div className="flex flex-col items-center gap-1 text-center xl:flex-row xl:gap-3 xl:text-left">
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  {user.role?.name || "User"}
                </p>
                {user.email && (
                  <>
                    <div className="hidden h-3.5 w-px bg-gray-300 dark:bg-gray-700 xl:block"></div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {user.email}
                    </p>
                  </>
                )}
              </div>
            </div>
            <div className="flex items-center order-2 gap-2 grow xl:order-3 xl:justify-end">
              {user.domains && user.domains.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {user.domains.map((domain) => (
                    <span
                      key={domain.id}
                      className="inline-block px-3 py-1.5 text-xs font-medium rounded-full bg-brand-100 text-brand-700 dark:bg-brand-900/20 dark:text-brand-400"
                    >
                      {domain.name}
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
