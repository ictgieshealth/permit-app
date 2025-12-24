"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Input from "@/components/form/input/InputField";
import Link from "next/link";
import { Task } from "@/types/task";
import { taskService } from "@/services/task.service";
import { projectService } from "@/services/project.service";
import { ConfirmModal } from "@/components/ui/modal/ConfirmModal";
import TaskActionsDropdown from "@/components/ui/dropdown/TaskActionsDropdown";
import { useAuth } from "@/context/AuthContext";

export default function TasksPage() {
  const router = useRouter();
  const { currentRole } = useAuth();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [projects, setProjects] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState({
    search: "",
    project_id: "",
    status_id: "",
    approval_status_id: "",
    assigned_id: "",
  });
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
  });
  const [deleteModal, setDeleteModal] = useState<{
    isOpen: boolean;
    taskId: number | null;
    taskCode: string;
  }>({ isOpen: false, taskId: null, taskCode: "" });
  const [deleting, setDeleting] = useState(false);
  const [statusModal, setStatusModal] = useState<{
    isOpen: boolean;
    task: Task | null;
    action: string;
    statusId: number;
  }>({ isOpen: false, task: null, action: "", statusId: 0 });
  const [typeModal, setTypeModal] = useState<{
    isOpen: boolean;
    task: Task | null;
    action: string;
    typeId: number;
  }>({ isOpen: false, task: null, action: "", typeId: 0 });
  const [reasonModal, setReasonModal] = useState<{
    isOpen: boolean;
    task: Task | null;
    reason: string;
  }>({ isOpen: false, task: null, reason: "" });
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
      const response = await taskService.getAll({
        page: pagination.page,
        limit: pagination.limit,
        search: filters.search || undefined,
        project_id: filters.project_id
          ? parseInt(filters.project_id)
          : undefined,
        status_id: filters.status_id ? parseInt(filters.status_id) : undefined,
        approval_status_id: filters.approval_status_id
          ? parseInt(filters.approval_status_id)
          : undefined,
        assigned_id: filters.assigned_id
          ? parseInt(filters.assigned_id)
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

  const handleDelete = (id: number, code: string) => {
    setDeleteModal({ isOpen: true, taskId: id, taskCode: code });
  };

  const confirmDelete = async () => {
    if (!deleteModal.taskId) return;

    setDeleting(true);
    try {
      await taskService.delete(deleteModal.taskId);
      setDeleteModal({ isOpen: false, taskId: null, taskCode: "" });
      loadTasks();
    } catch (err: any) {
      alert(err.message || "Failed to delete task");
    } finally {
      setDeleting(false);
    }
  };

  // Status change handlers
  const handleStatusChange = (task: Task, action: string, statusId: number) => {
    if (action === "done" && task.status_id !== 3 && task.status_id !== 37) {
      alert("Task must be in 'On Progress' or 'In Review' status to be marked as done.");
      return;
    }
    setStatusModal({ isOpen: true, task, action, statusId });
  };

  const confirmStatusChange = async () => {
    if (!statusModal.task) return;

    setProcessing(true);
    try {
      await taskService.changeStatus(statusModal.task.id, {
        status_id: statusModal.statusId,
      });
      setStatusModal({ isOpen: false, task: null, action: "", statusId: 0 });
      loadTasks();
    } catch (err: any) {
      alert(err.message || "Failed to change task status");
    } finally {
      setProcessing(false);
    }
  };

  // Type change handlers
  const handleTypeChange = (task: Task, action: string, typeId: number) => {
    setTypeModal({ isOpen: true, task, action, typeId });
  };

  const confirmTypeChange = async () => {
    if (!typeModal.task) return;

    setProcessing(true);
    try {
      await taskService.changeType(typeModal.task.id, {
        type_id: typeModal.typeId,
      });
      setTypeModal({ isOpen: false, task: null, action: "", typeId: 0 });
      loadTasks();
    } catch (err: any) {
      alert(err.message || "Failed to change task type");
    } finally {
      setProcessing(false);
    }
  };

  // Reason handler
  const handleAddReason = (task: Task) => {
    setReasonModal({ isOpen: true, task, reason: "" });
  };

  const confirmAddReason = async () => {
    if (!reasonModal.task || !reasonModal.reason) {
      alert("Please enter a reason");
      return;
    }

    setProcessing(true);
    try {
      await taskService.setReason(reasonModal.task.id, {
        reason: reasonModal.reason,
      });
      setReasonModal({ isOpen: false, task: null, reason: "" });
      loadTasks();
    } catch (err: any) {
      alert(err.message || "Failed to add reason");
    } finally {
      setProcessing(false);
    }
  };

  // Action handlers
  const handleTodo = (task: Task) => {
    handleStatusChange(task, "todo", 1);
  };

  const handleStart = (task: Task) => {
    handleStatusChange(task, "start", 3);
  };

  const handleOnHold = (task: Task) => {
    handleStatusChange(task, "onHold", 2);
  };

  const handleDone = (task: Task) => {
    handleStatusChange(task, "done", 4);
  };

  const handleInReview = (task: Task) => {
    router.push(`/tasks/in-review/${task.code}`);
  };

  const handleShow = (task: Task) => {
    router.push(`/tasks/show/${task.code}`);
  };

  const handleEdit = (task: Task) => {
    router.push(`/tasks/edit/${task.code}`);
  };

  const handleRevision = (task: Task) => {
    router.push(`/tasks/revision/${task.code}`);
  };

  const handleDevelopment = (task: Task) => {
    handleTypeChange(task, "development", 25);
  };

  const handleMaintenance = (task: Task) => {
    handleTypeChange(task, "maintenance", 24);
  };

  const getStatusBadge = (statusName?: string) => {
    const colors: Record<string, string> = {
      "To Do": "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300",
      "On Hold": "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
      "On Progress": "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300",
      Done: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
      "In Review": "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300",
      Revision: "bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300",
    };
    const colorClass =
      colors[statusName || ""] ||
      "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300";
    return (
      <span
        className={`px-2 py-1 text-xs font-medium rounded-full ${colorClass}`}
      >
        {statusName || "-"}
      </span>
    );
  };

  const getApprovalBadge = (approvalName?: string) => {
    const colors: Record<string, string> = {
      Waiting: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
      Reject: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
      Approve: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
      "Pending Manager": "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300",
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

  const getPriorityBadge = (priorityName?: string) => {
    const colors: Record<string, string> = {
      Low: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
      Medium: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
      High: "bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300",
      Critical: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
    };
    const colorClass =
      colors[priorityName || ""] ||
      "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300";
    return (
      <span
        className={`px-2 py-1 text-xs font-medium rounded-full ${colorClass}`}
      >
        {priorityName || "-"}
      </span>
    );
  };

  const totalPages = Math.ceil(pagination.total / pagination.limit);

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
            Task Management
          </h1>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Manage and track project tasks
          </p>
        </div>
        <Link href="/tasks/add">
          <Button>Create New Task</Button>
        </Link>
      </div>

      {/* Filters */}
      <div className="p-4 mb-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
          <div>
            <label className="block mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
              Search
            </label>
            <Input
              name="search"
              type="text"
              value={filters.search}
              onChange={handleFilterChange}
              placeholder="Search by code or title"
            />
          </div>
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
                  status_id: "",
                  approval_status_id: "",
                  assigned_id: "",
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
      <div className="bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
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
                  Priority
                </th>
                <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Status
                </th>
                <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Approval
                </th>
                <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Assignee
                </th>
                <th className="px-6 py-3 text-xs font-medium tracking-wider text-left text-gray-500 uppercase dark:text-gray-400">
                  Due Date
                </th>
                <th className="px-6 py-3 text-xs font-medium tracking-wider text-right text-gray-500 uppercase dark:text-gray-400">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200 dark:bg-gray-dark dark:divide-gray-700">
              {loading ? (
                <tr>
                  <td colSpan={9} className="px-6 py-4 text-center">
                    <div className="flex items-center justify-center">
                      <div className="w-8 h-8 border-4 border-blue-600 rounded-full border-t-transparent animate-spin"></div>
                    </div>
                  </td>
                </tr>
              ) : tasks.length === 0 ? (
                <tr>
                  <td
                    colSpan={9}
                    className="px-6 py-4 text-sm text-center text-gray-500 dark:text-gray-400"
                  >
                    No tasks found
                  </td>
                </tr>
              ) : (
                tasks.map((task, index) => (
                  <tr key={task.id} className="hover:bg-gray-50 dark:hover:bg-gray-800">
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
                      {getPriorityBadge(task.priority?.name)}
                    </td>
                    <td className="px-6 py-4 text-sm whitespace-nowrap">
                      {getStatusBadge(task.status_task?.name)}
                    </td>
                    <td className="px-6 py-4 text-sm whitespace-nowrap">
                      {getApprovalBadge(task.approval_status?.name)}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500 whitespace-nowrap dark:text-gray-400">
                      {task.assignee?.full_name || "-"}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500 whitespace-nowrap dark:text-gray-400">
                      {task.due_date
                        ? new Date(task.due_date).toLocaleDateString()
                        : "-"}
                    </td>
                    <td className="px-6 py-4 text-sm font-medium text-right whitespace-nowrap">
                      <TaskActionsDropdown
                        task={task}
                        index={index}
                        totalRows={tasks.length}
                        userRole={currentRole?.name || null}
                        onTodo={handleTodo}
                        onStart={handleStart}
                        onHold={handleOnHold}
                        onDone={handleDone}
                        onInReview={handleInReview}
                        onShow={handleShow}
                        onAddReason={handleAddReason}
                        onRevision={handleRevision}
                        onEdit={handleEdit}
                        onDelete={() => handleDelete(task.id, task.code)}
                        onDevelopment={handleDevelopment}
                        onMaintenance={handleMaintenance}
                      />
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between px-6 py-3 bg-gray-50 dark:bg-gray-800">
            <div className="text-sm text-gray-700 dark:text-gray-300">
              Showing {(pagination.page - 1) * pagination.limit + 1} to{" "}
              {Math.min(pagination.page * pagination.limit, pagination.total)} of{" "}
              {pagination.total} results
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

      {/* Delete Confirmation Modal */}
      <ConfirmModal
        isOpen={deleteModal.isOpen}
        onClose={() =>
          setDeleteModal({ isOpen: false, taskId: null, taskCode: "" })
        }
        onConfirm={confirmDelete}
        title="Delete Task"
        message={`Are you sure you want to delete task "${deleteModal.taskCode}"? This action cannot be undone.`}
        confirmText="Delete"
        cancelText="Cancel"
        loading={deleting}
      />

      {/* Status Change Confirmation Modal */}
      <ConfirmModal
        isOpen={statusModal.isOpen}
        onClose={() =>
          setStatusModal({ isOpen: false, task: null, action: "", statusId: 0 })
        }
        onConfirm={confirmStatusChange}
        title={`Change Task Status to ${statusModal.action}`}
        message={`Are you sure you want to change task "${statusModal.task?.code}" to ${statusModal.action}?`}
        confirmText="Confirm"
        cancelText="Cancel"
        loading={processing}
      />

      {/* Type Change Confirmation Modal */}
      <ConfirmModal
        isOpen={typeModal.isOpen}
        onClose={() =>
          setTypeModal({ isOpen: false, task: null, action: "", typeId: 0 })
        }
        onConfirm={confirmTypeChange}
        title={`Change Task Type to ${typeModal.action}`}
        message={`Are you sure you want to change task "${typeModal.task?.code}" type to ${typeModal.action}?`}
        confirmText="Confirm"
        cancelText="Cancel"
        loading={processing}
      />

      {/* Add Reason Modal */}
      {reasonModal.isOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
          <div className="w-full max-w-md p-6 bg-white rounded-lg dark:bg-gray-800">
            <h2 className="mb-4 text-xl font-semibold text-gray-900 dark:text-white">
              Add Reason for Task "{reasonModal.task?.code}"
            </h2>
            <textarea
              className="w-full px-3 py-2 border border-gray-300 rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              rows={4}
              placeholder="Enter reason..."
              value={reasonModal.reason}
              onChange={(e) =>
                setReasonModal({ ...reasonModal, reason: e.target.value })
              }
            />
            <div className="flex justify-end gap-2 mt-4">
              <Button
                onClick={() =>
                  setReasonModal({ isOpen: false, task: null, reason: "" })
                }
                variant="outline"
                disabled={processing}
              >
                Cancel
              </Button>
              <Button onClick={confirmAddReason} disabled={processing}>
                {processing ? "Saving..." : "Save Reason"}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
