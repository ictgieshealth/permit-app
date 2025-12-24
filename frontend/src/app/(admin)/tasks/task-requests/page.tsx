"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Input from "@/components/form/input/InputField";
import Link from "next/link";
import { Task } from "@/types/task";
import { taskService } from "@/services/task.service";
import { projectService } from "@/services/project.service";
import ActionDropdown from "@/components/ui/dropdown/ActionDropdown";
import { ConfirmModal } from "@/components/ui/modal/ConfirmModal";
import { useAuth } from "@/context/AuthContext";
import { APPROVAL_STATUS, TASK_TYPE } from "@/constants/taskConstants";

export default function TaskRequestsPage() {
  const { user, currentRole } = useAuth();
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

  const [typeChangeModal, setTypeChangeModal] = useState<{
    isOpen: boolean;
    task: Task | null;
    newTypeId: number | null;
    typeName: string;
  }>({ isOpen: false, task: null, newTypeId: null, typeName: "" });
  const [changingType, setChangingType] = useState(false);

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
      (at) => at.approval_status_id === APPROVAL_STATUS.WAITING
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

  const handleChangeType = (task: Task, typeId: number, typeName: string) => {
    setTypeChangeModal({
      isOpen: true,
      task,
      newTypeId: typeId,
      typeName,
    });
  };

  const confirmChangeType = async () => {
    if (!typeChangeModal.task || !typeChangeModal.newTypeId) return;

    setChangingType(true);
    try {
      await taskService.update(typeChangeModal.task.id, {
        project_id: typeChangeModal.task.project_id,
        title: typeChangeModal.task.title,
        description: typeChangeModal.task.description,
        priority_id: typeChangeModal.task.priority_id!,
        stack_id: typeChangeModal.task.stack_id!,
        type_id: typeChangeModal.newTypeId,
        assigned_id: typeChangeModal.task.assigned_id,
        due_date: typeChangeModal.task.due_date,
      });

      setTypeChangeModal({
        isOpen: false,
        task: null,
        newTypeId: null,
        typeName: "",
      });
      loadTasks();
    } catch (err: any) {
      alert(err.message || "Failed to change task type");
    } finally {
      setChangingType(false);
    }
  };

  const handleView = (task: Task) => {
    router.push(`/tasks/show/${task.code}`);
  };

  const handleEdit = (task: Task) => {
    router.push(`/tasks/edit/${task.code}`);
  };

  const handleDelete = async (task: Task) => {
    if (!confirm(`Are you sure you want to delete task ${task.code}?`)) return;

    try {
      await taskService.delete(task.id);
      loadTasks();
    } catch (err: any) {
      alert(err.message || "Failed to delete task");
    }
  };

  // Role helper functions
  const hasAdminOrDev = () => {
    if (!currentRole) return false;
    return currentRole.name === "Ticketing Developer" || currentRole.name === "Admin";
  };

  const isTicketingPIC = () => {
    if (!currentRole) return false;
    return currentRole.name === "Ticketing PIC";
  };

  const isTicketingManager = () => {
    if (!currentRole) return false;
    return currentRole.name === "Ticketing Manager";
  };

  const isTicketingHeadOfUnit = () => {
    if (!currentRole) return false;
    return currentRole.name === "Ticketing Head of Unit";
  };

  // Permission checks - based on AdminTaskRequests.vue logic
  const canApprove = (task: Task) => {
    if (hasAdminOrDev()) return true;
    
    // Check manager approval
    const headUnitApproval = task.approval_tasks?.find(
      (at) => at.sequence === 1 && at.approval_status_id === APPROVAL_STATUS.WAITING
    );
    if (isTicketingManager()) {
      const managerApproval = task.approval_tasks?.find(
        (at) => at.sequence === 2 && at.approval_status_id === APPROVAL_STATUS.WAITING
      );
      
      if (managerApproval && !headUnitApproval) return true;
    }
    
    // Check head of unit approval
    if (isTicketingHeadOfUnit()) {
      if (headUnitApproval) return true;
    }
    
    return false;
  };

  const canReject = (task: Task) => {
    return canApprove(task);
  };

  const canEdit = (task: Task) => {
    if (hasAdminOrDev()) return true;
    if (isTicketingPIC() && task.approval_status?.name === "Waiting") return true;
    return false;
  };

  const canDelete = (task: Task) => {
    return isTicketingPIC() && task.approval_status?.name === "Waiting";
  };

  const canChangeType = () => {
    // Only non-admin/dev can change type
    return !hasAdminOrDev();
  };

  const shouldShowAction = (task: Task) => {
    return (
      canApprove(task) ||
      canReject(task) ||
      canEdit(task) ||
      canDelete(task) ||
      canChangeType()
    );
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
                    (at) => at.approval_status_id === APPROVAL_STATUS.WAITING
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
                        {shouldShowAction(task) && (
                          <ActionDropdown
                            data={task}
                            index={tasks.indexOf(task)}
                            totalRows={tasks.length}
                            onView={handleView}
                            onEdit={handleEdit}
                            onDelete={handleDelete}
                            onApprove={(t) => openApprovalModal(t, "approve")}
                            onReject={(t) => openApprovalModal(t, "reject")}
                            onDevelopment={(t) =>
                              handleChangeType(t, TASK_TYPE.DEVELOPMENT, "Development")
                            }
                            onMaintenance={(t) =>
                              handleChangeType(t, TASK_TYPE.MAINTENANCE, "Maintenance")
                            }
                            hideView={false}
                            hideEdit={!canEdit(task)}
                            hideDelete={!canDelete(task)}
                            hideApprove={!canApprove(task)}
                            hideReject={!canReject(task)}
                            hideDevelopment={!canChangeType()}
                            hideMaintenance={!canChangeType()}
                          />
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

      {/* Type Change Modal */}
      <ConfirmModal
        isOpen={typeChangeModal.isOpen}
        onClose={() =>
          setTypeChangeModal({
            isOpen: false,
            task: null,
            newTypeId: null,
            typeName: "",
          })
        }
        onConfirm={confirmChangeType}
        title="Change Task Type"
        message={`Are you sure you want to change task "${typeChangeModal.task?.code}" type to ${typeChangeModal.typeName}?`}
        confirmText="Change Type"
        type="warning"
        loading={changingType}
      />
    </div>
  );
}
