"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Link from "next/link";
import { Project } from "@/types/project";
import { projectService } from "@/services/project.service";
import { authService } from "@/services/auth.service";
import { ConfirmModal } from "@/components/ui/modal/ConfirmModal";
import {
  PlusIcon,
  PencilIcon,
  TrashBinIcon,
  CheckCircleIcon,
  ArrowRightIcon,
  TimeIcon,
} from "@/icons";

export default function ProjectsPage() {
  const router = useRouter();

  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [total, setTotal] = useState(0);
  const [filters, setFilters] = useState({
    name: "",
    code: "",
    status: "" as string | null,
    project_status_id: "" as string | null,
  });
  const [deleteModal, setDeleteModal] = useState<{
    isOpen: boolean;
    projectId: number | null;
    projectName: string;
  }>({ isOpen: false, projectId: null, projectName: "" });
  const [deleting, setDeleting] = useState(false);
  const [statusChangeModal, setStatusChangeModal] = useState<{
    isOpen: boolean;
    project: Project | null;
    statusId: number | null;
    statusName: string;
  }>({ isOpen: false, project: null, statusId: null, statusName: "" });
  const [changingStatus, setChangingStatus] = useState(false);
  const [selectedStatus, setSelectedStatus] = useState("all");

  useEffect(() => {
    loadProjects();
  }, [page, filters]);

  const loadProjects = async () => {
    try {
      setLoading(true);
      const response = await projectService.getAll({
        name: filters.name || undefined,
        code: filters.code || undefined,
        status: filters.status === "active" ? true : filters.status === "inactive" ? false : undefined,
        project_status_id: filters.project_status_id ? Number(filters.project_status_id) : undefined,
        page,
        limit,
      });

      setProjects(response.data);
      setTotal(response.meta?.total || 0);
    } catch (err: any) {
      setError(err.response?.data?.message || "Failed to load projects");
    } finally {
      setLoading(false);
    }
  };

  const handleFilterChange = (key: string, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
    setPage(1);
  };

  const handleDelete = async () => {
    if (!deleteModal.projectId) return;

    try {
      setDeleting(true);
      setError("");
      await projectService.delete(deleteModal.projectId);
      setDeleteModal({ isOpen: false, projectId: null, projectName: "" });
      setSuccessMessage("Project deleted successfully");
      setTimeout(() => setSuccessMessage(""), 3000);
      loadProjects();
    } catch (err: any) {
      setError(err.response?.data?.message || "Failed to delete project. Only projects with Pending status can be deleted.");
      setTimeout(() => setError(""), 5000);
    } finally {
      setDeleting(false);
    }
  };

  const handleChangeStatus = (
    project: Project,
    statusId: number,
    statusName: string
  ) => {
    setStatusChangeModal({
      isOpen: true,
      project,
      statusId,
      statusName,
    });
  };

  const confirmChangeStatus = async () => {
    if (!statusChangeModal.project || !statusChangeModal.statusId) return;

    try {
      setChangingStatus(true);
      setError("");
      await projectService.changeStatus(statusChangeModal.project.id, {
        status_id: statusChangeModal.statusId
      });
      setStatusChangeModal({ isOpen: false, project: null, statusId: null, statusName: "" });
      setSuccessMessage(`Project status changed to ${statusChangeModal.statusName} successfully`);
      setTimeout(() => setSuccessMessage(""), 3000);
      loadProjects();
    } catch (err: any) {
      setError(err.response?.data?.message || "Failed to change project status");
      setTimeout(() => setError(""), 5000);
    } finally {
      setChangingStatus(false);
    }
  };

  const filterByStatus = (status: string) => {
    setSelectedStatus(status);
    setFilters((prev) => ({
      ...prev,
      status: status === "all" ? null : status,
    }));
    setPage(1);
  };

  const getStatusBadge = (projectStatus: any) => {
    if (!projectStatus) return null;

    const statusMap: Record<string, { color: string; text: string }> = {
      "Pending": { color: "bg-yellow-100 text-yellow-800", text: "Pending" },
      "On Hold": { color: "bg-orange-100 text-orange-800", text: "On Hold" },
      "On Progress": { color: "bg-blue-100 text-blue-800", text: "On Progress" },
      "Done": { color: "bg-green-100 text-green-800", text: "Done" },
    };

    const status = statusMap[projectStatus.name] || {
      color: "bg-gray-100 text-gray-800",
      text: projectStatus.name,
    };

    return (
      <span
        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${status.color}`}
      >
        {status.text}
      </span>
    );
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <div className="p-6 max-w-full overflow-x-hidden">
      <div className="mb-6">
        <div className="flex flex-wrap justify-between items-center mb-4 gap-4">
          <div>
            <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">
              Project Management
            </h1>
            <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
              Manage your projects and track their progress
            </p>
          </div>
          <Link href="/projects/create">
            <Button variant="primary" size="md">
              <PlusIcon className="h-5 w-5 mr-2" />
              Create Project
            </Button>
          </Link>
        </div>

        {/* Status Filter */}
        <div className="flex flex-wrap gap-2 mb-4">
          <button
            onClick={() => filterByStatus("all")}
            className={`px-4 py-2 rounded-lg ${selectedStatus === "all"
                ? "bg-blue-600 text-white"
                : "bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300"
              }`}
          >
            All
          </button>
          <button
            onClick={() => filterByStatus("active")}
            className={`px-4 py-2 rounded-lg ${selectedStatus === "active"
                ? "bg-blue-600 text-white"
                : "bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300"
              }`}
          >
            Active
          </button>
          <button
            onClick={() => filterByStatus("inactive")}
            className={`px-4 py-2 rounded-lg ${selectedStatus === "inactive"
                ? "bg-blue-600 text-white"
                : "bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300"
              }`}
          >
            Inactive
          </button>
        </div>

        {/* Search Filters */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Name
            </label>
            <input
              type="text"
              value={filters.name}
              onChange={(e) => handleFilterChange("name", e.target.value)}
              placeholder="Search by name..."
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Code
            </label>
            <input
              type="text"
              value={filters.code}
              onChange={(e) => handleFilterChange("code", e.target.value)}
              placeholder="Search by code..."
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
            />
          </div>
        </div>
      </div>

      {error && (
        <div className="mb-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p className="text-red-800 dark:text-red-400">{error}</p>
        </div>
      )}

      {successMessage && (
        <div className="mb-4 p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg">
          <p className="text-green-800 dark:text-green-400">{successMessage}</p>
        </div>
      )}

      {loading ? (
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      ) : (
        <>
          <div className="overflow-x-auto bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
            <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                <thead className="bg-gray-50 dark:bg-gray-700">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      No
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Name
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Code
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Description
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Project Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Users
                    </th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                  {projects.length === 0 ? (
                    <tr>
                      <td
                        colSpan={8}
                        className="px-6 py-4 text-center text-gray-500 dark:text-gray-400"
                      >
                        No projects found
                      </td>
                    </tr>
                  ) : (
                    projects.map((project, index) => (
                      <tr
                        key={project.id}
                        className="hover:bg-gray-50 dark:hover:bg-gray-700"
                      >
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
                          {(page - 1) * limit + index + 1}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm font-medium text-gray-900 dark:text-white">
                            {project.name}
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-500 dark:text-gray-400">
                            {project.code}
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="text-sm text-gray-500 dark:text-gray-400 max-w-xs truncate">
                            {project.description || "-"}
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span
                            className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${project.status
                                ? "bg-green-100 text-green-800"
                                : "bg-red-100 text-red-800"
                              }`}
                          >
                            {project.status ? "Active" : "Inactive"}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          {getStatusBadge(project.project_status)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-500 dark:text-gray-400">
                            {project.users?.length || 0} users
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          <div className="flex justify-end gap-2">
                            {/* Status Change Buttons */}
                            {project.project_status?.name === "Pending" && (
                              <button
                                onClick={() =>
                                  handleChangeStatus(project, 28, "On Progress")
                                }
                                className="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
                                title="Start Project"
                              >
                                <ArrowRightIcon className="h-5 w-5" />
                              </button>
                            )}
                            {project.project_status?.name === "On Progress" && (
                              <>
                                <button
                                  onClick={() =>
                                    handleChangeStatus(project, 27, "On Hold")
                                  }
                                  className="text-orange-600 hover:text-orange-900 dark:text-orange-400 dark:hover:text-orange-300"
                                  title="Put On Hold"
                                >
                                  <TimeIcon className="h-5 w-5" />
                                </button>
                                <button
                                  onClick={() =>
                                    handleChangeStatus(project, 29, "Done")
                                  }
                                  className="text-green-600 hover:text-green-900 dark:text-green-400 dark:hover:text-green-300"
                                  title="Mark as Done"
                                >
                                  <CheckCircleIcon className="h-5 w-5" />
                                </button>
                              </>
                            )}
                            {project.project_status?.name === "On Hold" && (
                              <button
                                onClick={() =>
                                  handleChangeStatus(project, 28, "On Progress")
                                }
                                className="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
                                title="Resume Project"
                              >
                                <ArrowRightIcon className="h-5 w-5" />
                              </button>
                            )}

                            {/* Edit Button */}
                            <Link href={`/projects/edit/${project.id}`}>
                              <button className="text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-300">
                                <PencilIcon className="h-5 w-5" />
                              </button>
                            </Link>

                            {/* Delete Button - Only for Pending status */}
                            {project.project_status?.name === "Pending" && (
                              <button
                                onClick={() =>
                                  setDeleteModal({
                                    isOpen: true,
                                    projectId: project.id,
                                    projectName: project.name,
                                  })
                                }
                                className="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300"
                              >
                                <TrashBinIcon className="h-5 w-5" />
                              </button>
                            )}
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
            <div className="mt-4 flex flex-wrap justify-between items-center gap-4">
              <div className="text-sm text-gray-700 dark:text-gray-300">
                Showing {(page - 1) * limit + 1} to{" "}
                {Math.min(page * limit, total)} of {total} results
              </div>
              <div className="flex flex-wrap gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  disabled={page === 1}
                >
                  Previous
                </Button>
                <span className="px-4 py-2 text-sm text-gray-700 dark:text-gray-300">
                  Page {page} of {totalPages}
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                  disabled={page === totalPages}
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </>
      )}

      <ConfirmModal
        isOpen={deleteModal.isOpen}
        onClose={() =>
          setDeleteModal({ isOpen: false, projectId: null, projectName: "" })
        }
        onConfirm={handleDelete}
        title="Delete Project"
        message={`Are you sure you want to delete project "${deleteModal.projectName}"? Only projects with Pending status can be deleted. This action cannot be undone.`}
        confirmText="Delete"
        type="danger"
        loading={deleting}
      />

      <ConfirmModal
        isOpen={statusChangeModal.isOpen}
        onClose={() =>
          setStatusChangeModal({ isOpen: false, project: null, statusId: null, statusName: "" })
        }
        onConfirm={confirmChangeStatus}
        title="Change Project Status"
        message={`Are you sure you want to change project "${statusChangeModal.project?.name}" status to ${statusChangeModal.statusName}?`}
        confirmText="Change Status"
        type="warning"
        loading={changingStatus}
      />
    </div>
  );
}
