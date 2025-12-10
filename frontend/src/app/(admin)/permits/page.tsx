"use client";
import { useState, useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Link from "next/link";
import { Permit } from "@/types/permit";
import { permitService } from "@/services/permit.service";
import { authService } from "@/services/auth.service";
import { ConfirmModal } from "@/components/ui/modal/ConfirmModal";

export default function PermitsPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const searchQuery = searchParams.get("search");
  
  const [permits, setPermits] = useState<Permit[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [total, setTotal] = useState(0);
  const [filters, setFilters] = useState({
    permit_no: "",
    equipment_name: "",
    status: "",
    application_type: "",
  });
  const [deleteModal, setDeleteModal] = useState<{
    isOpen: boolean;
    permitId: number | null;
    permitNo: string;
  }>({ isOpen: false, permitId: null, permitNo: "" });
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    loadPermits();
  }, [page, filters, searchQuery]);

  const loadPermits = async () => {
    try {
      setLoading(true);
      
      // If there's a search query, use search endpoint
      if (searchQuery) {
        const response = await permitService.search(searchQuery, {
          page,
          limit,
        });
        setPermits(response.data);
        setTotal(response.meta?.total || 0);
      } else {
        // Otherwise use normal getAll
        const response = await permitService.getAll({
          page,
          limit,
          ...filters,
        });
        setPermits(response.data);
        setTotal(response.meta?.total || 0);
      }
    } catch (err: any) {
      setError(err.message || "Failed to load permits");
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = (id: number, permitNo: string) => {
    setDeleteModal({ isOpen: true, permitId: id, permitNo: permitNo });
  };

  const confirmDelete = async () => {
    if (!deleteModal.permitId) return;

    setDeleting(true);
    try {
      await permitService.delete(deleteModal.permitId);
      setDeleteModal({ isOpen: false, permitId: null, permitNo: "" });
      loadPermits();
    } catch (err: any) {
      alert(err.message || "Failed to delete permit");
    } finally {
      setDeleting(false);
    }
  };

  const handleFilterChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    setFilters({
      ...filters,
      [e.target.name]: e.target.value,
    });
    setPage(1); // Reset to first page when filtering
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const getStatusBadge = (permit: Permit) => {
    const now = new Date();
    const expiryDate = new Date(permit.expiry_date);
    
    // Check if permit is expired based on date
    const isExpired = expiryDate < now && permit.status !== "inactive";
    const displayStatus = isExpired ? "expired" : permit.status;
    
    const statusClasses = {
      active: "bg-success-100 text-success-700 dark:bg-success-900/20 dark:text-success-400",
      expired: "bg-error-100 text-error-700 dark:bg-error-900/20 dark:text-error-400",
      inactive: "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400",
    };

    return (
      <span
        className={`px-2 py-1 text-xs font-medium rounded-full ${
          statusClasses[displayStatus as keyof typeof statusClasses] || statusClasses.inactive
        }`}
      >
        {displayStatus}
      </span>
    );
  };

  const totalPages = Math.ceil(total / limit);

  if (loading && permits.length === 0) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-500">Loading...</div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
            Permit Management
          </h1>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {searchQuery 
              ? `Search results for "${searchQuery}"`
              : "Manage all permits and their details"
            }
          </p>
        </div>
        <div className="flex items-center gap-3">
          {searchQuery && (
            <button
              onClick={() => router.push("/permits")}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 dark:bg-gray-800 dark:text-gray-300 dark:border-gray-700 dark:hover:bg-gray-700"
            >
              Clear Search
            </button>
          )}
          <Link href="/permits/create">
            <Button>+ Create Permit</Button>
          </Link>
        </div>
      </div>

      {/* Filters - Only show when not searching */}
      {!searchQuery && (
        <div className="p-4 mb-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
          <div className="grid grid-cols-1 gap-4 md:grid-cols-4">
            <input
              type="text"
              name="permit_no"
              placeholder="Search by Permit No..."
              value={filters.permit_no}
              onChange={handleFilterChange}
              className="px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
            />
            <input
              type="text"
              name="equipment_name"
              placeholder="Search by Equipment/Services..."
              value={filters.equipment_name}
              onChange={handleFilterChange}
              className="px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
            />
            <input
              type="text"
              name="application_type"
              placeholder="Search by Application Type..."
              value={filters.application_type}
              onChange={handleFilterChange}
              className="px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
            />
            <select
              name="status"
              value={filters.status}
              onChange={handleFilterChange}
              className="px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
            >
              <option value="">All Status</option>
              <option value="active">Active</option>
              <option value="expired">Expired</option>
              <option value="inactive">Inactive</option>
            </select>
          </div>
        </div>
      )}

      {error && (
        <div className="p-3 mb-4 text-sm border rounded-lg bg-error-50 border-error-200 text-error-700 dark:bg-error-900/20">
          {error}
        </div>
      )}

      {/* Table */}
      <div className="overflow-hidden bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th className="px-4 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Permit No
                </th>
                <th className="px-4 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Category
                </th>
                <th className="px-4 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Equipment/Services
                </th>
                <th className="px-4 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Domain
                </th>
                <th className="px-4 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Division
                </th>
                <th className="px-4 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Application Type
                </th>
                <th className="px-4 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Expiry Date
                </th>
                <th className="px-4 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Status
                </th>
                <th className="px-4 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200 dark:divide-gray-800">
              {permits.length === 0 ? (
                <tr>
                  <td
                    colSpan={9}
                    className="px-4 py-8 text-sm text-center text-gray-500"
                  >
                    No permits found
                  </td>
                </tr>
              ) : (
                permits.map((permit) => (
                  <tr
                    key={permit.id}
                    className="hover:bg-gray-50 dark:hover:bg-gray-800/50"
                  >
                    <td className="px-4 py-3 text-sm font-medium text-gray-800 dark:text-white">
                      {permit.permit_no}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-300">
                      {permit.permit_type?.name || "-"}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-300">
                      {permit.name}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-300">
                      {permit.domain?.name || "-"}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-300">
                      {permit.division?.name || "-"}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-300">
                      {permit.application_type}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-300">
                      {formatDate(permit.expiry_date)}
                    </td>
                    <td className="px-4 py-3 text-sm">
                      {getStatusBadge(permit)}
                    </td>
                    <td className="px-4 py-3 text-sm">
                      <div className="flex gap-2">
                        <Link href={`/permits/${permit.id}`}>
                          <button 
                            className="p-2 text-brand-500 transition-colors rounded-lg hover:bg-brand-50 hover:text-brand-600 dark:hover:bg-brand-900/20"
                            title="View"
                          >
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                            </svg>
                          </button>
                        </Link>
                        <Link href={`/permits/${permit.id}/edit`}>
                          <button 
                            className="p-2 text-blue-500 transition-colors rounded-lg hover:bg-blue-50 hover:text-blue-600 dark:hover:bg-blue-900/20"
                            title="Edit"
                          >
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                            </svg>
                          </button>
                        </Link>
                        <button
                          onClick={() => handleDelete(permit.id, permit.permit_no)}
                          className="p-2 text-error-500 transition-colors rounded-lg hover:bg-error-50 hover:text-error-600 dark:hover:bg-error-900/20"
                          title="Delete"
                        >
                          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                          </svg>
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between px-4 py-3 border-t border-gray-200 dark:border-gray-800">
            <div className="text-sm text-gray-500 dark:text-gray-400">
              Showing {(page - 1) * limit + 1} to{" "}
              {Math.min(page * limit, total)} of {total} results
            </div>
            <div className="flex gap-2">
              <Button
                onClick={() => setPage(page - 1)}
                disabled={page === 1}
                size="sm"
                variant="outline"
              >
                Previous
              </Button>
              <Button
                onClick={() => setPage(page + 1)}
                disabled={page >= totalPages}
                size="sm"
                variant="outline"
              >
                Next
              </Button>
            </div>
          </div>
        )}
      </div>

      {/* Delete Confirmation Modal */}
      <ConfirmModal
        isOpen={deleteModal.isOpen}
        onClose={() => setDeleteModal({ isOpen: false, permitId: null, permitNo: "" })}
        onConfirm={confirmDelete}
        title="Delete Permit"
        message={`Are you sure you want to delete permit "${deleteModal.permitNo}"? This action cannot be undone.`}
        confirmText="Delete"
        cancelText="Cancel"
        type="danger"
        loading={deleting}
      />
    </div>
  );
}
