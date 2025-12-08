"use client";
import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Link from "next/link";
import { Permit } from "@/types/permit";
import { permitService } from "@/services/permit.service";

export default function ViewPermitPage() {
  const router = useRouter();
  const params = useParams();
  const permitId = parseInt(params.id as string);

  const [permit, setPermit] = useState<Permit | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    loadPermit();
  }, [permitId]);

  const loadPermit = async () => {
    try {
      const data = await permitService.getById(permitId);
      setPermit(data);
    } catch (err: any) {
      setError(err.message || "Failed to load permit");
    } finally {
      setLoading(false);
    }
  };

  const handleDownloadDocument = async () => {
    try {
      await permitService.downloadDocument(permitId);
    } catch (err: any) {
      setError(err.message || "Failed to download document");
    }
  };

  const formatFileSize = (bytes?: number | null) => {
    if (!bytes) return "-";
    if (bytes < 1024) return bytes + " B";
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + " KB";
    return (bytes / (1024 * 1024)).toFixed(2) + " MB";
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("id-ID", {
      year: "numeric",
      month: "long",
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
        className={`px-3 py-1 text-sm font-medium rounded-full ${
          statusClasses[status as keyof typeof statusClasses] ||
          statusClasses.inactive
        }`}
      >
        {status.toUpperCase()}
      </span>
    );
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-500">Loading...</div>
      </div>
    );
  }

  if (error || !permit) {
    return (
      <div className="p-6">
        <div className="p-4 text-sm border rounded-lg bg-error-50 border-error-200 text-error-700">
          {error || "Permit not found"}
        </div>
        <Link href="/permits" className="inline-block mt-4">
          <Button>Back to Permits</Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-4xl p-6 mx-auto">
      <div className="mb-6">
        <div className="flex items-center gap-2 mb-2">
          <Link
            href="/permits"
            className="text-gray-500 hover:text-gray-700 dark:text-gray-400"
          >
            ← Back
          </Link>
        </div>
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
              Permit Details
            </h1>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              View detailed information about this permit
            </p>
          </div>
          <Link href={`/permits/${permitId}/edit`}>
            <Button>Edit Permit</Button>
          </Link>
        </div>
      </div>

      <div className="p-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <div className="space-y-6">
          {/* Status */}
          <div className="flex items-center justify-between pb-4 border-b border-gray-200 dark:border-gray-800">
            <h2 className="text-lg font-semibold text-gray-800 dark:text-white">
              Status
            </h2>
            {getStatusBadge(permit.status)}
          </div>

          {/* Basic Information */}
          <div>
            <h3 className="mb-4 text-sm font-medium text-gray-500 uppercase dark:text-gray-400">
              Basic Information
            </h3>
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Permit Number
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {permit.permit_no}
                </p>
              </div>
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Application Type
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {permit.application_type}
                </p>
              </div>
            </div>
          </div>

          {/* Domain and Permit Type */}
          <div>
            <h3 className="mb-4 text-sm font-medium text-gray-500 uppercase dark:text-gray-400">
              Domain, Division & Category
            </h3>
            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Domain
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {permit.domain?.name || "-"}
                </p>
                <p className="text-xs text-gray-500 dark:text-gray-400">
                  {permit.domain?.code || "-"}
                </p>
              </div>
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Division
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {permit.division?.name || "-"}
                </p>
                <p className="text-xs text-gray-500 dark:text-gray-400">
                  {permit.division?.code || "-"}
                </p>
              </div>
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Category
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {permit.permit_type?.name || "-"}
                </p>
              </div>
            </div>
          </div>

          {/* Equipment/Services */}
          <div>
            <h3 className="mb-4 text-sm font-medium text-gray-500 uppercase dark:text-gray-400">
              Permit Name
            </h3>
            <p className="text-sm text-gray-900 dark:text-white">
              {permit.name}
            </p>
          </div>

          {/* Dates */}
          <div>
            <h3 className="mb-4 text-sm font-medium text-gray-500 uppercase dark:text-gray-400">
              Validity Period
            </h3>
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Expiry Date
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {formatDate(permit.expiry_date)}
                </p>
              </div>
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Effective Term
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {permit.effective_term || "-"}
                </p>
              </div>
            </div>
          </div>

          {/* Responsible Person */}
          <div>
            <h3 className="mb-4 text-sm font-medium text-gray-500 uppercase dark:text-gray-400">
              Responsible Person
            </h3>
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Permit Responsible
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {permit.responsible_person?.full_name || "-"}
                </p>
                {permit.responsible_person && (
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    {permit.responsible_person.email}
                  </p>
                )}
              </div>
              
            </div>
          </div>

          {/* Document Information */}
          <div>
            <h3 className="mb-4 text-sm font-medium text-gray-500 uppercase dark:text-gray-400">
              Document Information
            </h3>
            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Document Name
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {permit.doc_name || "-"}
                </p>
              </div>
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Document Number
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {permit.doc_number || "-"}
                </p>
              </div>
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Document Responsible
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {permit.responsible_doc_person?.full_name || "-"}
                </p>
                {permit.responsible_doc_person && (
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    {permit.responsible_doc_person.email}
                  </p>
                )}
              </div>
            </div>

            {/* Uploaded File Information */}
            {permit.doc_file_name && (
              <div className="p-4 mt-4 border border-gray-200 rounded-lg bg-gray-50 dark:bg-gray-800 dark:border-gray-700">
                <h4 className="mb-3 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Uploaded Document
                </h4>
                <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
                  <div>
                    <label className="block mb-1 text-xs font-medium text-gray-600 dark:text-gray-400">
                      File Name
                    </label>
                    <p className="text-sm text-gray-900 break-all dark:text-white">
                      {permit.doc_file_name}
                    </p>
                  </div>
                  <div>
                    <label className="block mb-1 text-xs font-medium text-gray-600 dark:text-gray-400">
                      File Size
                    </label>
                    <p className="text-sm text-gray-900 dark:text-white">
                      {formatFileSize(permit.doc_file_size)}
                    </p>
                  </div>
                  <div>
                    <label className="block mb-1 text-xs font-medium text-gray-600 dark:text-gray-400">
                      File Type
                    </label>
                    <p className="text-sm text-gray-900 dark:text-white">
                      {permit.doc_file_type || "-"}
                    </p>
                  </div>
                </div>
                <div className="mt-3">
                  <Button
                    type="button"
                    onClick={handleDownloadDocument}
                    className="text-sm bg-blue-600 hover:bg-blue-700"
                  >
                    📥 Download Document
                  </Button>
                </div>
              </div>
            )}
          </div>

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
                  {formatDate(permit.created_at)}
                </p>
              </div>
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Last Updated
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {formatDate(permit.updated_at)}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
