"use client";
import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Link from "next/link";
import { Domain } from "@/types/domain";
import { domainService } from "@/services/domain.service";

export default function ViewDomainPage() {
  const router = useRouter();
  const params = useParams();
  const domainId = parseInt(params.id as string);

  const [domain, setDomain] = useState<Domain | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    loadDomain();
  }, [domainId]);

  const loadDomain = async () => {
    try {
      const data = await domainService.getById(domainId);
      setDomain(data);
    } catch (err: any) {
      setError(err.message || "Failed to load domain");
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

  if (error || !domain) {
    return (
      <div className="p-6">
        <div className="p-4 text-sm border rounded-lg bg-error-50 border-error-200 text-error-700">
          {error || "Domain not found"}
        </div>
        <Link href="/domains" className="inline-block mt-4">
          <Button>Back to Domains</Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-4xl p-6 mx-auto">
      <div className="mb-6">
        <div className="flex items-center gap-2 mb-2">
          <Link href="/domains" className="text-gray-500 hover:text-gray-700 dark:text-gray-400">
            ← Back
          </Link>
        </div>
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-800 dark:text-white">Domain Details</h1>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              View detailed information about this domain
            </p>
          </div>
          <Link href={`/domains/${domainId}/edit`}>
            <Button>Edit Domain</Button>
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
                  Domain Code
                </label>
                <p className="text-sm text-gray-900 dark:text-white">{domain.code}</p>
              </div>
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Domain Name
                </label>
                <p className="text-sm text-gray-900 dark:text-white">{domain.name}</p>
              </div>
            </div>
          </div>

          {/* Description */}
          {domain.description && (
            <div>
              <h3 className="mb-4 text-sm font-medium text-gray-500 uppercase dark:text-gray-400">
                Description
              </h3>
              <p className="text-sm text-gray-900 dark:text-white whitespace-pre-wrap">
                {domain.description}
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
                  {formatDate(domain.created_at)}
                </p>
              </div>
              <div>
                <label className="block mb-1 text-sm font-medium text-gray-700 dark:text-gray-300">
                  Last Updated
                </label>
                <p className="text-sm text-gray-900 dark:text-white">
                  {formatDate(domain.updated_at)}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
