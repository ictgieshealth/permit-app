"use client";
import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Link from "next/link";
import { PermitType } from "@/types/permitType";
import { permitTypeService } from "@/services/permitType.service";

export default function ViewPermitTypePage() {
  const router = useRouter();
  const params = useParams();
  const permitTypeId = parseInt(params.id as string);

  const [permitType, setPermitType] = useState<PermitType | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    loadPermitType();
  }, [permitTypeId]);

  const loadPermitType = async () => {
    try {
      const data = await permitTypeService.getById(permitTypeId);
      setPermitType(data);
    } catch (err: any) {
      setError(err.message || "Failed to load permit type");
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("id-ID", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-500">Loading...</div>
      </div>
    );
  }

  if (error || !permitType) {
    return (
      <div className="p-6">
        <div className="p-4 text-sm border rounded-lg bg-error-50 border-error-200 text-error-700">
          {error || "Permit type not found"}
        </div>
        <Link href="/permit-types" className="inline-block mt-4">
          <Button>Back to Permit Types</Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-4xl p-6 mx-auto">
      <div className="mb-6">
        <div className="flex items-center gap-2 mb-2">
          <Link href="/permit-types" className="text-gray-500 hover:text-gray-700 dark:text-gray-400">
            ← Back
          </Link>
        </div>
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-800 dark:text-white">Permit Type Details</h1>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              View detailed information about this permit type
            </p>
          </div>
          <Link href={`/permit-types/${permitTypeId}/edit`}>
            <Button>Edit Permit Type</Button>
          </Link>
        </div>
      </div>

      <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <div className="space-y-6">
          {/* Basic Information */}
          <div>
            <h3 className="mb-4 text-sm font-medium text-gray-500 uppercase dark:text-gray-400">
              Basic Information
            </h3>
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Permit Type Code
                </label>
                <p className="text-sm text-gray-900 dark:text-white">{permitType.code}</p>
              </div>
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Permit Type Name
                </label>
                <p className="text-sm text-gray-900 dark:text-white">{permitType.name}</p>
              </div>
            </div>
          </div>

          {/* Description */}
          {permitType.description && (
            <div>
              <h3 className="mb-4 text-sm font-medium text-gray-500 uppercase dark:text-gray-400">
                Description
              </h3>
              <p className="text-sm text-gray-900 dark:text-white whitespace-pre-wrap">
                {permitType.description}
              </p>
            </div>
          )}

          {/* Timestamps */}
          <div className="pt-4 border-t border-gray-200 dark:border-gray-800">
            <h3 className="mb-4 text-sm font-medium text-gray-500 uppercase dark:text-gray-400">
              System Information
            </h3>
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Created At
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {formatDate(permitType.created_at)}
                </p>
              </div>
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Last Updated
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {formatDate(permitType.updated_at)}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
