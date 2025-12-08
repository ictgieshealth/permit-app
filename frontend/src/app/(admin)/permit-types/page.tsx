"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Input from "@/components/form/input/InputField";
import Link from "next/link";
import { PermitType } from "@/types/permitType";
import { permitTypeService } from "@/services/permitType.service";
import { ConfirmModal } from "@/components/ui/modal/ConfirmModal";

export default function PermitTypesPage() {
  const router = useRouter();
  const [permitTypes, setPermitTypes] = useState<PermitType[]>([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState({
    code: "",
    name: "",
  });
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
  });
  const [deleteModal, setDeleteModal] = useState<{
    isOpen: boolean;
    permitTypeId: number | null;
    permitTypeName: string;
  }>({ isOpen: false, permitTypeId: null, permitTypeName: "" });
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    loadPermitTypes();
  }, [pagination.page, filters]);

  const loadPermitTypes = async () => {
    try {
      setLoading(true);
      const response = await permitTypeService.getAll({
        page: pagination.page,
        limit: pagination.limit,
        code: filters.code || undefined,
        name: filters.name || undefined,
      });
      setPermitTypes(response.data);
      setPagination({ ...pagination, total: response.meta?.total || 0 });
    } catch (err) {
      console.error("Failed to load permit types:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleFilterChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFilters({ ...filters, [name]: value });
    setPagination({ ...pagination, page: 1 });
  };

  const handleDelete = (id: number, name: string) => {
    setDeleteModal({ isOpen: true, permitTypeId: id, permitTypeName: name });
  };

  const confirmDelete = async () => {
    if (!deleteModal.permitTypeId) return;

    setDeleting(true);
    try {
      await permitTypeService.delete(deleteModal.permitTypeId);
      setDeleteModal({ isOpen: false, permitTypeId: null, permitTypeName: "" });
      loadPermitTypes();
    } catch (err: any) {
      alert(err.message || "Failed to delete permit type");
    } finally {
      setDeleting(false);
    }
  };

  const totalPages = Math.ceil(pagination.total / pagination.limit);

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
            Permit Type Management
          </h1>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Manage different types of permits in the system
          </p>
        </div>
        <Link href="/permit-types/create">
          <Button>Create New Permit Type</Button>
        </Link>
      </div>

      {/* Filters */}
      <div className="p-4 mb-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
          <Input
            name="code"
            type="text"
            value={filters.code}
            onChange={handleFilterChange}
            placeholder="Search by code"
          />
          <Input
            name="name"
            type="text"
            value={filters.name}
            onChange={handleFilterChange}
            placeholder="Search by name"
          />
          <div className="flex items-end">
            <Button
              onClick={() => {
                setFilters({ code: "", name: "" });
                setPagination({ ...pagination, page: 1 });
              }}
              className="w-full bg-gray-500 hover:bg-gray-600"
            >
              Reset Filters
            </Button>
          </div>
        </div>
      </div>

      {/* Table */}
      <div className="overflow-hidden bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th className="px-4 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Code
                </th>
                <th className="px-4 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Name
                </th>
                <th className="px-4 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Description
                </th>
                <th className="px-4 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200 dark:divide-gray-800">
              {loading ? (
                <tr>
                  <td colSpan={4} className="px-4 py-8 text-sm text-center text-gray-500">
                    Loading...
                  </td>
                </tr>
              ) : permitTypes.length === 0 ? (
                <tr>
                  <td colSpan={4} className="px-4 py-8 text-sm text-center text-gray-500">
                    No permit types found
                  </td>
                </tr>
              ) : (
                permitTypes.map((permitType) => (
                  <tr key={permitType.id} className="hover:bg-gray-50 dark:hover:bg-gray-800/50">
                    <td className="px-4 py-3 text-sm font-medium text-gray-800 dark:text-white">
                      {permitType.code}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-300">
                      {permitType.name}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-300">
                      {permitType.description || "-"}
                    </td>
                    <td className="px-4 py-3 text-sm">
                      <div className="flex gap-2">
                        <Link href={`/permit-types/${permitType.id}`}>
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
                        <Link href={`/permit-types/${permitType.id}/edit`}>
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
                          onClick={() => handleDelete(permitType.id, permitType.name)}
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
              Showing {(pagination.page - 1) * pagination.limit + 1} to{" "}
              {Math.min(pagination.page * pagination.limit, pagination.total)} of {pagination.total}{" "}
              results
            </div>
            <div className="flex gap-2">
              <Button
                onClick={() => setPagination({ ...pagination, page: pagination.page - 1 })}
                disabled={pagination.page === 1}
                className="bg-gray-500 hover:bg-gray-600"
              >
                Previous
              </Button>
              <Button
                onClick={() => setPagination({ ...pagination, page: pagination.page + 1 })}
                disabled={pagination.page >= totalPages}
                className="bg-gray-500 hover:bg-gray-600"
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
        onClose={() => setDeleteModal({ isOpen: false, permitTypeId: null, permitTypeName: "" })}
        onConfirm={confirmDelete}
        title="Delete Permit Type"
        message={`Are you sure you want to delete "${deleteModal.permitTypeName}"? This action cannot be undone.`}
        confirmText="Delete"
        cancelText="Cancel"
        type="danger"
        loading={deleting}
      />
    </div>
  );
}
