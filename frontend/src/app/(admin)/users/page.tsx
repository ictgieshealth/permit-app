"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Link from "next/link";
import { User } from "@/types/user";
import { Role } from "@/types/role";
import { Domain } from "@/types/domain";
import { authService } from "@/services/api.service";
import { userService } from "@/services/user.service";
import { roleService } from "@/services/role.service";
import { domainService } from "@/services/domain.service";

export default function UsersPage() {
  const router = useRouter();
  const [users, setUsers] = useState<User[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [domains, setDomains] = useState<Domain[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [filters, setFilters] = useState({
    username: "",
    email: "",
    role_id: "",
    domain_id: "",
    is_active: "",
  });
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
  });

  useEffect(() => {
    // Check if user is admin
    if (!authService.hasRole(["admin"])) {
      router.push("/dashboard");
      return;
    }

    loadRoles();
    loadDomains();
    loadUsers();
  }, [pagination.page, filters]);

  const loadUsers = async () => {
    try {
      setLoading(true);
      const params: any = {
        page: pagination.page,
        limit: pagination.limit,
      };

      if (filters.username) params.username = filters.username;
      if (filters.email) params.email = filters.email;
      if (filters.role_id) params.role_id = parseInt(filters.role_id);
      if (filters.domain_id) params.domain_id = parseInt(filters.domain_id);
      if (filters.is_active) params.is_active = filters.is_active === "true";

      const response = await userService.getAll(params);
      setUsers(response.data);
      if (response.meta) {
        setPagination((prev) => ({ ...prev, total: response.meta!.total }));
      }
    } catch (err: any) {
      setError(err.message || "Failed to load users");
    } finally {
      setLoading(false);
    }
  };

  const loadRoles = async () => {
    try {
      const response = await roleService.getAll({ limit: 100 });
      setRoles(response.data);
    } catch (err) {
      console.error("Failed to load roles:", err);
    }
  };

  const loadDomains = async () => {
    try {
      const response = await domainService.getAll({ limit: 100 });
      setDomains(response.data as unknown as Domain[]);
    } catch (err) {
      console.error("Failed to load domains:", err);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this user?")) return;

    try {
      await userService.delete(id);
      loadUsers();
    } catch (err: any) {
      alert(err.message || "Failed to delete user");
    }
  };

  const handleFilterChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFilters({
      ...filters,
      [e.target.name]: e.target.value,
    });
    setPagination((prev) => ({ ...prev, page: 1 }));
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    loadUsers();
  };

  const totalPages = Math.ceil(pagination.total / pagination.limit);

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-800 dark:text-white">User Management</h1>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Manage system users and their permissions
          </p>
        </div>
        <Link href="/users/create">
          <Button>Add New User</Button>
        </Link>
      </div>

      {/* Filters */}
      <div className="p-4 mb-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <form onSubmit={handleSearch} className="grid grid-cols-1 gap-4 md:grid-cols-5">
          <input
            type="text"
            name="username"
            placeholder="Search username..."
            value={filters.username}
            onChange={handleFilterChange}
            className="px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
          />
          <input
            type="text"
            name="email"
            placeholder="Search email..."
            value={filters.email}
            onChange={handleFilterChange}
            className="px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
          />
          <select
            name="role_id"
            value={filters.role_id}
            onChange={handleFilterChange}
            className="px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
          >
            <option value="">All Roles</option>
            {roles.map((role) => (
              <option key={role.id} value={role.id}>
                {role.name}
              </option>
            ))}
          </select>
          <select
            name="domain_id"
            value={filters.domain_id}
            onChange={handleFilterChange}
            className="px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
          >
            <option value="">All Domains</option>
            {domains.map((domain) => (
              <option key={domain.id} value={domain.id}>
                {domain.name}
              </option>
            ))}
          </select>
          <select
            name="is_active"
            value={filters.is_active}
            onChange={handleFilterChange}
            className="px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
          >
            <option value="">All Status</option>
            <option value="true">Active</option>
            <option value="false">Inactive</option>
          </select>
        </form>
      </div>

      {error && (
        <div className="p-4 mb-6 border rounded-lg bg-error-50 border-error-200 text-error-700 dark:bg-error-900/20">
          {error}
        </div>
      )}

      {/* Users Table */}
      <div className="overflow-hidden bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <div className="overflow-x-auto">
          <table className="w-full text-sm text-left">
            <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-800 dark:text-gray-400">
              <tr>
                <th className="px-6 py-3">Username</th>
                <th className="px-6 py-3">Full Name</th>
                <th className="px-6 py-3">Email</th>
                <th className="px-6 py-3">Role</th>
                <th className="px-6 py-3">Domains</th>
                <th className="px-6 py-3">Status</th>
                <th className="px-6 py-3">Actions</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={7} className="px-6 py-8 text-center text-gray-500">
                    Loading...
                  </td>
                </tr>
              ) : users.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-6 py-8 text-center text-gray-500">
                    No users found
                  </td>
                </tr>
              ) : (
                users.map((user) => (
                  <tr
                    key={user.id}
                    className="border-b dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800/50"
                  >
                    <td className="px-6 py-4 font-medium text-gray-900 dark:text-white">
                      {user.username}
                    </td>
                    <td className="px-6 py-4 text-gray-700 dark:text-gray-300">
                      {user.full_name}
                    </td>
                    <td className="px-6 py-4 text-gray-700 dark:text-gray-300">
                      {user.email}
                    </td>
                    <td className="px-6 py-4">
                      <span className="inline-block px-2 py-1 text-xs font-medium rounded-full bg-brand-100 text-brand-700 dark:bg-brand-900/20 dark:text-brand-400">
                        {user.role?.name || "N/A"}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex flex-wrap gap-1">
                        {user.domains && user.domains.length > 0 ? (
                          user.domains.map((domain) => (
                            <span
                              key={domain.id}
                              className="inline-block px-2 py-0.5 text-xs rounded-full bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300"
                            >
                              {domain.code}
                            </span>
                          ))
                        ) : (
                          <span className="text-xs text-gray-500">No domains</span>
                        )}
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <span
                        className={`inline-block px-2 py-1 text-xs font-medium rounded-full ${
                          user.is_active
                            ? "bg-success-100 text-success-700 dark:bg-success-900/20 dark:text-success-400"
                            : "bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300"
                        }`}
                      >
                        {user.is_active ? "Active" : "Inactive"}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex gap-2">
                        <Link
                          href={`/users/${user.id}/edit`}
                          className="text-brand-600 hover:text-brand-700 dark:text-brand-400"
                        >
                          Edit
                        </Link>
                        <button
                          onClick={() => handleDelete(user.id)}
                          className="text-error-600 hover:text-error-700 dark:text-error-400"
                        >
                          Delete
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
          <div className="flex items-center justify-between px-6 py-4 border-t dark:border-gray-800">
            <div className="text-sm text-gray-700 dark:text-gray-300">
              Showing {(pagination.page - 1) * pagination.limit + 1} to{" "}
              {Math.min(pagination.page * pagination.limit, pagination.total)} of{" "}
              {pagination.total} results
            </div>
            <div className="flex gap-2">
              <button
                onClick={() => setPagination((prev) => ({ ...prev, page: prev.page - 1 }))}
                disabled={pagination.page === 1}
                className="px-3 py-1 text-sm border rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 dark:border-gray-700 dark:hover:bg-gray-800"
              >
                Previous
              </button>
              <button
                onClick={() => setPagination((prev) => ({ ...prev, page: prev.page + 1 }))}
                disabled={pagination.page === totalPages}
                className="px-3 py-1 text-sm border rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 dark:border-gray-700 dark:hover:bg-gray-800"
              >
                Next
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
