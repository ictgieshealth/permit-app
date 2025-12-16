"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Input from "@/components/form/input/InputField";
import Link from "next/link";
import { Task } from "@/types/task";
import { taskService } from "@/services/task.service";
import { projectService } from "@/services/project.service";

export default function TaskRequestsPage() {
  const router = useRouter();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [projects, setProjects] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState({
    search: "",
    project_id: "",
    approval_status_id: "",
  });
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
  });

  const [approvalModal, setApprovalModal] = useState<{
    isOpen: boolean;
    task: Task | null;
    approvalTaskId: number | null;
    action: "approve" | "reject" | null;
  }>({ isOpen: false, task: null, approvalTaskId: null, action: null });
  const [note, setNote] = useState("");
  const [processing, setProcessing] = useState(false);

  useEffect(() => {
    loadProjects();
    loadTasks();
  }, [pagination.page, filters]);

  const loadProjects = async () => {
    try {
      const response = await projectService.getAll({ limit: 100 });
      setProjects(response.data);
    } catch (err) {
      console.error("Failed to load projects:", err);
    }
  };

  const loadTasks = async () => {
    try {
      setLoading(true);
      // Use getAllRequests instead of getAll to get all tasks without approval_status_id filter
      const response = await taskService.getAllRequests({
        page: pagination.page,
        limit: pagination.limit,
        search: filters.search || undefined,
        project_id: filters.project_id
          ? parseInt(filters.project_id)
          : undefined,
        approval_status_id: filters.approval_status_id
          ? parseInt(filters.approval_status_id)
          : undefined,
      });
      setTasks(response.data);
      setPagination({ ...pagination, total: response.meta?.total || 0 });
    } catch (err) {
      console.error("Failed to load tasks:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleFilterChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setFilters({ ...filters, [name]: value });
    setPagination({ ...pagination, page: 1 });
  };

  const openApprovalModal = (
    task: Task,
    action: "approve" | "reject"
  ) => {
    // Find the next approval task that needs action
    const pendingApproval = task.approval_tasks?.find(
      (at) => at.approval_status_id === 20 // Waiting status
    );

    if (pendingApproval) {
      setApprovalModal({
        isOpen: true,
        task,
        approvalTaskId: pendingApproval.id,
        action,
      });
      setNote("");
    } else {
      alert("No pending approval found for this task");
    }
  };

  const handleApproval = async () => {
    if (!approvalModal.task || !approvalModal.approvalTaskId || !approvalModal.action)
      return;

    setProcessing(true);
    try {
      if (approvalModal.action === "approve") {
        await taskService.approveTask(
          approvalModal.task.id,
          approvalModal.approvalTaskId,
          { note: note || undefined }
        );
      } else {
        await taskService.rejectTask(
          approvalModal.task.id,
          approvalModal.approvalTaskId,
          { note: note || undefined }
        );
      }

      setApprovalModal({
        isOpen: false,
        task: null,
        approvalTaskId: null,
        action: null,
      });
      setNote("");
      loadTasks();
    } catch (err: any) {
      alert(err.message || "Failed to process approval");
    } finally {
      setProcessing(false);
    }
  };

  const getApprovalBadge = (approvalName?: string) => {
    const colors: Record<string, string> = {
      Waiting:
        "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
      Reject: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
      Approve:
        "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
      "Pending Manager":
        "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300",
    };
    const colorClass =
      colors[approvalName || ""] ||
      "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300";
    return (
      <span
        className={`px-2 py-1 text-xs font-medium rounded-full ${colorClass}`}
      >
        {approvalName || "-"}
      </span>
    );
  };

  const totalPages = Math.ceil(pagination.total / pagination.limit);

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
            Task Approval Requests
          </h1>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Review and approve or reject pending tasks
          </p>
        </div>
      </div>

      {/* Filters */}
      <div className="p-4 mb-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
          <Input
            name="search"
            type="text"
            value={filters.search}
            onChange={handleFilterChange}
            placeholder="Search by code or title"
          />
          <div>
            <label className="block mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
              Project
            </label>
            <select
              name="project_id"
              value={filters.project_id}
              onChange={handleFilterChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
            >
              <option value="">All Projects</option>
              {projects.map((project) => (
                <option key={project.id} value={project.id}>
                  {project.name}
                </option>
              ))}
            </select>
          </div>
          <div className="flex items-end">
            <Button
              onClick={() => {
                setFilters({
                  search: "",
                  project_id: "",
                  approval_status_id: "",
                });
                setPagination({ ...pagination, page: 1 });
              }}
              variant="outline"
            >
              Clear Filters
            </Button>
          </div>
        </div>
      </div>

      {/* Tasks Table */}
      <div className="overflow-hidden bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Code
                </th>
                <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Title
                </th>
                <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Project
                </th>
                <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Status
                </th>
                <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Created By
                </th>
                <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Pending Sequence
                </th>
                <th className="px-6 py-3 text-xs font-medium tracking-wider text-right text-gray-500 uppercase dark:text-gray-400">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200 dark:bg-gray-dark dark:divide-gray-700">
              {loading ? (
                <tr>
                  <td colSpan={7} className="px-6 py-4 text-center">
                    <div className="flex items-center justify-center">
                      <div className="w-8 h-8 border-4 border-blue-600 rounded-full border-t-transparent animate-spin"></div>
                    </div>
                  </td>
                </tr>
              ) : tasks.length === 0 ? (
                <tr>
                  <td
                    colSpan={7}
                    className="px-6 py-4 text-sm text-center text-gray-500 dark:text-gray-400"
                  >
                    No pending approval requests
                  </td>
                </tr>
              ) : (
                tasks.map((task) => {
                  const pendingApproval = task.approval_tasks?.find(
                    (at) => at.approval_status_id === 20
                  );

                  return (
                    <tr
                      key={task.id}
                      className="hover:bg-gray-50 dark:hover:bg-gray-800"
                    >
                      <td className="px-6 py-4 text-sm font-medium text-gray-900 whitespace-nowrap dark:text-white">
                        {task.code}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-900 dark:text-white">
                        <div className="max-w-xs truncate">{task.title}</div>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500 whitespace-nowrap dark:text-gray-400">
                        {task.project?.name || "-"}
                      </td>
                      <td className="px-6 py-4 text-sm whitespace-nowrap">
                        {getApprovalBadge(task.approval_status?.name)}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500 whitespace-nowrap dark:text-gray-400">
                        {task.creator?.full_name || "-"}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500 whitespace-nowrap dark:text-gray-400">
                        {pendingApproval
                          ? `Sequence ${pendingApproval.sequence}`
                          : "-"}
                      </td>
                      <td className="px-6 py-4 text-sm font-medium text-right whitespace-nowrap">
                        <Link
                          href={`/tasks/show/${task.code}`}
                          className="mr-3 text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
                        >
                          View
                        </Link>
                        {pendingApproval && (
                          <>
                            <button
                              onClick={() => openApprovalModal(task, "approve")}
                              className="mr-3 text-green-600 hover:text-green-900 dark:text-green-400 dark:hover:text-green-300"
                            >
                              Approve
                            </button>
                            <button
                              onClick={() => openApprovalModal(task, "reject")}
                              className="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300"
                            >
                              Reject
                            </button>
                          </>
                        )}
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between px-6 py-3 bg-gray-50 dark:bg-gray-800">
            <div className="text-sm text-gray-700 dark:text-gray-300">
              Showing {(pagination.page - 1) * pagination.limit + 1} to{" "}
              {Math.min(pagination.page * pagination.limit, pagination.total)}{" "}
              of {pagination.total} results
            </div>
            <div className="flex gap-2">
              <Button
                onClick={() =>
                  setPagination({ ...pagination, page: pagination.page - 1 })
                }
                disabled={pagination.page === 1}
                variant="outline"
                size="sm"
              >
                Previous
              </Button>
              <Button
                onClick={() =>
                  setPagination({ ...pagination, page: pagination.page + 1 })
                }
                disabled={pagination.page >= totalPages}
                variant="outline"
                size="sm"
              >
                Next
              </Button>
            </div>
          </div>
        )}
      </div>

      {/* Approval Modal */}
      {approvalModal.isOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
          <div className="w-full max-w-md p-6 bg-white rounded-lg dark:bg-gray-dark">
            <h2 className="mb-4 text-xl font-bold text-gray-800 dark:text-white">
              {approvalModal.action === "approve"
                ? "Approve Task"
                : "Reject Task"}
            </h2>
            <p className="mb-4 text-sm text-gray-600 dark:text-gray-400">
              Task: {approvalModal.task?.code} - {approvalModal.task?.title}
            </p>
            <div className="mb-4">
              <label className="block mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
                Note (Optional)
              </label>
              <textarea
                value={note}
                onChange={(e) => setNote(e.target.value)}
                rows={3}
                placeholder="Add a note for this action"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-800 dark:border-gray-700 dark:text-white"
              />
            </div>
            <div className="flex gap-3">
              <Button
                onClick={handleApproval}
                disabled={processing}
                className={
                  approvalModal.action === "approve"
                    ? "bg-green-600 hover:bg-green-700"
                    : "bg-red-600 hover:bg-red-700"
                }
              >
                {processing
                  ? "Processing..."
                  : approvalModal.action === "approve"
                    ? "Confirm Approve"
                    : "Confirm Reject"}
              </Button>
              <Button
                onClick={() =>
                  setApprovalModal({
                    isOpen: false,
                    task: null,
                    approvalTaskId: null,
                    action: null,
                  })
                }
                variant="outline"
                disabled={processing}
              >
                Cancel
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
